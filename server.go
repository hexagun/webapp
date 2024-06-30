package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func setConfig() {
	viper.SetConfigName("config")  // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/config") // path to look for the config file in
	viper.AddConfigPath(".")       // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
	}

	viper.SetDefault("webapp.port", 8888)
}

func main() {
	setConfig()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	port := fmt.Sprintf("%s%d", ":", viper.GetInt("webapp.port"))
	log.Fatal(http.ListenAndServe(port, nil))
}
