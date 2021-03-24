package main

import (
	hardware "github.com/adelrosarioh/tarantulas/collectors/hardware"
)

func main() {

	if err := hardware.Run(); err != nil {
		panic(err)
	}
}
