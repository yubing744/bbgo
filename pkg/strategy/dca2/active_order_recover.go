package dca2

import (
	"context"
	"time"

	"github.com/c9s/bbgo/pkg/exchange/retry"
	"github.com/c9s/bbgo/pkg/strategy/common"
	"github.com/c9s/bbgo/pkg/util"
)

func (s *Strategy) recoverPeriodically(ctx context.Context) {
	s.logger.Info("[DCA] monitor and recover periodically")
	interval := util.MillisecondsJitter(10*time.Minute, 5*60*1000)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.recoverActiveOrders(ctx); err != nil {
				s.logger.WithError(err).Warn(err, "failed to recover active orders")
			}
		}
	}
}

func (s *Strategy) recoverActiveOrders(ctx context.Context) error {
	openOrders, err := retry.QueryOpenOrdersUntilSuccessfulLite(ctx, s.ExchangeSession.Exchange, s.Symbol)
	if err != nil {
		s.logger.WithError(err).Warn("failed to query open orders")
		return err
	}

	activeOrders := s.OrderExecutor.ActiveMakerOrders().Orders()

	// update num of open orders metrics
	if metricsNumOfOpenOrders != nil {
		metricsNumOfOpenOrders.With(baseLabels).Set(float64(len(openOrders)))
	}

	// update num of active orders metrics
	if metricsNumOfActiveOrders != nil {
		metricsNumOfActiveOrders.With(baseLabels).Set(float64(len(activeOrders)))
	}

	opts := common.SyncActiveOrdersOpts{
		Logger:          s.logger,
		Exchange:        s.ExchangeSession.Exchange,
		ActiveOrderBook: s.OrderExecutor.ActiveMakerOrders(),
		OpenOrders:      openOrders,
	}

	return common.SyncActiveOrders(ctx, opts)
}
