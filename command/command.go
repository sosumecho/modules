package command

import (
	"github.com/sosumecho/modules/utils"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var CMD *Command

type Command struct {
	command *cobra.Command
}

func (c *Command) AddCommand(cmd ...Commander) *Command {
	for _, item := range cmd {
		c.command.AddCommand(item.Command())
	}
	return c
}

func (c *Command) Run() error {
	return c.command.Execute()
}

func (c *Command) SetDefault(cmd Commander) *Command {
	com, _, err := c.command.Find(os.Args[1:])
	firstArg := utils.FirstElement(os.Args[1:])
	if err == nil && com.Use == c.command.Use && firstArg != "-h" && firstArg != "--help" {
		args := append([]string{cmd.Command().Use}, os.Args[1:]...)
		c.command.SetArgs(args)
	}
	return c
}

func (c *Command) Bootstrap(fn func(command *cobra.Command, args []string)) *Command {
	c.command.PersistentPreRun = fn
	return c
}

func (c *Command) AddDependency(fn func(*Command)) *Command {
	fn(c)
	return c
}

func init() {
	name := path.Base(os.Args[0])
	cmd := &cobra.Command{
		Use:   name,
		Short: name,
		Long:  name,
	}
	CMD = &Command{
		command: cmd,
	}
}
