package peterstd

import (
	"fmt"

	"github.com/RichardKnop/machinery/v1"
	redisbackend "github.com/RichardKnop/machinery/v1/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v1/brokers/redis"
	"github.com/RichardKnop/machinery/v1/config"
)

type DistaskConfig struct {
	DefaultQueue    string      `yaml:"default_queue" envconfig:"DEFAULT_QUEUE"`
	ResultsExpireIn int         `yaml:"results_expire_in" envconfig:"RESULTS_EXPIRE_IN"`
	WorkerName      string      `yaml:"worker_name"`
	Redis           RedisConfig `yaml:"redis"`
}

func (c DistaskConfig) NewServer() *machinery.Server {
	cnf := &config.Config{
		Broker:          "amqp://guest:guest@localhost:5672/",
		DefaultQueue:    c.DefaultQueue,
		ResultBackend:   "amqp://guest:guest@localhost:5672/",
		ResultsExpireIn: c.ResultsExpireIn,
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			DelayedTasksPollPeriod: 20,
		},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		panic(err)
	}

	redisURL := fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
	broker := redisbroker.New(cnf, redisURL, c.Redis.Auth, "", c.Redis.DB)
	backend := redisbackend.New(cnf, redisURL, c.Redis.Auth, "", c.Redis.DB)

	server.SetBroker(broker)
	server.SetBackend(backend)

	return server
}

func (c DistaskConfig) NewWorker() *DistaskWorker {
	server := c.NewServer()
	return &DistaskWorker{
		server: server,
		Name:   c.WorkerName,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
type DistaskWorker struct {
	server *machinery.Server
	Name   string
}

func (w *DistaskWorker) Register(name string, task interface{}) error {
	return w.server.RegisterTask(name, task)
}

func (w *DistaskWorker) ServeForever(errorChan chan<- error) {
	if cap(errorChan) == 0 {
		panic("Capacity of error channel shoule > 0")
	}
	w.server.NewWorker(w.Name, 0).LaunchAsync(errorChan)
}
