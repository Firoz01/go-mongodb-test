package app

import (
	"context"
	"github.com/Firoz01/go-mongodb-test/config"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"github.com/Firoz01/go-mongodb-test/web"
	"log"
	"sync"
	"time"
)

type Application struct {
	wg sync.WaitGroup
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Init() {
	config.LoadConfig()
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mongodb.GetClient(ctx, cfg.MongodbURL, cfg.MongodbDBName)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB client: %v", err)
	}

}

func (app *Application) Run() {
	web.StartServer(&app.wg)
}

func (app *Application) Wait() {
	// wait for all the top level goroutines
	app.wg.Wait()
}

func (app *Application) Cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongodb.Disconnect(ctx) // Graceful shutdown to disconnect MongoDB
}
