package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	sendgrid "terraform-provider-sendgrid/internal/provider"
)

func main() {

	var debug bool

	opts := providerserver.ServeOpts{
		Address: "hashicorp.com/edu/sendgrid",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), sendgrid.New("0.0.1"), opts)

	if err != nil {
		log.Fatal(err)
	}
}
