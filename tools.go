//go:build tools

package main

import (
	_ "github.com/air-verse/air"
	_ "github.com/go-task/task/v3/cmd/task"
	_ "github.com/suzuki-shunsuke/pinact/cmd/pinact"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
