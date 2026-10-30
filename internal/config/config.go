package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr    string
	Timeout time.Duration
	Rate    int
}

func Load() (Config, error) {
	c := Config{Addr: "127.0.0.1:8080", Timeout: 2 * time.Second, Rate: 60}
	if v := os.Getenv("ADDR"); v != "" {
		c.Addr = v
	}
	if _, _, err := net.SplitHostPort(c.Addr); err != nil {
		return c, fmt.Errorf("ADDR must be host:port: %w", err)
	}
	if v := os.Getenv("REQUEST_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d < time.Millisecond || d > 30*time.Second {
			return c, fmt.Errorf("REQUEST_TIMEOUT must be 1ms..30s")
		}
		c.Timeout = d
	}
	if v := os.Getenv("RATE_PER_MINUTE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 10000 {
			return c, fmt.Errorf("RATE_PER_MINUTE must be 1..10000")
		}
		c.Rate = n
	}
	return c, nil
}
