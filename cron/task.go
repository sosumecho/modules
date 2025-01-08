package cron

type Task struct{}

func (Task) Type() string {
	return TypeTask
}
