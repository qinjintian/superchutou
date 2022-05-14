package main

import (
	"github.com/qinjintian/superchutou/app/core"
	"github.com/qinjintian/superchutou/app/service"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func bootstrap(app *Application) error {
	// 配置加载
	err := app.Container.Provide(func() (*service.Config, error) {
		config, err := ioutil.ReadFile(app.cfgPath)
		if err != nil {
			return nil, err
		}

		cfg := &service.Config{}

		if err = yaml.Unmarshal(config, cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	})

	// 湘潭大学接口服务
	err = app.Container.Provide(func() (*core.XTDXService, error) {
		return core.NewXTDXService()
	})

	// 业务处理
	err = app.Container.Provide(func(cfg *service.Config, xd *core.XTDXService) (*service.Service, error) {
		return service.NewService(cfg, xd)
	})

	if err != nil {
		return nil
	}

	return nil
}