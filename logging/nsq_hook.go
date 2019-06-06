package logging

import (
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

type nsqHook struct {
	level     int
	formatter logrus.Formatter
	address   string
	producer  *nsq.Producer
	topic     string
}

func NSQHook(level int, formatter logrus.Formatter, address, topic string) (logrus.Hook, error) {
	producer, err := nsq.NewProducer(address, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	return &nsqHook{
		level:     level,
		formatter: formatter,
		address:   address,
		producer:  producer,
		topic:     topic,
	}, nil
}

func (h *nsqHook) Fire(entry *logrus.Entry) error {
	bytes, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	return h.producer.Publish(h.topic, bytes)
}

func (h *nsqHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, 0)
	for _, lvl := range logrus.AllLevels {
		if lvl <= logrus.Level(h.level) {
			levels = append(levels, lvl)
		}
	}
	return levels
}
