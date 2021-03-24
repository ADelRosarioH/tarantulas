package main

import "github.com/adelrosarioh/tarantulas/collectors/agricultural"

func main() {

	if err := agricultural.Run(); err != nil {
		panic(err)
	}
}
