package main

import (
	"fmt"
	"os"
	"runtime"

	. "github.com/nttcom/ecl2mond"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error:%s\n", err)
			os.Exit(1)
		}
	}()

	os.Exit(_main())
}

func _main() int {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(1)
	}

	cli := NewEcl2mond()

	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}

	return 0
}
