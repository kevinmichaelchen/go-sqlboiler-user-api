package kontext

import (
	"context"

	"github.com/kevinmichaelchen/go-sqlboiler-user-api/pkg/headers"
	"github.com/rs/zerolog/log"
	grpcGo "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrParseMetadataFailure = status.Error(codes.PermissionDenied, "failed to parse grpc metadata headers")
)

type Key string

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpcGo.UnaryServerInfo, handler grpcGo.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrParseMetadataFailure
	}

	var userID string
	r := md[headers.UserID]
	if len(r) > 0 {
		userID = r[0]
	}

	ctx = context.WithValue(ctx, Key(headers.UserID), userID)

	return handler(ctx, req)
}

func GetRequesterID(ctx context.Context) string {
	v := ctx.Value(Key(headers.UserID))
	if s, ok := v.(string); ok {
		return s
	}

	// TODO this can be handled robustly... beyond merely logging...
	log.Warn().Msgf("requester ID was not a string: %v", v)

	return ""
}
