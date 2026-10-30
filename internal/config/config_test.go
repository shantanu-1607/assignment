package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("ADDR", "")
	t.Setenv("REQUEST_TIMEOUT", "")
	t.Setenv("RATE_PER_MINUTE", "")
	c, err := Load()
	if err != nil || c.Addr != "127.0.0.1:8080" || c.Rate != 60 {
		t.Fatalf("%+v %v", c, err)
	}
	for _, tc := range []struct{ key, value string }{{"ADDR", "invalid"}, {"REQUEST_TIMEOUT", "0s"}, {"REQUEST_TIMEOUT", "31s"}, {"REQUEST_TIMEOUT", "no"}, {"RATE_PER_MINUTE", "0"}, {"RATE_PER_MINUTE", "NaN"}, {"RATE_PER_MINUTE", "10001"}} {
		t.Run(tc.key+tc.value, func(t *testing.T) {
			t.Setenv(tc.key, tc.value)
			if _, err := Load(); err == nil {
				t.Fatal("invalid configuration accepted")
			}
		})
	}
	t.Setenv("ADDR", "0.0.0.0:9000")
	t.Setenv("REQUEST_TIMEOUT", "1s")
	t.Setenv("RATE_PER_MINUTE", "10")
	c, err = Load()
	if err != nil || c.Rate != 10 {
		t.Fatal(c, err)
	}
}
