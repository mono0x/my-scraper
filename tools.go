//go:build tools

package main

import (
	_ "github.com/air-verse/air"
	_ "github.com/go-task/task/v3/cmd/task"
	_ "gotest.tools/gotestsum"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
