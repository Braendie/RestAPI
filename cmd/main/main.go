package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Braendie/RestAPI/internal/user"
	"github.com/julienschmidt/httprouter"
)


// main creates router, handler and it registers all the handlers by Register().
// it also begins start function, which starts server
func main() {
	fmt.Println("create router")
	router := httprouter.New()

	handler := user.NewHandler()
	handler.Register(router)

	start(router)
}

// start begins listen and serve
func start(router *httprouter.Router) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
