package main

import (
	"net/http"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/transport"
)

func main() {
	hnd := transport.NewGatewayHandler()

	srv, err := api.NewServer(hnd)
	if err != nil {
		panic(err)
	}

	if err = http.ListenAndServe(":8080", srv); err != nil {
		panic(err)
	}
}
