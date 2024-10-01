package workers

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/adapters"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/domains/orders"
)

func ProcessUncompletedOrders(ctx context.Context, storager orders.OrderStorager, accuralAPI *adapters.AccuralAPI) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			orders, err := storager.GetUncompletedOrders(ctx)
			if err != nil {
				log.Error("failed to get orders")
				continue
			}
			log.WithField("orders", len(orders)).Info("got orders")

			for _, o := range orders {
				order, err := adapters.Retry(func ()(*adapters.Order, error) {
					return accuralAPI.GetOrderAccuralStatus(ctx, o.Number)
				}, 3)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				if order == nil {
					continue
				}

				err = storager.Update(ctx, o.Number, order.Status, order.Accural)
				if err != nil {
					log.Error(err.Error())
				}
			}
		}
	}
}
