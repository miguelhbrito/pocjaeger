package main

import (
	"github.com/pocjaeger/pkg/serverOne"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	closer := tracing.InitJaeger("server-one-tracing")
	defer closer.Close()

	//TODO maybe change this to use midtracing
	http.Handle("/", http.HandlerFunc(serverOne.MyTracingHandlerServerOne))
	http.Handle("/serverTwoResponse", http.HandlerFunc(serverOne.HandlerServerTwoResponse))

	log.Info().Msg("Server one listening on 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal().Err(err)
	}
}