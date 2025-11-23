package main

import (
	"fmt"
	"os"

	"github.com/Shresth72/ark/internal/jepsen"
	"github.com/alecthomas/kong"
)

type EchoCmd struct{}
type UniqueIdsCmd struct{}

type Cli struct {
	Echo      EchoCmd      `cmd:"" help:"Test Echo"`
	UniqueIds UniqueIdsCmd `cmd:"" help:"Test UniqueIds"`
}

func main() {
	cli := &Cli{}
	ctx := kong.Parse(cli)

	handlers := map[string]jepsen.Node{
		"echo":       &jepsen.EchoNode{},
		"unique-ids": &jepsen.UniqueIdsNode{},
	}

	node, ok := handlers[ctx.Command()]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", ctx.Command())
		os.Exit(1)
	}

	fmt.Printf("Running %s...\n", ctx.Command())
	if err := jepsen.MainLoop(node, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
