package worker

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/vgarvardt/gue/v5"
	"github.com/vgarvardt/gue/v5/adapter"
	"github.com/vgarvardt/gue/v5/adapter/libpq"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/task"
)

var (
	SlowWheel       *gue.Client
	FastWheel       *asynq.Client
	fastWheelServer *asynq.Server
)

// Task list
const (
	TaskDeleteAccount = "user:delete"
	TaskRunCheck      = "check:run"
	TaskStartMonitor  = "monitor:start"
)

type Worker struct {
	runCheckTask     *task.RunCheckTask
	startMonitorTask *task.StartMonitorTask
}

func NewWorker(runCheckTask *task.RunCheckTask, startMonitorTask *task.StartMonitorTask) *Worker {
	return &Worker{runCheckTask: runCheckTask, startMonitorTask: startMonitorTask}
}

func (w *Worker) StartGue(ctx context.Context) error {
	tracingID := pkg.GetTracingID(ctx)
	poolAdapter := libpq.NewConnPool(infra.DB)

	var err error
	SlowWheel, err = gue.NewClient(poolAdapter)
	if err != nil {
		return err
	}

	workMap := gue.WorkMap{
		TaskRunCheck:     w.runCheckTask.Do,
		TaskStartMonitor: w.startMonitorTask.Do,
	}

	// create a pool of workers
	workers, err := gue.NewWorkerPool(
		SlowWheel, workMap,
		config.App.WorkerPool,
		gue.WithPoolLogger(adapter.NewStdLogger()),
		gue.WithPoolPollInterval(500*time.Millisecond),
		gue.WithPoolPollStrategy(gue.RunAtPollStrategy),
	)
	if err != nil {
		return err
	}

	// work jobs in goroutine
	go func() {
		if err := workers.Run(ctx); err != nil {
			panic(err)
		}
	}()
	lgr.Print(tracingID, "slow worker started with", config.App.WorkerPool, "worker pool")
	return nil
}

func (w *Worker) StartAsynq(ctx context.Context) error {
	tracingID := pkg.GetTracingID(ctx)

	redisClientOpt := asynq.RedisClientOpt{
		Addr: config.App.RedisQueue, Username: config.App.RedisQueueUser, Password: config.App.RedisQueuePass,
	}
	FastWheel = asynq.NewClient(redisClientOpt)
	fastWheelServer = asynq.NewServer(redisClientOpt, asynq.Config{Concurrency: config.App.WorkerPool})

	mux := asynq.NewServeMux()
	mux.Handle(TaskStartMonitor, w.startMonitorTask)
	mux.Handle(TaskRunCheck, w.runCheckTask)

	go func() {
		if err := fastWheelServer.Run(mux); err != nil {
			panic(err)
		}
	}()
	lgr.Print(tracingID, "fast worker started with", config.App.WorkerPool, "worker pool")
	return nil
}

func Shutdown() {
	fastWheelServer.Shutdown()
	if err := FastWheel.Close(); err != nil {
		sentry.CaptureException(err)
	}
}
