package cron

type Command struct {
}

func (c Command) Type() string {
	return TypeCommand
}
