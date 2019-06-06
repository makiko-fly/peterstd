package logging

import (
	"io"

	"github.com/sirupsen/logrus"
	"gitlab.wallstcn.com/spider/peterstd/writers"
)

type dailyRotatingHook struct {
	level     int
	formatter logrus.Formatter

	writer io.Writer
	date   string
}

func DailyRotatingHook(level int, formatter logrus.Formatter, dir, format string, maxBackup int) logrus.Hook {
	hook := &dailyRotatingHook{
		level:     level,
		formatter: formatter,
		writer:    writers.NewDailyRotatingWriter(dir, format, maxBackup),
	}
	return hook
}

func (h *dailyRotatingHook) Fire(entry *logrus.Entry) error {
	bytes, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(bytes)
	return err
}

func (h *dailyRotatingHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, 0)
	for _, lvl := range logrus.AllLevels {
		if lvl <= logrus.Level(h.level) {
			levels = append(levels, lvl)
		}
	}
	return levels
}
