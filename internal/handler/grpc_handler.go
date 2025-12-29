package handler

import (
	"context"
	"github.com/RussiaFPS/shortlink/api/proto"
	"github.com/RussiaFPS/shortlink/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCURLShortenerHandler struct {
	proto.UnimplementedShortenerServiceServer
	urlShortener service.URLService
}

func NewGRPCURLShortenerHandler(s *grpc.Server, urlShortener service.URLService) {
	grpcHandler := &GRPCURLShortenerHandler{
		urlShortener: urlShortener,
	}
	proto.RegisterShortenerServiceServer(s, grpcHandler)
}

func (h *GRPCURLShortenerHandler) ShortenURL(ctx context.Context, r *proto.URLShortenRequest) (*proto.URLShortenResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}

	shortURL, ok, err := h.urlShortener.Shorten(ctx, r.Url, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if ok {
		return nil, status.Error(codes.AlreadyExists, "this url already exits")
	}

	return &proto.URLShortenResponse{Result: shortURL}, nil
}

func (h *GRPCURLShortenerHandler) ExpandURL(ctx context.Context, r *proto.URLExpandRequest) (*proto.URLExpandResponse, error) {
	longURL, exists := h.urlShortener.GetOriginal(ctx, r.Id, ctx.Value(UserIDKey).(string))
	if !exists {
		return nil, status.Error(codes.NotFound, "url not found")
	}

	return &proto.URLExpandResponse{Result: longURL.OriginalURL}, nil
}

func (h *GRPCURLShortenerHandler) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*proto.UserURLsResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}

	urls, err := h.urlShortener.GetShortenedURLByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var result []*proto.URLData
	for _, u := range urls {
		result = append(result, &proto.URLData{
			ShortUrl:    u.ShortURL,
			OriginalUrl: u.OriginalURL,
		})
	}
	return &proto.UserURLsResponse{Url: result}, nil
}
