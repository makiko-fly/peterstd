package logging

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type fileHook struct {
	level     int
	formatter logrus.Formatter
	filepath  string
	writer    io.Writer
}

func FileHook(level int, formatter logrus.Formatter, filepath string) (logrus.Hook, error) {
	w, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666|os.ModeSticky)
	if err != nil {
		return nil, err
	}
	return &fileHook{
		level:     level,
		formatter: formatter,
		filepath:  filepath,
		writer:    w,
	}, nil
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	bytes, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	h.writer.Write(bytes)
	return nil
}

func (h *fileHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, 0)
	for _, lvl := range logrus.AllLevels {
		if lvl <= logrus.Level(h.level) {
			levels = append(levels, lvl)
		}
	}
	return levels
}
