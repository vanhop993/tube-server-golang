package common

import (
	"github.com/lib/pq"
	"reflect"
	"strings"
)

const (
	Postgre   = "postgre"
	Cassandra = "cassandra"
)

func ArrayValueForScansSQL(structValue interface{}, fieldColumnDB []string, dbName string) []interface{} {
	s := reflect.ValueOf(structValue).Elem()
	t := reflect.TypeOf(structValue).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, len(fieldColumnDB))
	for i := 0; i < len(fieldColumnDB); i++ {
		for j := 0; j < numCols; j++ {
			tagJson := strings.Split(t.Field(j).Tag.Get("json"), ",")
			if fieldColumnDB[i] == strings.ToLower(tagJson[0]) {
				field := s.Field(j)
				switch dbName {
				case Postgre:
					if reflect.TypeOf(field.Addr().Interface()).Elem().Kind() == reflect.Slice || reflect.TypeOf(field.Addr().Interface()).Elem().Kind() == reflect.Array {
						columns[i] = pq.Array(field.Addr().Interface())
					} else {
						columns[i] = field.Addr().Interface()
					}
					break
				case Cassandra:
					columns[i] = field.Addr().Interface()
					break
				}
				break
			}
		}
	}
	return columns
}
