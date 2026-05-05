package worker

import (
	"context"
	"fmt"
	"time"

	"cryplio/internal/domain/trading"
	"cryplio/pkg/config"

	"github.com/hibiken/asynq"
)

const (
	TaskTradeReconcile   = "trade:reconcile"
	TaskTradeAutoDispute = "trade:auto_dispute"
)

type Worker struct {
	server       *asynq.Server
	tradeService trading.TradeService
}

func NewWorker(cfg *config.Config, tradeService trading.TradeService) *Worker {
	redisConn := asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}

	server := asynq.NewServer(
		redisConn,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &Worker{
		server:       server,
		tradeService: tradeService,
	}
}

func (w *Worker) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskTradeReconcile, w.handleTradeReconcile)
	mux.HandleFunc(TaskTradeAutoDispute, w.handleTradeAutoDispute)

	if err := w.server.Run(mux); err != nil {
		return fmt.Errorf("failed to start asynq server: %w", err)
	}
	return nil
}

func (w *Worker) handleTradeReconcile(ctx context.Context, t *asynq.Task) error {
	count, err := w.tradeService.ReconcileExpiredTrades(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		fmt.Printf("Reconciled %d expired trades\n", count)
	}
	return nil
}

func (w *Worker) handleTradeAutoDispute(ctx context.Context, t *asynq.Task) error {
	// 2 hour grace period for auto-dispute
	count, err := w.tradeService.FlagAutoDisputesForOverduePaidTrades(ctx, 2*time.Hour)
	if err != nil {
		return err
	}
	if count > 0 {
		fmt.Printf("Flagged %d trades for auto-dispute\n", count)
	}
	return nil
}

// Scheduler sets up periodic tasks
type Scheduler struct {
	inspector *asynq.Scheduler
}

func NewScheduler(cfg *config.Config) *Scheduler {
	redisConn := asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	return &Scheduler{
		inspector: asynq.NewScheduler(redisConn, &asynq.SchedulerOpts{}),
	}
}

func (s *Scheduler) Start() error {
	// Every 5 minutes
	_, err := s.inspector.Register("@every 5m", asynq.NewTask(TaskTradeReconcile, nil))
	if err != nil {
		return err
	}

	// Every 15 minutes
	_, err = s.inspector.Register("@every 15m", asynq.NewTask(TaskTradeAutoDispute, nil))
	if err != nil {
		return err
	}

	return s.inspector.Run()
}
