//+build tools

package main

import (
	_ "github.com/lestrrat/go-server-starter/cmd/start_server"
	_ "github.com/mattn/goveralls"
	_ "honnef.co/go/tools/cmd/megacheck"
)
