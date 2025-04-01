//go:build tools

package main

import (
	_ "github.com/air-verse/air"
	_ "github.com/go-task/task/v3/cmd/task"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
