package domain

import (
	"context"
	"time"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type JobDomain struct{}

func NewJobDomain() *JobDomain {
	return &JobDomain{}
}

func (j *JobDomain) ListJobsToRun(ctx context.Context, from, to int) ([]model.Job, error) {
	now := times.Now()
	prev := now.Add(time.Second * time.Duration(from))
	later := now.Add(time.Second * time.Duration(to))

	condition := Job.NextRunAt.GT(TimestampT(prev)).
		AND(Job.NextRunAt.LT(TimestampT(later))).
		AND(Job.LastRanAt.LT(TimestampT(prev)))

	condition = condition.OR(Job.LastRanAt.IS_NULL())
	condition = condition.AND(Job.On.EQ(Bool(true)))

	stmt := SELECT(Job.AllColumns).FROM(Job).WHERE(condition)

	var jobs []model.Job
	err := stmt.QueryContext(ctx, infra.DB, &jobs)
	return jobs, err
}

func (j *JobDomain) ListRecurringJobs(ctx context.Context, from, to int) ([]model.Job, error) {
	stmt := SELECT(Job.AllColumns).FROM(Job).WHERE(Job.On.EQ(Bool(true)).AND(Job.Recurring.EQ(Bool(true))))

	var jobs []model.Job
	err := stmt.QueryContext(ctx, infra.DB, &jobs)
	return jobs, err
}

func (j *JobDomain) UpdateRunning(
	ctx context.Context,
	id int64,
	lastRunAt, nextRunAt *time.Time,
	status resource.JobStatus,
) (*model.Job, error) {
	if !status.Valid() {
		return nil, constant.ErrInvalidJobStatus
	}
	statusValue := int32(status)

	now := times.Now()
	job := &model.Job{
		Status:    &statusValue,
		LastRanAt: lastRunAt,
		NextRunAt: nextRunAt,
		UpdatedAt: now,
	}

	updateStmt := Job.UPDATE(Job.Status, Job.LastRanAt, Job.NextRunAt, Job.UpdatedAt).
		MODEL(job).WHERE(Job.ID.EQ(Int(id))).RETURNING(Job.AllColumns)

	err := updateStmt.QueryContext(ctx, infra.DB, job)
	return job, err
}

func (j *JobDomain) UpdateStatus(
	ctx context.Context,
	id int64,
	status resource.JobStatus,
) (*model.Job, error) {
	if !status.Valid() {
		return nil, constant.ErrInvalidJobStatus
	}
	statusValue := int32(status)

	now := times.Now()
	job := &model.Job{
		Status:    &statusValue,
		UpdatedAt: now,
	}

	updateStmt := Job.UPDATE(Job.Status, Job.UpdatedAt).MODEL(job).WHERE(Job.ID.EQ(Int(id))).RETURNING(Job.AllColumns)

	err := updateStmt.QueryContext(ctx, infra.DB, job)
	return job, err
}
