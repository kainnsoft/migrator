package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kainnsoft/migrator/internal/entity"
	"github.com/kainnsoft/migrator/internal/repository"
)

type (
	IPGUsecase interface {
		UseCaseGetCountPayments(context.Context) (int, error)
		UseCaseInsertPayment(ctx context.Context, payment *entity.Payment) (rowID int64, err error)
		ConsumePayment(data []byte)
		ClosePGUseCase()
	}
	pgUsecase struct {
		Repo       repository.IPGRepo
		workerpool int
		jobs       chan *entity.Payment
		res        chan string
		done       chan struct{}
	}
)

func NewPGUsecase(workerpool int, repo repository.IPGRepo) IPGUsecase {
	r := &pgUsecase{
		Repo:       repo,
		workerpool: workerpool,
		jobs:       make(chan *entity.Payment),
		res:        make(chan string),
		done:       make(chan struct{}),
	}
	go r.workersInit()

	return r
}

func (u *pgUsecase) UseCaseGetCountPayments(ctx context.Context) (int, error) {
	return u.Repo.GetCountPayments(ctx)
}

func (u *pgUsecase) UseCaseInsertPayment(ctx context.Context, payment *entity.Payment) (rowID int64, err error) {
	return u.Repo.InsertPayment(ctx, payment)
}

func (u *pgUsecase) ConsumePayment(data []byte) {
	var paymant = &entity.Payment{}
	if err := json.Unmarshal(data, paymant); err != nil {
		fmt.Println("ConsumePayment unmarshalling error:", err)
	}
	// Отправка заданий в поток. Синхронно.
	u.jobs <- paymant
}

func (u *pgUsecase) ClosePGUseCase() {
	u.done <- struct{}{}
}
