package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/srgyrn/lucky-38/pkg/config"
	"github.com/srgyrn/lucky-38/pkg/creating"
	"github.com/srgyrn/lucky-38/pkg/drawing"
	"github.com/srgyrn/lucky-38/pkg/listing"
	"github.com/srgyrn/lucky-38/pkg/rest"
	"github.com/srgyrn/lucky-38/pkg/storage"
)

func main() {
	conf, err := config.Load(".")
	if err != nil {
		log.Fatal(err.Error())
	}

	repository, err := storage.NewRepository(conf.Driver, conf.Source)
	if err != nil {
		log.Fatal(err.Error())
	}

	router := rest.Handler(creating.NewService(repository), listing.NewService(repository), drawing.NewService(repository))

	fmt.Printf("Your digital croupier is now available at: localhost:3000\n")
	log.Fatal(http.ListenAndServe(":3000", router))
}
