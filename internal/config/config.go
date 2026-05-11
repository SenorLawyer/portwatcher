package config

import (
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Interval         time.Duration
	Docker           bool
	HistoryPath      string
	HistoryRetention time.Duration
}

func Default() Config {
	return Config{
		Interval:         time.Second,
		Docker:           true,
		HistoryPath:      defaultHistoryPath(),
		HistoryRetention: 7 * 24 * time.Hour,
	}
}

func defaultHistoryPath() string {
	dir, err := os.UserCacheDir()
	if err != nil || dir == "" {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "portwatch", "history.jsonl")
}
