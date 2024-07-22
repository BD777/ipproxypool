package config

type Config struct {
	Mode string // "debug", "release"
	Host string // default "0.0.0.0"
	Port int    // default 9002
}

var DefaultConfig = &Config{
	Mode: "release",
	Host: "0.0.0.0",
	Port: 9002,
}
