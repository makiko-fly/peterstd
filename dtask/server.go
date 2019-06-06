package dtask

import (
	"fmt"

	"github.com/RichardKnop/machinery/v1"
	redisbackend "github.com/RichardKnop/machinery/v1/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v1/brokers/redis"
	"github.com/RichardKnop/machinery/v1/config"
)

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Auth string `yaml:"auth"`
	DB   int    `yaml:"db"`
}

type Config struct {
	DefaultQueue    string      `yaml:"default_queue"`
	ResultsExpireIn int         `yaml:"results_expire_in"`
	WorkerName      string      `yaml:"worker_name"`
	Redis           RedisConfig `yaml:"redis"`
}

func (c Config) NewServer() *machinery.Server {
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

func (c Config) NewWorker() *Worker {
	server := c.NewServer()
	return &Worker{
		server: server,
		Name:   c.WorkerName,
	}
}
