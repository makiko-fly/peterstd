package logging

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	global *logrus.Entry
)

func GetLogger() *logrus.Entry {
	return global
}

type LoggingConfig struct {
	AppName   string       `yaml:"app_name"`
	Enable    bool         `yaml:"enable"`
	Level     int          `yaml:"level"`
	Formatter string       `yaml:"formatter"`
	Hooks     []HookConfig `yaml:"hooks"`
}

func (c LoggingConfig) Init() error {
	logger := logrus.New()
	logger.SetLevel(logrus.Level(c.Level))
	logger.Formatter = getLogrusFormatter(c.Formatter)
	global = logrus.NewEntry(logger)

	if c.AppName != "" {
		global = logger.WithField("app", c.AppName)
	}

	if !c.Enable {
		logger.Out = nil
	}

	for _, config := range c.Hooks {
		if !config.Enable {
			continue
		}
		hook, err := config.Build()
		if err != nil {
			return err
		}
		logger.AddHook(hook)
	}
	return nil
}

type HookConfig struct {
	Name      string `yaml:"name"`
	Enable    bool   `yaml:"enable"`
	Level     int    `yaml:"level"`
	Formatter string `yaml:"formatter"`

	// FileHook config
	Filepath string `yaml:"filepath"`

	// TCPHook config
	TCPAddress string `yaml:"tcp_address"`

	// DailyRotatingHook config
	RotateDir string `yaml:"rotate_dir"`
	Filename  string `yaml:"filename"`
	MaxBackup int    `yaml:"max_backup"`

	// NSQHook config
	NSQDAddress string `yaml:"nsqd_address"`
	Topic       string `yaml:"topic"`
}

func (c HookConfig) Build() (logrus.Hook, error) {
	switch strings.ToLower(c.Name) {
	case "filehook":
		return c.BuildFileHook()
	case "tcphook":
		return c.BuildTCPHook()
	case "dailyrotatinghook":
		return c.BuildDailyRotatingHook()
	case "nsqhook":
		return c.BuildNSQHook()
	default:
		return nil, errors.New("not available hook")
	}
}

func getLogrusFormatter(formatter string) logrus.Formatter {
	switch formatter {
	case "text":
		return &logrus.TextFormatter{}
	case "json":
		return &logrus.JSONFormatter{}
	case "raw":
		return &RawFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}

func (c HookConfig) BuildFileHook() (logrus.Hook, error) {
	return FileHook(
		c.Level,
		getLogrusFormatter(c.Formatter),
		c.Filepath,
	)
}

func (c HookConfig) BuildTCPHook() (logrus.Hook, error) {
	return TCPHook(
		c.Level,
		getLogrusFormatter(c.Formatter),
		c.TCPAddress,
	)
}

func (c HookConfig) BuildDailyRotatingHook() (logrus.Hook, error) {
	return DailyRotatingHook(
		c.Level,
		getLogrusFormatter(c.Formatter),
		c.RotateDir,
		c.Filename,
		c.MaxBackup,
	), nil
}

func (c HookConfig) BuildNSQHook() (logrus.Hook, error) {
	return NSQHook(
		c.Level,
		getLogrusFormatter(c.Formatter),
		c.NSQDAddress,
		c.Topic,
	)
}
