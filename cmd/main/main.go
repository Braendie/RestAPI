package main

import (
	"net"
	"net/http"
	"time"

	"github.com/Braendie/RestAPI/internal/user"
	"github.com/Braendie/RestAPI/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

// main creates router, handler, logger and it registers all the handlers by Register().
// it also begins start function, which starts server
func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	logger.Info("create user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, logger)
}

// start begins listen and serve
func start(router *httprouter.Router, logger logging.Logger) {
	logger.Info("start application")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("server is listening port 0.0.0.0:8080")
	logger.Fatal(server.Serve(listener))
}
