package logging

import (
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

type tcpHook struct {
	level     int
	formatter logrus.Formatter
	address   string
	writer    io.Writer
}

func TCPHook(level int, formatter logrus.Formatter, address string) (logrus.Hook, error) {
	w, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &tcpHook{
		level:     level,
		formatter: formatter,
		address:   address,
		writer:    w,
	}, nil
}

func (h *tcpHook) Fire(entry *logrus.Entry) error {
	bytes, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	h.writer.Write(bytes)
	return nil
}

func (h *tcpHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, 0)
	for _, lvl := range logrus.AllLevels {
		if lvl <= logrus.Level(h.level) {
			levels = append(levels, lvl)
		}
	}
	return levels
}
