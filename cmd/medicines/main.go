package main

import "github.com/adelrosarioh/tarantulas/collectors/medicines"

func main() {

	if err := medicines.Run(); err != nil {
		panic(err)
	}
}
