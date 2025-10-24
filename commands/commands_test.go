package commands

import (
	"log/slog"
	"os"
	"testing"
)

func init() {
	defaultLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(defaultLogger)
}

func TestMain(m *testing.M) {
	m.Run()
}
