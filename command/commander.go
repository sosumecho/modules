package command

import "github.com/spf13/cobra"

type Commander interface {
	Command() *cobra.Command
	Run(cmd *cobra.Command, args []string)
	SubCommands() []Commander
}

type BaseCommand struct{}

func (b BaseCommand) Command() *cobra.Command {
	panic("implement me")
}

func (b BaseCommand) Run(cmd *cobra.Command, args []string) {
}

func (b BaseCommand) SubCommands() []Commander {
	return nil
}
