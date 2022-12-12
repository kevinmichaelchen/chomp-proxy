package app

import (
	"github.com/bufbuild/connect-go"
	modService "github.com/kevinmichaelchen/chomp-proxy/internal/app/service"
	"github.com/kevinmichaelchen/chomp-proxy/internal/service"
	modConnect "github.com/kevinmichaelchen/chomp-proxy/pkg/fxmod/connect"
	"github.com/kevinmichaelchen/chomp-proxy/pkg/fxmod/logging"
	"go.buf.build/bufbuild/connect-go/kevinmichaelchen/chompapis/chomp/v1beta1/chompv1beta1connect"
	"go.uber.org/fx"
)

var Module = fx.Options(
	modConnect.CreateModule(&modConnect.ModuleOptions{
		HandlerProvider: func(svc *service.Service) modConnect.HandlerOutput {
			compress1KB := connect.WithCompressMinBytes(1024)
			// Register our Connect-Go server
			path, h := chompv1beta1connect.NewChompServiceHandler(
				svc,
				compress1KB,
			)
			return modConnect.HandlerOutput{
				Path:    path,
				Handler: h,
			}
		},
		ServiceName: chompv1beta1connect.ChompServiceName,
		Services: []string{
			"chompv1beta1.ChompService",
		},
	}),
	logging.Module,
	modService.Module,
)
