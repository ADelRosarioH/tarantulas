package main

import "github.com/adelrosarioh/tarantulas/collectors/dairy"

func main() {

	if err := dairy.Run(); err != nil {
		panic(err)
	}
}
