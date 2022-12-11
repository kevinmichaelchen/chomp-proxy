package main

import (
	"github.com/kevinmichaelchen/chomp-proxy/internal/app"
	"go.uber.org/fx"
)

func main() {
	a := fx.New(
		app.Module,
		// TODO configure logrus
		// For details, see https://github.com/uber-go/fx/blob/master/fxevent/zap.go
	)
	a.Run()
}
