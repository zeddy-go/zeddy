package config

import (
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/container"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	ModeLocal   = "local"
	ModeDevelop = "develop"
	ModeStaging = "staging"
	ModeRelease = "release"
)

func WithPath(path string) func(*Module) {
	return func(module *Module) {
		module.path = path
	}
}

func WithContent(content string) func(*Module) {
	return func(module *Module) {
		module.config = content
	}
}

func NewModule(opts ...func(*Module)) *Module {
	m := &Module{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

type Module struct {
	config string
	path   string
}

func (m Module) Name() string {
	return "config"
}

func (m Module) Init() (err error) {
	c := viper.New()
	c.SetConfigType("yaml")

	if m.path != "" {
		m.config = readConfig(m.path)
	}

	err = c.ReadConfig(strings.NewReader(m.config))
	if err != nil {
		return
	}
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	c.AutomaticEnv()
	err = container.Bind[*viper.Viper](c, container.AsSingleton())
	if err != nil {
		return
	}

	setLog(c)

	return
}

func setLog(c *viper.Viper) {
	opts := &slog.HandlerOptions{
		AddSource: true,
	}
	if c.GetString("mode") == ModeLocal || c.GetBool("showDebugLog") {
		opts.Level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
}

func readConfig(path string) (result string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return
	}

	return string(content)
}
