package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/anhvandev/doc-kit/internal/cli"
)

func main() {
	err := cli.Execute()
	if err == nil {
		return
	}
	var ee *cli.ExitError
	if errors.As(err, &ee) {
		if ee.Msg != "" {
			fmt.Fprintln(os.Stderr, ee.Msg)
		}
		os.Exit(ee.Code)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
