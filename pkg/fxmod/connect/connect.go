package connect

import (
	"context"
	"errors"
	"fmt"
	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/kevinmichaelchen/chomp-proxy/pkg/cors"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"time"
)

func CreateModule(opts *ModuleOptions) fx.Option {
	return fx.Module("grpc",
		fx.Provide(
			opts.HandlerProvider,
			func() *ModuleOptions {
				return opts
			},
			NewConfig,
			NewServer,
		),
		fx.Invoke(
			Register,
		),
	)
}

type HandlerOutput struct {
	Path    string
	Handler http.Handler
}

type ModuleOptions struct {
	HandlerProvider any
	ServiceName     string
	Services        []string
}

type Config struct {
	ConnectConfig *NestedConfig `env:",prefix=GRPC_CONNECT_"`
}

type NestedConfig struct {
	Host string `env:"HOST,default=localhost"`
	Port int    `env:"PORT,required"`
}

func NewConfig() (cfg Config, err error) {
	err = envconfig.Process(context.Background(), &cfg)
	return
}

func NewServer(lc fx.Lifecycle, cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	addr := fmt.Sprintf("%s:%d", cfg.ConnectConfig.Host, cfg.ConnectConfig.Port)
	srv := &http.Server{
		Addr: addr,
		// Use h2c, so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(
			cors.NewCORS().Handler(mux),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go func() {
				err := srv.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logrus.WithError(err).Error("connect-go ListenAndServe failed")
				}
			}()
			logrus.WithField("address", addr).Info("Listening for connect-go")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return mux
}

func Register(opts *ModuleOptions, mux *http.ServeMux, h HandlerOutput) {
	checker := grpchealth.NewStaticChecker(
		// protoc-gen-connect-go generates package-level constants
		// for these fully-qualified protobuf service names, so we'd be able
		// to reference foov1beta1.FooService as opposed to foo.v1beta1.FooService.
		opts.Services...,
	)
	mux.Handle(grpchealth.NewHandler(checker))
	mux.Handle(h.Path, h.Handler)

	compress1KB := connect.WithCompressMinBytes(1024)
	mux.Handle(grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(opts.ServiceName),
		compress1KB,
	))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(opts.ServiceName),
		compress1KB,
	))
}
