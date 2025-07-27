// Package tasks provides a command line interface for running various tasks
package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/utking/spaces/cmd/tasks/migrate"
	"github.com/utking/spaces/cmd/tasks/server"
)

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Task runner",
	Long:  `Task runner to start serve HTTP or run other CLI commands`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Usage()
		}
	},
}

func InitTasks() {
	server.Init(rootCmd)
	migrate.Init(rootCmd)
}

// Execute runs the root command and handles errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
