package scheduler

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/kavehjamshidi/arvan-challenge/service/quota_reset/contract"
	"os"
	"os/signal"
	"time"
)

func Schedule(qoutaResetService contract.QuotaResetService) {
	s := gocron.NewScheduler(time.Local)
	s.SingletonModeAll()

	s.Every(1).Minute().Do(func() {
		qoutaResetService.ResetUserQuota(context.Background())
	})

	s.StartAsync()
	defer s.Stop()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
