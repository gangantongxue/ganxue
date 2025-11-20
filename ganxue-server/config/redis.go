package config

type RedisConf struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db"`
}
