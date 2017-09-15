package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {

	//init logger
	var logger log.Logger
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))

	//init logging middleware with logger
	log := NewLogger(logger)

	mux := http.NewServeMux()
	mux.Handle("/hc", Action(healthcheck))

	mux.Handle("/api/shipmentEnvironmentsByBarge", log(Action(shipmentEnvironmentsByBarge)))

	mux.Handle("/", log(http.FileServer(http.Dir("freeboard"))))

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	fmt.Printf("server listening on port %s\n", port)
	http.ListenAndServe(":"+port, mux)
}

func healthcheck(r *http.Request) *Response {
	return DataJSON(http.StatusOK, "I'm healthy", nil)
}
