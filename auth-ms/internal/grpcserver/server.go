package grpcserver

import (
	"net"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/neglarken/clickhead/auth-ms/config"
	"github.com/neglarken/clickhead/auth-ms/internal/auth"
	"github.com/neglarken/clickhead/auth-ms/internal/hasher"
	"github.com/neglarken/clickhead/auth-ms/internal/postgres"
	"github.com/neglarken/clickhead/auth-ms/internal/repo"
	"github.com/neglarken/clickhead/auth-ms/internal/service"
	pb "github.com/neglarken/clickhead/auth-ms/protobuf"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

var (
	customFunc = func(code codes.Code) logrus.Level {
		if code == codes.OK {
			return logrus.InfoLevel
		}
		return logrus.ErrorLevel
	}
)

type Server struct {
	pb.UnimplementedAuthMsServiceServer
	port    string
	service *service.Service
	logger  *logrus.Logger
}

func NewServer(port string, service *service.Service) *Server {
	s := &Server{
		port:    port,
		service: service,
		logger:  logrus.New(),
	}
	return s
}

func Start(cfg *config.Config) error {
	db, err := postgres.NewDB(cfg.URL)
	if err != nil {
		return err
	}

	defer db.Close()

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		return err
	}

	r := repo.NewRepository(db)
	hasher := hasher.NewSHA1Hasher(cfg.Salt)
	manager, err := auth.NewManager(cfg.SignKey)
	if err != nil {
		return err
	}
	serv := service.NewService(*r, hasher, manager, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	authInterceptor := auth.NewAuthInterceptor(manager, accessibleMethods())

	s := NewServer(cfg.Port, serv)

	logrusEntry := logrus.NewEntry(s.logger)

	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(customFunc),
	}

	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
		),
	)
	reflection.Register(grpcServer)

	pb.RegisterAuthMsServiceServer(grpcServer, s)

	s.logger.Infof("Starting server on %s", cfg.Port)

	return grpcServer.Serve(lis)
}

func accessibleMethods() map[string]bool {
	const authServicePath = "/protobuf.AuthMsService/"
	return map[string]bool{
		authServicePath + "SignUp": true,
		authServicePath + "SignIn": true,
		authServicePath + "WhoAmI": false,
	}
}
