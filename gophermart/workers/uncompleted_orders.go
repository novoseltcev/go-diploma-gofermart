package workers

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/adapters"
	"github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
)

type OrderStorager interface {
	GetUncompletedOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrder(ctx context.Context, number, status string, accural *models.Money) error
	UpdateBalance(ctx context.Context, userID models.UserID, accural models.Money) error
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type repository struct {
	db *sqlx.DB
}

func NewOrderStorager(db *sqlx.DB) OrderStorager {
	return &repository{db: db}	
}


func (r *repository) GetUncompletedOrders(ctx context.Context) (result []models.Order, err error) {
	err = r.db.SelectContext(ctx, &result, "SELECT number, status, accrual::numeric as accrual, user_id, uploaded_at FROM orders WHERE status = 'NEW' or status = 'PROCESSING' ORDER BY uploaded_at ASC")
	return result, err
}

func (r *repository) UpdateOrder(ctx context.Context, number string, status string, accrual *models.Money) error {
	_, err := r.db.ExecContext(ctx, "UPDATE orders SET status = $1, accrual = $2::MONEY WHERE number = $3", status, accrual, number)
	return err
}

func (r *repository) UpdateBalance(ctx context.Context, userID models.UserID, accural models.Money) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET balance = balance + $1::MONEY WHERE id = $2", accural, userID)
	return err
}

func (r *repository) Begin(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "BEGIN")
	return err
}

func (r *repository) Commit(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "COMMIT")
	return err
}

func (r *repository) Rollback(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "ROLLBACK")
	return err
}


func ProcessUncompletedOrders(ctx context.Context, storager OrderStorager, accuralAPI *adapters.AccrualAPI) {
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

				log.WithFields(log.Fields{"order": o.Number, "status": order.Status, "accural": order.Accrual}).Debug("got order")

				_ = storager.Begin(ctx)
				if err := storager.UpdateOrder(ctx, o.Number, order.Status, order.Accrual); err != nil {
					log.Error(err.Error())
					_ = storager.Rollback(ctx)
					continue
				}
				if order.Status == "PROCESSED" {
					if err := storager.UpdateBalance(ctx, o.UserID, *order.Accrual); err != nil {
						log.Error(err.Error())
						_ = storager.Rollback(ctx)
						continue
					}
				}
				_ = storager.Commit(ctx)
			}
		}
	}
}
