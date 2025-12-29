package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/audit"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/handler"
	"github.com/RussiaFPS/shortlink/internal/logger"
	"github.com/RussiaFPS/shortlink/internal/repository"
	"github.com/RussiaFPS/shortlink/internal/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

const lenShortURL = 8

func main() {
	printOrDefault := func(value string) string {
		if value == "" {
			return "N/A"
		}
		return value
	}

	fmt.Printf("Build version: %s\n", printOrDefault(buildVersion))
	fmt.Printf("Build date: %s\n", printOrDefault(buildDate))
	fmt.Printf("Build commit: %s\n", printOrDefault(buildCommit))

	cfg := config.NewConfig()

	if err := logger.NewLogger(); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	auditService := audit.NewAuditService()
	if cfg.AuditFile != "" {
		auditService.Register(audit.NewLogObserver(cfg.AuditFile))
	}
	if cfg.AuditURL != "" {
		auditService.Register(audit.NewHTTPObserver(cfg.AuditURL))
	}

	router := gin.Default()
	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, lenShortURL, rep, auditService)
	handler.NewURLShortenerHandler(router, cfg, ser)

	defer rep.Close()
	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router.Handler(),
	}

	go func() {
		var err error
		if cfg.EnableHTTPS {
			log.Printf("Starting server with HTTPS")
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			log.Printf("Starting server without HTTPS")
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started at %v", cfg.ServerAddr)

	// gRPC server
	go func() {
		listen, err := net.Listen("tcp", cfg.GRPCAddr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))
		handler.NewGRPCURLShortenerHandler(s, ser)

		log.Printf("gRPC server listening at %v", listen.Addr())
		if err := s.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, uHandler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, errors.New("authorization token is not provided")
	}

	ctx = context.WithValue(ctx, "userID", values[0])
	return uHandler(ctx, req)
}
