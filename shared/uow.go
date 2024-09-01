package shared

import (
	"context"

	"github.com/jmoiron/sqlx"
)


type UnitOfWork struct {
	Tx *sqlx.Tx
	conn *sqlx.Conn
	apllied bool
}

func (uow *UnitOfWork) Apply() error {
	uow.apllied = true
	return uow.Tx.Commit()
}

func (uow *UnitOfWork) Discard() error {
	return uow.Tx.Rollback()
}

func (uow *UnitOfWork) Close() error {
	defer uow.conn.Close()
	if !uow.apllied {
		return uow.Discard()
	}
	return nil
}


type UOWPool = func(context.Context) *UnitOfWork
func NewUOWPool(db *sqlx.DB) UOWPool {
	return func(ctx context.Context) *UnitOfWork {
		conn, err := db.Connx(ctx)
		if err != nil {
			panic(err)
		}

		tx, err := conn.BeginTxx(ctx, nil)
		if err != nil {
			panic(err)
		}

		return &UnitOfWork{tx, conn, false}
	}
}
