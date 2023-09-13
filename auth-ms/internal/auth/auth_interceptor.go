package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager        TokenManager
	accessibleMethods map[string]bool
}

func NewAuthInterceptor(jwtManager TokenManager, accessibleMethods map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleMethods}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		id, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		if id == "" {
			return handler(ctx, req)
		}
		return handler(
			context.WithValue(ctx, "id", id),
			req,
		)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (string, error) {
	accessibleMethod, ok := interceptor.accessibleMethods[method]
	if !ok {
		return "", status.Error(codes.PermissionDenied, "no permission to access this RPC")
	}

	if accessibleMethod {
		return "", nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	values = strings.Split(values[0], " ")
	if values[0] != "Bearer" {
		return "", status.Errorf(codes.Unauthenticated, "token is not a Bearer token")
	}

	accessToken := values[1]

	id, err := interceptor.jwtManager.Parse(accessToken)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return id, nil
}
