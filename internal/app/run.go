package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/kainnsoft/migrator/pkg/httpserver"
	"github.com/kainnsoft/migrator/pkg/mysql"
	"github.com/kainnsoft/migrator/pkg/postgres"
	"github.com/kainnsoft/migrator/pkg/rmq"

	"github.com/kainnsoft/migrator/config"
	"github.com/kainnsoft/migrator/internal/controller/http/v1"
	"github.com/kainnsoft/migrator/internal/entity"
	"github.com/kainnsoft/migrator/internal/handler"
	"github.com/kainnsoft/migrator/internal/repository"
	"github.com/kainnsoft/migrator/internal/usecase"

	"github.com/fasthttp/router"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	var (
		logger       = zap.Must(zap.NewProduction())
		mySqlDB      mysql.IMySQL
		pgDB         postgres.IPgDB
		rmqChannel   rmq.IRmq
		mySQLRepo    repository.IMySQLRepo
		pgRepo       repository.IPGRepo
		rmqRepo      repository.IRmqRepo
		rmqConsumer  handler.IConsumer
		mySQLUsecase usecase.IMySQLUsecase
		pgUsecase    usecase.IPGUsecase
		httpHandlers v1.IHttpHandlers
		mux          = router.New()
		httpServer   *httpserver.Server
		err          error
	)

	// db mySQL
	if mySqlDB, err = mysql.New(&cfg.MySql); err != nil {
		logger.Fatal("app - Run - mySql db open:", zap.Error(err))
	}
	// db postgres
	if pgDB, err = postgres.New(&cfg.PG); err != nil {
		logger.Fatal("app - Run - can't create PG connection:", zap.Error(err))
	}
	if err = pgDB.MigrationUp(); err != nil {
		logger.Fatal("pg migration error:", zap.Error(err))
	}
	// rmq
	if rmqChannel, err = rmq.NewRmqConn(&cfg.RMQ); err != nil {
		logger.Fatal("app - Run - can't create RMQ connection:", zap.Error(err))
	}
	if err = rmqInit(rmqChannel); err != nil {
		logger.Fatal("app - Run - can't init RMQ :", zap.Error(err))
	}

	// repo mySQL
	mySQLRepo = repository.NewMySQLRepo(mySqlDB.DB(), logger)
	// repo postgres
	pgRepo = repository.NewPGRepo(pgDB.DB(), logger)
	// repo rmq publisher
	rmqRepo = repository.NewRmqRepo(rmqChannel.Rmq(), logger)

	// usecase mySQL
	mySQLUsecase = usecase.NewMySQLUsecase(mySQLRepo, rmqRepo)
	// usecase postgres
	pgUsecase = usecase.NewPGUsecase(cfg.WorkerCount, pgRepo)

	// handlers usecases
	httpHandlers = handler.NewHandler(mySQLUsecase, pgUsecase, logger)
	// rmq consumer
	rmqConsumer = handler.NewListener(rmqChannel, entity.QueuePayments, pgUsecase)

	// http
	v1.NewRouter(mux, httpHandlers, logger)
	httpServer = httpserver.New(mux, cfg.HTTP)
	if httpServer != nil {
		logger.Info("app - Run - httpServer has run on addr", zap.String("", httpServer.GetAddr()))
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify: %w", zap.Error(err))
	}

	// Shutdown
	rmqConsumer.CloseConsumer()
	logger.Info("app - Run - rmqConsumer.CloseConsumer(): OK")

	if err = mySqlDB.Close(); err != nil {
		logger.Error("app - Run - mySql.Close():", zap.Error(err))
	}
	logger.Info("app - Run - mySql.Close(): OK")

	if err = pgDB.Close(); err != nil {
		logger.Error("app - Run - pgDB.Close():", zap.Error(err))
	}
	logger.Info("app - Run - pgDB.Close(): OK")

	if err = rmqChannel.CloseRmqChan(); err != nil {
		logger.Error("app - Run - rmqChannel.CloseRmqChan():", zap.Error(err))
	}
	logger.Info("app - Run - rmqChannel.CloseRmqChan(): OK")

	err = httpServer.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown:", zap.Error(err))
	}

	_ = logger.Sync()
}

func rmqInit(rmqChannel rmq.IRmq) (err error) {
	if err = rmqChannel.ExchangeDeclare(entity.Exchange, entity.ExchangeType); err != nil {
		return fmt.Errorf("rmqChannel.ExchangeDeclare error: %v", err)
	}

	var q amqp.Queue
	if q, err = rmqChannel.QueueDeclare(entity.QueuePayments); err != nil {
		return fmt.Errorf("rmqChannel.QueueDeclare error: %v", err)
	}

	if err = rmqChannel.QueueBind(q, entity.RoutingKeyPayments, entity.Exchange); err != nil {
		return fmt.Errorf("rmqChannel.QueueBind error: %v", err)
	}

	return nil
}
