package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/Braendie/RestAPI/internal/config"
	"github.com/Braendie/RestAPI/internal/user"
	"github.com/Braendie/RestAPI/internal/user/db"
	"github.com/Braendie/RestAPI/pkg/client/mongodb"
	"github.com/Braendie/RestAPI/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

// main creates router, handler, logger, storage and it registers all the handlers by Register().
// it uses config from config.yml using GetConfig function.
// it also begins start function, which starts server.
func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	mongoDBClient, err := mongodb.NewClient(context.Background(), cfg.MongoDB.Host, cfg.MongoDB.Port, cfg.MongoDB.Username, cfg.MongoDB.Password, cfg.MongoDB.Database, cfg.MongoDB.AuthDB)
	if err != nil {
		panic(err)
	}
	
	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	logger.Info("create user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, logger, cfg)
}

// start begins listen and serve.
func start(router *httprouter.Router, logger *logging.Logger, cfg *config.Config) {
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Info("server is listening unix socket")
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Fatal(server.Serve(listener))
}
