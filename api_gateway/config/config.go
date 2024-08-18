package config

type App struct {
	HttpServerAddress      string `mapstructure:"http_server" default:"localhost:3001"`
	Cache                  Cache  `mapstructure:"cache"`
	LogLevel               string `mapstructure:"log_level" default:"info"`
	Version                string `mapstructure:"version" default:"0.0.1"`
	AuthServiceGrpcAddress string `mapstructure:"auth_grpc_server_address" default:"localhost:50051"`
	AuthBroker             Broker `mapstructure:"auth_broker"`
	UserServiceGrpcAddress string `mapstructure:"user_grpc_server_address" default:"localhost:50052"`
	UserBroker             Broker `mapstructure:"user_broker"`
	CoinServiceGrpcAddress string `mapstructure:"coin_grpc_server_address" default:"localhost:50053"`
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
