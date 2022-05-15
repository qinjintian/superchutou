package main

import (
	"flag"
	"github.com/qinjintian/superchutou/app/service"
	"go.uber.org/dig"
)

type Application struct {
	Container *dig.Container
	cfgPath   string
}

var app = &Application{}

func init() {
	flag.StringVar(&app.cfgPath, "c", "configs/config.yaml", "config file path")

	flag.Parse()
}

func main() {
	app.Container = dig.New()

	err := bootstrap(app)
	if err != nil {
		panic(err)
	}

	err = app.Container.Invoke(func(svc *service.Service) {
		svc.Run()
	})

	if err != nil {
		panic(err)
	}
}
