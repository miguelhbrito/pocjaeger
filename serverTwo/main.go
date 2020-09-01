package main

import (
	"github.com/pocjaeger/pkg/serverTwo"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	closer := tracing.InitJaeger("server-two-tracing")
	defer closer.Close()

	http.Handle("/server-two", http.HandlerFunc(serverTwo.MyTracingHandlerServerTwo))

	log.Info().Msg("Server two listening on 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err)
	}
}
