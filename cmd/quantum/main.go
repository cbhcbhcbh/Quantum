package main

import (
	"fmt"
	"os"

	"github.com/cbhcbhcbh/Quantum/internal/quantum"
)

func main() {
	command := quantum.NewQuantumCommand()

	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
