package main

import (
	"github.com/codeshelldev/goredacted/cmd"
	_ "github.com/codeshelldev/goredacted/internals/parser/parsers"
)

func main() {
	cmd.Execute()
}