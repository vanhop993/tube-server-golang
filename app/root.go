package app

import (
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
)

type Root struct {
	Server     ServerConfig    `mapstructure:"server"`
	Mongo      MongoConfig     `mapstructure:"mongo"`
	Cassandra  CassandraConfig `mapstructure:"cassandra"`
	Postgre    PostgreConfig   `mapstructure:"postgre"`
	OpenDb     int             `mapstructure:"openDb"`
	Log        log.Config      `mapstructure:"log"`
	MiddleWare mid.LogConfig   `mapstructure:"middleware"`
	Key        string          `mapstructure:"key"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port *int64 `mapstructure:"port"`
}

type MongoConfig struct {
	Uri      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type CassandraConfig struct {
	Uri      string `mapstructure:"uri" json:"uri,omitempty" gorm:"column:uri" bson:"uri,omitempty" dynamodbav:"uri,omitempty" firestore:"uri,omitempty"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type PostgreConfig struct {
	Driver         string `mapstructure:"driver"`
	DataSourceName string `mapstructure:"data_source_name"`
}
