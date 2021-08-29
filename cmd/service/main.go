package main

import (
	"github.com/anantadwi13/letsencrypt-manager/internal"
)

func main() {
	config := internal.NewConfig(internal.ConfigParam{
		PublicStaticPath: "./public",
		ApiPort:          5555,
	})

	s := internal.NewService(config)
	s.Start()
}
