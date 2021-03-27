package utils

import (
	"os"

	"github.com/gocolly/redisstorage"
)

func Storage(prefix string) *redisstorage.Storage {

	endpoint := os.Getenv("REDIS_ENDPOINT")
	password := os.Getenv("REDIS_ENDPOINT")

	storage := &redisstorage.Storage{
		Address:  endpoint,
		Password: password,
		Prefix:   prefix,
	}

	return storage
}
