package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func initLogging() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()
}

func setConfig() {
	viper.SetConfigName("config")  // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/config") // path to look for the config file in
	viper.AddConfigPath(".")       // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Error().Msg(fmt.Sprintf("fatal error config file: %w", err))
		} else {
			// Config file was found but another error was produced
			log.Error().Msg(fmt.Sprintf("fatal error config file: %w", err))
		}
	}

	viper.SetDefault("webapp.port", 8888)
	log.Info().Msg("configs:")
	log.Info().Msg(fmt.Sprintf("%s%d", "webapp.port:", viper.GetInt("webapp.port")))
}

func main() {
	initLogging()

	setConfig()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	port := fmt.Sprintf("%s%d", ":", viper.GetInt("webapp.port"))
	log.Info().Msg("init server.")
	log.Fatal().Err(http.ListenAndServe(port, nil))
	log.Info().Msg("shutting out server.")
}
