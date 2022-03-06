package main

import (
	"bacon/cmd"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	cmd.Execute()
}
