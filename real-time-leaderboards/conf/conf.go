package conf

import "fmt"

type Config struct {
	ServerConfig ServerConfig
	DBConfig DBConfig
	RedisConfig RedisConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DBConfig struct {
	Host        string
	Port        int
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type RedisConfig struct {
	Host string
	Port int
	Password string
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=postgres password=postgres dbname=postgres sslmode=disable", d.Host, d.Port)
}
