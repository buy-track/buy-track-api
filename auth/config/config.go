package config

type App struct {
	GrpcServerAddress string   `mapstructure:"grpc_server_address" default:"localhost:50051"`
	Broker            Broker   `mapstructure:"broker"`
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
	Name     string `mapstructure:"name" default:"auth_table"`
	User     string `mapstructure:"username" default:"root"`
}

type Cache struct {
	Driver   string `mapstructure:"driver" default:"redis"`
	Password string `mapstructure:"password" default:""`
	Host     string `mapstructure:"host" default:"127.0.0.1"`
	Port     string `mapstructure:"port" default:"6379"`
	//these are for redis
	DB     int    `mapstructure:"db" default:"0"`
	Prefix string `mapstructure:"prefix" default:"auth_cache_prefix"`
}

type Broker struct {
	Driver   string `mapstructure:"driver" default:"redis"`
	Password string `mapstructure:"password" default:""`
	Host     string `mapstructure:"host" default:"127.0.0.1"`
	Port     string `mapstructure:"port" default:"6379"`
	DB       int    `mapstructure:"db" default:"0"`
}

type Redis struct {
	Password string `mapstructure:"password" default:""`
	Host     string `mapstructure:"host" default:"127.0.0.1"`
	Port     string `mapstructure:"port" default:"6379"`
	DB       int    `mapstructure:"db" default:"2"`
}
