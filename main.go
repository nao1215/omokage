package main

import (
	"os"

	"github.com/nao1215/dyer/cmd"
)

func main() {
	workDir, err := os.Getwd()
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	app := cmd.NewApp(os.Stdout, os.Stderr, workDir)
	os.Exit(app.Run(os.Args[1:]))
}
