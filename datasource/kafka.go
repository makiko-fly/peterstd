package datasource

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"gitlab.wallstcn.com/spider/peterstd/util"
)

type KafkaConsumerConfig struct {
	BrokerList []string `yaml:"broker_list"`
	Topic      string   `yaml:"topic"`
	Group      string   `yaml:"group"`
}

func (c KafkaConsumerConfig) NewKafkaConsumer(topic string) *KafkaConsumer {
	client, err := sarama.NewClient(c.BrokerList, nil)
	if err != nil {
		panic(util.AppendLineNumToErr(err))
	}
	consumerGroup, err := sarama.NewConsumerGroupFromClient(c.Group, client)
	if err != nil {
		panic(util.AppendLineNumToErr(err))
	}
	return &KafkaConsumer{Topic: topic, Client: client, ConsumerGroup: consumerGroup, Handler: nil}

}

// =====================================================================================================================

type KafkaConsumer struct {
	Topic string
	// store client and consumer group so we can close them in the event of a error or a crash
	Client        sarama.Client
	ConsumerGroup sarama.ConsumerGroup
	Handler       sarama.ConsumerGroupHandler
}

func (c KafkaConsumer) ServeForever(errChan chan<- error) {
	if c.Handler == nil {
		errChan <- util.AppendLineNumToErr(errors.New("KafkaConsumerGroup's Handler is null"))
		return
	}
	go func() {
		for {
			err := c.ConsumerGroup.Consume(context.Background(), []string{c.Topic}, c.Handler)
			if err != nil {
				errChan <- util.AppendLineNumToErr(err)
				// clean
				c.ConsumerGroup.Close()
				c.Client.Close()
				return
			}
		}
	}()
}

// only one handler
func (c KafkaConsumer) AddHandler(handler sarama.ConsumerGroupHandler) {
	c.Handler = handler
}
