package datasource

import (
	"github.com/nsqio/go-nsq"
	"gitlab.wallstcn.com/spider/peterstd/util"
	"io/ioutil"
	"log"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// NSQ Producer

type NSQProducerConfig struct {
	NSQDAddress string `yaml:"nsqd_address"`
}

func (c NSQProducerConfig) NewNSQProducer() *nsq.Producer {
	config := nsq.NewConfig()
	config.WriteTimeout = time.Second * 10
	p, err := nsq.NewProducer(c.NSQDAddress, config)
	if err != nil {
		panic(util.AppendLineNumToErr(err))
	}
	if err := p.Ping(); err != nil {
		// ping failed
		panic(util.AppendLineNumToErr(err))
	}

	return p
}

// ---------------------------------------------------------------------------------------------------------------------
// NSQ Consumer

type NSQConsumer struct {
	consumer *nsq.Consumer
	NSQConsumerConfig
}

type NSQConsumerConfig struct {
	Topic         string   `yaml:"topic"`
	Channel       string   `yaml:"channel"`
	NSQDAddress   []string `yaml:"nsqd_address"`
	LookupAddress []string `yaml:"lookup_address"`
}

func NewNSQConsumer(c NSQConsumerConfig, handlers ...nsq.Handler) *NSQConsumer {
	return c.NewNSQConsumer(handlers...)
}

func (c NSQConsumerConfig) NewNSQConsumer(handlers ...nsq.Handler) *NSQConsumer {
	config := nsq.NewConfig()
	config.MaxInFlight = 1000
	config.LookupdPollInterval = time.Second
	consumer, err := nsq.NewConsumer(c.Topic, c.Channel, config)
	logger := log.New(ioutil.Discard, "", 0)
	consumer.SetLogger(logger, 0)
	if err != nil {
		panic(err)
	}

	co := &NSQConsumer{
		consumer:          consumer,
		NSQConsumerConfig: c,
	}

	for _, handler := range handlers {
		co.consumer.AddHandler(handler)
	}

	return co
}

func (c *NSQConsumer) AddHandler(handler nsq.Handler) {
	c.consumer.AddHandler(handler)
}

func (c *NSQConsumer) AddConcurrencyHandler(handler nsq.Handler, concurrency int) {
	c.consumer.AddConcurrentHandlers(handler, concurrency)
}

func (c *NSQConsumer) ServeForever(errorChan chan<- error) {
	if cap(errorChan) == 0 {
		panic("Capacity of error channel should > 0")
	}

	var err error
	if len(c.LookupAddress) > 0 {
		err = c.consumer.ConnectToNSQLookupds(c.LookupAddress)
	} else {
		err = c.consumer.ConnectToNSQDs(c.NSQDAddress)
	}

	if err != nil {
		errorChan <- err
		return
	}
}
