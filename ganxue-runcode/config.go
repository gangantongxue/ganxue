package main

import "os"

type redisCfg struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

var RedisCfg = redisCfg{
	Addr:     "redis:6379",
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       0,
}
