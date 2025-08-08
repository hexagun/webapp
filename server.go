package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	viper.SetDefault("environment", "default")
	log.Info().Msg("configs:")
	log.Info().Msg(fmt.Sprintf("%s%s", "environment:", viper.GetString("environment")))
	log.Info().Msg(fmt.Sprintf("%s%d", "webapp.port:", viper.GetInt("webapp.port")))
}

func ServeReactApp(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("route /")
	// Get the absolute path to the dist directory
	absPath, err := filepath.Abs("frontend/dist")
	if err != nil {
		http.Error(w, "Could not find dist directory", http.StatusInternalServerError)
		return
	}

	// Serve the file based on URL path
	path := filepath.Join(absPath, r.URL.Path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If file doesn't exist, serve index.html for client-side routing
		http.ServeFile(w, r, filepath.Join(absPath, "index.html"))
		return
	}

	// Serve the static file
	http.FileServer(http.Dir(absPath)).ServeHTTP(w, r)

}

func main() {
	initLogging()

	setConfig()

	// Serve static files from the dist directory
	http.HandleFunc("/*", ServeReactApp)

	port := fmt.Sprintf("%s%d", ":", viper.GetInt("webapp.port"))
	log.Info().Msg("init server.")
	log.Fatal().Err(http.ListenAndServe(port, nil))
	log.Info().Msg("shutting out server.")
}
