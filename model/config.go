package model

type Config struct {
	AlExtens    []string
	MaxFileTask int
	MaxParallel int
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		MaxFileTask: 3,
		MaxParallel: 3,

		AlExtens: []string{".pdf", ".jpeg", ".jpg"},
		Port:     ":8080",
	}
}
