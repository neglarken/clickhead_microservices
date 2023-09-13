package grpcserver

import (
	"net"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/neglarken/clickhead/some-ms/config"
	"github.com/neglarken/clickhead/some-ms/internal/postgres"
	"github.com/neglarken/clickhead/some-ms/internal/repo"
	pb "github.com/neglarken/clickhead/some-ms/protobuf"
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
	pb.UnimplementedSomeMsServiceServer
	Port   string
	Repo   *repo.Repository
	Logger *logrus.Logger
}

func NewServer(port string, repo *repo.Repository) *Server {
	s := &Server{
		Port:   port,
		Repo:   repo,
		Logger: logrus.New(),
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
	s := NewServer(cfg.Port, r)

	logrusEntry := logrus.NewEntry(s.Logger)

	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(customFunc),
	}

	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
		),
	)
	reflection.Register(grpcServer)
	pb.RegisterSomeMsServiceServer(grpcServer, s)

	s.Logger.Infof("Starting server on %s", cfg.Port)
	return grpcServer.Serve(lis)
}
