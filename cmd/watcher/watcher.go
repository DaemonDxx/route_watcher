package main

import (
	"context"
	"fmt"
	"log"
	"watcher/config"
	"watcher/internal/api"
	"watcher/internal/notifier"
	"watcher/internal/scheduler"
)

func main() {

	cfg := config.Config{}
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	c, err := api.New()
	if err != nil {
		log.Fatal(err)
	}

	tg, err := notifier.NewTelegram(cfg.TGToken)
	if err != nil {
		log.Fatalf(err.Error())
	}

	sch := scheduler.New()

	log.Printf("start watching for route %d by chat user %d \n\r", cfg.RouteWatching, cfg.ChaiID)
	if err := sch.Start("@every 10m", func(ctx context.Context) error {
		log.Println("Check available ticket...")
		ok, err := c.IsAvailableRoute(ctx, cfg.RouteWatching)
		if err != nil {
			tg.Notify(cfg.ChaiID, fmt.Sprintf("[ERROR] %s", err.Error()))
			return fmt.Errorf("error checking route: %s", err.Error())
		}

		if !ok {
			log.Println("route is not available")
			return nil
		}

		if err := tg.Notify(cfg.ChaiID, "Продажа билетов открыта!"); err != nil {
			return fmt.Errorf("error sending notification: %s", err.Error())
		}

		return nil

	}); err != nil {
		log.Fatal(err)
	}

}
