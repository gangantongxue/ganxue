package config

type MongoConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password" env:"MONGO_PASSWORD"`
	Database string `yaml:"database"`
}
