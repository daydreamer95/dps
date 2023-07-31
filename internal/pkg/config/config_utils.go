package config

import (
	"dps/logger"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func loadFromEnv(prefix string, conf interface{}) {
	rv := reflect.ValueOf(conf)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		logger.Fatal(fmt.Sprintf("should be a valid pointer, but %s found", reflect.TypeOf(conf).Name()))
	}
	loadFromEnvInner(prefix, rv.Elem(), "")
}

func loadFromEnvInner(prefix string, conf reflect.Value, defaultValue string) {
	kind := conf.Kind()
	switch kind {
	case reflect.Struct:
		t := conf.Type()
		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag
			loadFromEnvInner(prefix+"_"+tag.Get("yaml"), conf.Field(i), tag.Get("default"))
		}
	case reflect.String:
		str := os.Getenv(toUnderscoreUpper(prefix))
		if str == "" {
			str = defaultValue
		}
		conf.Set(reflect.ValueOf(str))
	case reflect.Int64:
		str := os.Getenv(toUnderscoreUpper(prefix))
		if str == "" {
			str = defaultValue
		}
		if str == "" {
			str = "0"
		}
		conf.Set(reflect.ValueOf(int64(mustAtoi(str))))
	default:
		panic(fmt.Errorf("unsupported type: %s", conf.Type().Name()))
	}
}

func toUnderscoreUpper(key string) string {
	key = strings.Trim(key, "_")
	matchLastCap := regexp.MustCompile("([A-Z])([A-Z][a-z])")
	s2 := matchLastCap.ReplaceAllString(key, "${1}_${2}")

	matchFirstCap := regexp.MustCompile("([a-z])([A-Z]+)")
	s2 = matchFirstCap.ReplaceAllString(s2, "${1}_${2}")
	// logger.Infof("loading from env: %s", strings.ToUpper(s2))
	return strings.ToUpper(s2)
}

func checkConfig(conf *Type) error {
	switch conf.Store.Driver {
	case Mysql, Postgres:
		if conf.Store.Host == "" {
			return errors.New("Db host not valid ")
		}
		if conf.Store.Port == 0 {
			return errors.New("Db port not valid ")
		}
		if conf.Store.User == "" {
			return errors.New("Db user not valid ")
		}
	}
	return nil
}

// MustAtoi is string to int
func mustAtoi(s string) int {
	r, err := strconv.Atoi(s)
	if err != nil {
		logger.Error("convert to int error" + s)
	}
	return r
}
