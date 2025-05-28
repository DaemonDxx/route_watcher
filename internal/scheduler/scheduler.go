package scheduler

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

type Scheduler struct {
	c *cron.Cron
}

func New() *Scheduler {
	c := cron.New()
	return &Scheduler{
		c: c,
	}
}

func (s *Scheduler) Start(spec string, fn func(ctx context.Context) error) error {
	_, err := s.c.AddFunc(spec, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := fn(ctx); err != nil {
			log.Println(fmt.Sprintf("error executing task: %s", err.Error()))
		}
	})

	if err != nil {
		return fmt.Errorf("add task failed: %s", err.Error())
	}

	done := make(chan bool)

	s.c.Start()

	<-done

	return nil

}
func (s *Scheduler) Stop() {}
