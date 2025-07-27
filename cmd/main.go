// Package main implements the main entry point for the application.
package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"

	"gogs.utking.net/utking/spaces/cmd/tasks"
)

func main() {
	// Initialize the tasks command
	tasks.InitTasks()
	// Execute the tasks command
	tasks.Execute()
}
