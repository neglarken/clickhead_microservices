package httpproxy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/neglarken/clickhead/api-gate/config"
	authMsService "github.com/neglarken/clickhead/api-gate/services/auth-ms/protobuf"
	someMsService "github.com/neglarken/clickhead/api-gate/services/some-ms/protobuf"
)

func Start(cfg *config.Config) error {

	grpcGwMux := runtime.NewServeMux()

	//----------------------------------------------------------------
	// настройка подключений со стороны gRPC
	//----------------------------------------------------------------
	//Подключение к сервису SomeMs
	grpcSomeMsConn, err := grpc.Dial(
		cfg.SomeMsServerAddress,
		//grpc.WithPerRPCCredentials(&reqData{}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

	mux.Handle("/api/v1/", grpcGwMux)
	mux.HandleFunc("/", helloworld)

	fmt.Println("starting HTTP server at " + cfg.Port)
	return http.ListenAndServe(cfg.Port, mux)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "URL:", r.URL.String())
}

// func AccessLogInterceptor(
// 	ctx context.Context,
// 	method string,
// 	req interface{},
// 	reply interface{},
// 	cc *grpc.ClientConn,
// 	invoker grpc.UnaryInvoker,
// 	opts ...grpc.CallOption,
// ) error {
// 	md,_:=metadata.FromOutgoingContext(ctx)
// 	start:=time.Now()

// 	var traceId,userId,userRole string
// 	if len(md["authorization"])>0{
// 		tokenString:= md["authorization"][0]
// 		if tokenString!=""{
// 			err,token:=userService.CheckGetJWTToken(tokenString)
// 			if err!=nil{
// 				return err
// 			}
// 			userId=fmt.Sprintf("%s",token["UserID"])
// 			userRole=fmt.Sprintf("%s",token["UserRole"])
// 		}
// 	}
// 	//Присваиваю ID запроса
// 	traceId=fmt.Sprintf("%d",time.Now().UTC().UnixNano())

// 	callContext:=context.Background()
// 	mdOut:=metadata.Pairs(
// 		"trace-id",traceId,
// 		"user-id",userId,
// 		"user-role",userRole,
// 	)
// 	callContext=metadata.NewOutgoingContext(callContext,mdOut)

// 	err:=invoker(callContext,method,req,reply,cc, opts...)

// 	msg:=fmt.Sprintf("Call:%v, traceId: %v, userId: %v, userRole: %v, time: %v", method,traceId,userId,userRole,time.Since(start))
// 	app.AccesLog(msg)

// 	return err
// }
