package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	mw "restproject/internal/api/middlewares"
	"restproject/internal/api/router"
	"restproject/pkg/utils"
	"time"

	"golang.org/x/net/http2"
)

func main() {

	const port = 8080

	// load the TLS cert and key
	const cert = "server.cert"
	const key = "server.key"

	// configure TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	rl := mw.NewRateLimiter(5, time.Minute)

	router := router.Router()
	securityMux := utils.ApplyMiddlewares(router, mw.Comression, mw.SecurityHeaders, rl.RateLimit, mw.Cors)

	// create a custom server
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   securityMux,
		TLSConfig: tlsConfig,
	}

	// enable http2
	if err := http2.ConfigureServer(server, &http2.Server{}); err != nil {
		log.Fatalf("Failed to configure HTTP/2: %v", err)
	}

	fmt.Println("Server started on port:", port)

	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("error starting server", err)
	}

}
