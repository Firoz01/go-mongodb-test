package web

import (
	"fmt"
	"github.com/Firoz01/go-mongodb-test/config"
	"github.com/Firoz01/go-mongodb-test/web/middlewares"
	"log/slog"
	"net/http"
	"sync"
)

func StartServer(wg *sync.WaitGroup) {
	manager := middlewares.NewManager()
	mux := http.NewServeMux()

	InitRoutes(mux, manager)

	handler := middlewares.EnableCors(mux)

	wg.Add(1)

	go func() {
		defer wg.Done()

		conf := config.GetConfig()

		addr := fmt.Sprintf(":%d", conf.HttpPort)
		slog.Info(fmt.Sprintf("Listening at %s", addr))

		if err := http.ListenAndServe(addr, handler); err != nil {
			slog.Error(err.Error())
		}
	}()
}
