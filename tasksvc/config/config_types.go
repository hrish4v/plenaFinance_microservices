package config

import "time"

type SvcBase struct {
	TaskSvc string `mapstructure:"usersvc"`
}

type StartupConfig struct {
	Debug   bool          `mapstructure:"debug"`
	Stage   string        `mapstructure:"stage"`
	Server  Server        `mapstructure:"server"`
	Timeout time.Duration `mapstructure:"timeout"`
	SvcBase SvcBase       `mapstructure:"svcbase"`
	TaskDB  DBConfig      `mapstructure:"taskDB"`
}

type DBConfig struct {
	Type           string `mapstructure:"type"`
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"name"`
	MaxConnections int    `mapstructure:"maxConnections"`
	Database       string `mapstructure:"database"`
}

type Server struct {
	Port string `mapstructure:"address"`
}
