package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/docker/docker/client"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func main() {
	// Create a docker client from the environment configured docker host.
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	app := kingpin.New("swarm-api", "A simple HTTP server to communicate with Docker Swarm")

	// We don't actually care about the cli version at all, use docker client.
	app.Version(cli.ClientVersion())

	server := app.Command("serve", "Start the API server.").Alias("server")
	host := server.Flag("host", "http host interface").Envar("HOST").Default("localhost").String()
	port := server.Flag("port", "http listening port").Envar("PORT").Default("8080").String()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case server.FullCommand():
		StartAPIServer(cli, net.JoinHostPort(*host, *port))
	}
}

func StartAPIServer(cli *client.Client, listenAddr string) {
	// Initialize http routing muxer
	r := chi.NewRouter()

	// Base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Handler method to retrieve a swarm join token from the docker client.
	r.Get("/swarm/token/{tokenType:(manager|worker)}", func(w http.ResponseWriter, r *http.Request) {
		swarm, err := cli.SwarmInspect(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var token string
		switch chi.URLParam(r, "tokenType") {
		case "manager":
			token = swarm.JoinTokens.Manager
		case "worker":
			token = swarm.JoinTokens.Worker
		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		render.Respond(w, r, map[string]string{
			"token": token,
		})
	})

	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	log.Printf("http server listening at http://%s\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
