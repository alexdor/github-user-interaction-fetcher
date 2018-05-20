package main

import (
	"fmt"
	"os"

	"github.com/alexdor/github-user-interaction-fetcher/cmd"
)

func main() {
	if err := cmd.Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
