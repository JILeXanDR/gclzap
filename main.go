package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"logging/logging"
)

func main() {
	cfg, err := loadConfig("./etc/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %+v", err)
		return
	}

	credentialsJSON, err := os.ReadFile(cfg.GCL.ServiceAccountPath)
	if err != nil {
		log.Fatalf("failed to read file %s: %+v", cfg.GCL.ServiceAccountPath, err)
		return
	}

	logger, err := logging.New(logging.Options{
		Level: cfg.Logging.Level,
		GCL: logging.GoogleCloudLoggingOptions{
			CredentialsJSON: credentialsJSON,
			ProjectID:       cfg.GCL.ProjectID,
			LogID:           cfg.GCL.LogID,
			Level:           cfg.GCL.Level,
		},
	})
	if err != nil {
		log.Fatalf("failed to build logger: %+v", err)
		return
	}

	randomLogs(logger)
}

func randomLogs(logger *zap.Logger) {
	go updateCtx()

	storeCtx(genCtxWithRandReqID())

	for {
		ctx := loadCtx()

		func() {
			newCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			callFunc(newCtx)
		}()

		var wg sync.WaitGroup

		wg.Add(11)

		go func() {
			defer wg.Done()
			time.Sleep(1 * time.Second)
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger)
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger.Named("internal"))
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger.Named("internal").Named("v1"))
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger.Named("internal").Named("v2"))
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger.Named("http"))
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger.Named("db").Named("mysql"))
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger.Named("redis").Named("mysql"))
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger.Named("redis").Named("clickhouse"))
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger.Named("user").Named("repo"))
		}()

		go func() {
			defer wg.Done()
			randLevel(ctx, logger.Named("item").Named("repo"))
		}()

		wg.Wait()
	}
}

func callFunc(ctx context.Context) {
	select {
	case <-ctx.Done():

	}
}

var ctxAtomic = atomic.Value{}

func storeCtx(ctx context.Context) {
	ctxAtomic.Store(ctx)
}

func loadCtx() context.Context {
	ctx, ok := ctxAtomic.Load().(context.Context)
	if !ok {
		return context.Background()
	}
	return ctx
}

func genCtxWithRandReqID() context.Context {
	return context.WithValue(context.Background(), "reqid", uuid.New().String())
}

func updateCtx() {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	for {
		select {
		case _, ok := <-tick.C:
			if !ok {
				return
			}
			storeCtx(genCtxWithRandReqID())
		default:

		}
	}
}

var loggers = []func(*zap.Logger){
	func(logger *zap.Logger) {
		logger.Error("test error", zap.Error(errors.New("wrong id")))
	},
	func(logger *zap.Logger) {
		logger.Warn("wrong user id", zap.Int("userid", 3))
	},
	func(logger *zap.Logger) {
		logger.Error("failed to fetch URL",
			zap.Int("attempt", 3),
			zap.Duration("backoff", time.Second),
		)
	},
	func(logger *zap.Logger) {
		logger.Info("info")
	},
	func(logger *zap.Logger) {
		logger.Debug("received value", zap.Int("value", rand.Intn(999)))
	},
	func(logger *zap.Logger) {
		logger.Error("failed to ping a service", zap.String("service", "s1"), zap.Error(fmt.Errorf("pinging service: %w", errors.New("wrong status code 500"))))
	},
	func(logger *zap.Logger) {
		logger.Error("failed to ping a service", zap.String("service", "s1"), zap.Error(fmt.Errorf("pinging service: %w", errors.New("wrong status code 500"))))
	},
}

func randLevel(ctx context.Context, logger *zap.Logger) {
	reqID, _ := ctx.Value("reqid").(string)
	logger = logger.With(zap.String("requestId", reqID))
	loggers[rand.Intn(len(loggers))](logger)
}
