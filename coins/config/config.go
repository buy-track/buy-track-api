package config

type App struct {
	GrpcServerAddress string   `mapstructure:"grpc_server_address" default:"localhost:50053"`
	RedisDB           Redis    `mapstructure:"redis"`
	Database          Database `mapstructure:"database"`
	Cache             Cache    `mapstructure:"cache"`
	LogLevel          string   `mapstructure:"log_level" default:"info"`
	Version           string   `mapstructure:"version" default:"0.0.1"`
}

type Database struct {
	Driver   string `mapstructure:"driver" default:"postgres"`
	Password string `mapstructure:"password" default:""`
	Host     string `mapstructure:"host" default:"127.0.0.1"`
	Port     string `mapstructure:"port" default:"5432"`
	Name     string `mapstructure:"name" default:"coin_table"`
	User     string `mapstructure:"username" default:"root"`
}

type Cache struct {
	Driver   string `mapstructure:"driver" default:"redis"`
	Password string `mapstructure:"password" default:""`
	Host     string `mapstructure:"host" default:"127.0.0.1"`
	Port     string `mapstructure:"port" default:"6379"`
	//these are for redis
	DB     int    `mapstructure:"db" default:"0"`
	Prefix string `mapstructure:"prefix" default:"coin_cache_prefix"`
}

type Redis struct {
	Password string `mapstructure:"password" default:""`
	Host     string `mapstructure:"host" default:"127.0.0.1"`
	Port     string `mapstructure:"port" default:"6379"`
	DB       int    `mapstructure:"db" default:"1"`
	Prefix   string `mapstructure:"prefix" default:"coins_service_db_"`
}
