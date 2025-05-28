package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	TGToken       string
	RouteWatching int
	ChaiID        int64
}

func (c *Config) Init() error {
	token := os.Getenv("WATCHER_TG_TOKEN")
	if token == "" {
		return errors.New("WATCHER_TG_TOKEN environment variable is required")
	}

	routeStr := os.Getenv("WATCHER_ROUTE_WATCHING")
	if routeStr == "" {
		return errors.New("WATCHER_ROUTE_WATCHING environment variable is required")
	}

	route, err := strconv.Atoi(routeStr)
	if err != nil {
		return fmt.Errorf("WATCHER_ROUTE_WATCHING environment variable must be a number: %w", err)
	}

	chatStr := os.Getenv("WATCHER_CHAT_ID")
	if chatStr == "" {
		return errors.New("WATCHER_CHAT_ID environment variable is required")
	}

	chatID, err := strconv.ParseInt(chatStr, 10, 64)
	if err != nil {
		return fmt.Errorf("WATCHER_CHAT_ID environment variable must be a number: %w", err)
	}

	c.TGToken = token
	c.RouteWatching = route
	c.ChaiID = chatID

	return nil
}
