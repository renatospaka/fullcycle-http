package main

import (
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/renatospaka/customer-processor-test/pkg/http"
)

func main() {
	// Setting timezone to America/Sao_Paulo
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Warn().Msgf("error when setting timezone to 'America/Sao_Paulo': %s", err.Error())
		time.Now().In(time.UTC)
	} else {
		time.Now().In(loc)
	}

	// Start HTTP server
	httpServer := http.New(
		http.WithAddress("0.0.0.0"),
		http.WithPort(8080),
		http.WithVersion("0.1.0"),
	)
	if httpServer == nil {
		log.Fatal().Msg("abnormal termination: http server not initialized")
	}
	defer httpServer.Close()

	connected := make(chan bool)
	defer close(connected)

	// Serve HTTP
	if _, err = httpServer.Serve(connected); err != nil {
		log.Fatal().Msgf("abnormal termination: %s", err.Error())
	}
	if <-connected {
		// params.HTTP = server
		log.Info().Msgf("serving at %s", httpServer.Addr())
	}

	for {
		select {}
	}
}
