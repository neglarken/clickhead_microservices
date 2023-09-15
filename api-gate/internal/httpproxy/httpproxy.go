package httpproxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/neglarken/clickhead/api-gate/config"
	"github.com/neglarken/clickhead/api-gate/pkg/auth"
	authMsService "github.com/neglarken/clickhead/api-gate/services/auth-ms/protobuf"
	someMsService "github.com/neglarken/clickhead/api-gate/services/some-ms/protobuf"
)

var (
	signingKey = "qwerty123"
	opts       []grpc.DialOption
)

func Start(cfg *config.Config) error {

	grpcGwMux := runtime.NewServeMux()

	//----------------------------------------------------------------
	// настройка подключений со стороны gRPC
	//----------------------------------------------------------------
	//Подключение к сервису SomeMs
	grpcSomeMsConn, err := grpc.Dial(
		cfg.SomeMsServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(AccessInterceptor),
	)
	if err != nil {
		return fmt.Errorf("Failed to connect to User service: %s", err)
	}
	defer grpcSomeMsConn.Close()

	err = someMsService.RegisterSomeMsServiceHandler(
		context.Background(),
		grpcGwMux,
		grpcSomeMsConn,
	)
	if err != nil {
		return fmt.Errorf("Failed to start HTTP server: %s", err)
	}

	//----------------------------------------------------------------
	//Подключение к сервису AuthMs
	grpcAuthMsConn, err := grpc.Dial(
		cfg.AuthMsServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("Failed to connect to Post service: %s", err)
	}
	defer grpcAuthMsConn.Close()

	err = authMsService.RegisterAuthMsServiceHandler(
		context.Background(),
		grpcGwMux,
		grpcAuthMsConn,
	)
	if err != nil {
		return fmt.Errorf("Failed to start HTTP server: %s", err)
	}

	//----------------------------------------------------------------
	//	Настройка маршрутов с стороны REST
	//----------------------------------------------------------------
	mux := http.NewServeMux()

	mux.Handle("/", grpcGwMux)

	fmt.Println("starting HTTP server at " + cfg.Port)
	return http.ListenAndServe(cfg.Port, mux)
}

func AccessInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	md, _ := metadata.FromOutgoingContext(ctx)

	if len(md["authorization"]) == 0 {
		return errors.New("authorization token is not provided")
	}
	tokenString := md["authorization"][0]
	values := strings.Split(tokenString, " ")
	if values[0] != "Bearer" {
		return errors.New("token is not a Bearer token")
	}
	tokenString = values[1]
	if tokenString == "" {
		return errors.New("token is empty")
	}

	userId, err := auth.Parse(tokenString, signingKey)
	if err != nil {
		return err
	}

	callContext := context.Background()

	mdOut := metadata.Pairs(
		"user-id", userId,
	)
	callContext = metadata.NewOutgoingContext(callContext, mdOut)

	err = invoker(callContext, method, req, reply, cc, opts...)
	if err != nil {
		return err
	}

	return nil
}
