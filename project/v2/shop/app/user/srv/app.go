package user

import (
	"fmt"
	"go-learn/project/v2/shop/app/user/srv/config"
	"go-learn/project/v2/shop/pkg/app"
)

func NewApp(basename string) *app.App {
	cfg := config.New()
	appl := app.NewApp("order", "shop",
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
		app.WithNoConfig(),
	)

	return appl
}

func run(cfg *config.Config) app.RunFunc {
	return func(basename string) error {
		fmt.Println(cfg.Log.Level)
		return nil
	}
}
