package worker

import (
	"context"

	"github.com/vgarvardt/gue/v5"
	"github.com/vgarvardt/gue/v5/adapter/libpq"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/task"
)

var Client *gue.Client

// Task list
const (
	TaskDeleteAccount = "user:delete"
	TaskRunCheck      = "check:run"
)

type Worker struct {
	runCheckTask *task.RunCheckTask
}

func NewWorker(runCheckTask *task.RunCheckTask) *Worker {
	return &Worker{runCheckTask: runCheckTask}
}

func (w *Worker) Start(ctx context.Context) error {
	tracingID := pkg.GetTracingID(ctx)
	poolAdapter := libpq.NewConnPool(infra.DB)

	Client, err := gue.NewClient(poolAdapter)
	if err != nil {
		return err
	}

	workMap := gue.WorkMap{
		TaskRunCheck: w.runCheckTask.Do,
	}

	// create a pool of workers
	workers, err := gue.NewWorkerPool(Client, workMap, config.App.WorkerPool)
	if err != nil {
		return err
	}

	// work jobs in goroutine
	go func() {
		err := workers.Run(ctx)
		if err != nil {
			panic(err)
		}
	}()
	lgr.Default.Print(tracingID, "worker started with", config.App.WorkerPool, "worker pool")
	return nil
}
