package usecase

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/kainnsoft/migrator/internal/entity"
)

var (
	// Рабочий поток.
	worker = func(ctx context.Context,
		u *pgUsecase,
		workerID int,
		jobs <-chan *entity.Payment,
		results chan<- string,
	) {
		fmt.Printf("Start worker №%d\n", workerID)
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Stop worker №%d\n", workerID)
				return
			case job := <-jobs:
				var (
					id  int64
					err error
				)
				if id, err = u.Repo.InsertPaymentWithConflictResolwing(ctx, job); err != nil {
					results <- err.Error()
					return
				}
				results <- strconv.Itoa(int(id))
			}
		}
	}
)

func (u *pgUsecase) workersInit() {
	var (
		// Максимальное количество рабочих равно количеству ядер в системе.
		W      = u.workerpool
		ctx    context.Context
		cancel context.CancelFunc
	)
	if W == 0 {
		W = runtime.NumCPU()
	}
	// Запуск рабочих в отдельных потоках. Асинхронно.
	ctx, cancel = context.WithCancel(context.Background())
	for i := 0; i < W; i++ {
		go worker(ctx, u, i, u.jobs, u.res)
	}

	// Поток обработки результатов. Асинхронно.
	go func(ctx context.Context, ch chan string) {
		for {
			select {
			case <-ctx.Done():
				return
			case val := <-ch:
				fmt.Printf("insert payment result is: %v\n", val)
			}
		}
	}(ctx, u.res)

	<-u.done
	cancel()
	time.Sleep(time.Second)
	close(u.jobs)
}
