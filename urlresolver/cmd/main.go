package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DefaultFallbackUrl = "https://golang.org/"
)

type UrlResolverHandler struct {
	config UrlResolverConfig
}

type UrlResolverConfig struct {
	paths       map[string]interface{} `json:"paths"`
	fallbackUrl string                 `json:"-"`
}

func (c *UrlResolverConfig) GetLink(path string) string {
	url, ok := c.paths[path]
	if !ok {
		url = c.fallbackUrl
	}
	return url.(string)
}

func ParseUrlResolverConfig(configFile string) *UrlResolverConfig {
	fileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Config file not found: %v", err)
	}
	conf := &UrlResolverConfig{
		fallbackUrl: DefaultFallbackUrl,
	}
	err = json.Unmarshal(fileContent, &conf)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	return conf
}

func CreateUrlResolverHandlerFromFile(configFilePath string) *UrlResolverHandler {
	config := ParseUrlResolverConfig(configFilePath)
	return &UrlResolverHandler{
		config: *config,
	}
}

func CreateDefaultUrlResolverHandler() *UrlResolverHandler {
	return &UrlResolverHandler{
		UrlResolverConfig{
			paths: map[string]interface{}{
				"/go-http":    "https://golang.org/pkg/net/http",
				"/go-gophers": "https://github.com/shalakhin/gophericons/blob/master/preview.jpg",
			},
			fallbackUrl: DefaultFallbackUrl,
		},
	}
}

func (b *UrlResolverHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	url := b.config.GetLink(path)
	http.Redirect(rw, req, url, http.StatusPermanentRedirect)
}

func main() {
	mux := http.NewServeMux()

	configPath := flag.String("f", "", "File path to config file")
	flag.Parse()
	var handler http.Handler
	if *configPath != "" {
		handler = CreateUrlResolverHandlerFromFile(*configPath)
	} else {
		handler = CreateDefaultUrlResolverHandler()
	}
	mux.Handle("/", handler)

	server := &http.Server{
		Handler:      mux,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server listen and server error: %v", err)
		}
	}()

	waitForShutdown(server)
}

func waitForShutdown(server *http.Server) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigint

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
