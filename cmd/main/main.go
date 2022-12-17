package main

import (
	"fmt"
	"os"
	"x/pkg/args"
	"x/pkg/config"
)

func main() {
	arguments, err := args.ParseArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg, err := config.New(cwd, arguments.ConfigFiles)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", cfg)
	fmt.Printf("%#v\n", arguments)
}
