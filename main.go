package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/claymoret/flumebeat/beater"
)

func main() {
	err := beat.Run("flumebeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
