package main

import (
	"fmt"
	"os"

	quantum "github.com/cbhcbhcbh/Quantum/internal/apiserver"
)

func main() {
	command := quantum.NewQuantumCommand()

	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
