package main

import (
	"os"

	"github.com/arnie97/emu-log/tasks"
)

func main() {
	tasks.CmdParser(os.Args...)
}
