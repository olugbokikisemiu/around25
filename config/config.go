package config

type Config struct {
	OrderTTL   int64  `envconfig:"LOCATION_HISTORY_TTL_SECONDS"`
	ServerPort string `envconfig:"HISTORY_SERVER_LISTEN_ADDR" default:"8080"`
}
