package app

import (
	"github.com/Firoz01/go-mongodb-test/config"
	"github.com/Firoz01/go-mongodb-test/logger"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"github.com/Firoz01/go-mongodb-test/web"
	"sync"
)

type Application struct {
	wg sync.WaitGroup
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Init() {
	config.LoadConfig()
	conf := config.GetConfig()
	logger.SetupLogger(conf.ServiceName)
	mongodb.InitMongoDB()

}

func (app *Application) Run() {
	web.StartServer(&app.wg)
}

func (app *Application) Wait() {
	// wait for all the top level goroutines
	app.wg.Wait()
}

func (app *Application) Cleanup() {
	mongodb.Disconnect() // Graceful shutdown to disconnect MongoDB
}
