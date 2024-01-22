package config

type Config struct {
	BaseUrl string
	Port    string
}

func NewConfig() *Config {
	return &Config{
		BaseUrl: "http://localhost:9991",
		Port:    "9991",
	}
}
