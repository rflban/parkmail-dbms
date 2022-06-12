package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

const (
	defaultProfile = "local"
)

type Conf struct {
	Server struct {
		Port int
	}
	Database struct {
		Host          string
		Name          string
		Username      string
		Password      string
		Port          int
		MaxConns      int
		MinConns      int
		MaxIdleTimeNS time.Duration
	}
}

func defaultConf() Conf {
	conf := Conf{}

	conf.Server.Port = 8080

	conf.Database.Host = "localhost"
	conf.Database.Name = "localhost"
	conf.Database.Username = "postgres"
	conf.Database.Password = "postgres"
	conf.Database.Port = 5432
	conf.Database.MaxConns = 10
	conf.Database.MinConns = 2
	conf.Database.MaxIdleTimeNS = 60_000_000_000

	return conf
}

func getConfig(ctx context.Context) (*Conf, error) {
	conf := defaultConf()

	viper.SetEnvPrefix("PMDBMS_FORUM")
	viper.SetDefault("PROFILE", defaultProfile)

	if err := viper.BindEnv("PROFILE"); err != nil {
		return nil, err
	}

	profile, ok := viper.Get("PROFILE").(string)
	if !ok {
		return nil, fmt.Errorf("ENV PMDBMS_FORUM_PROFILE invalid value")
	}

	viper.SetConfigName("config." + profile)
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs/forum")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println(err.Error())
		} else {
			fmt.Println(err.Error())
		}
	} else {
		if serverConf, ok := viper.Get("server").(map[string]interface{}); ok {
			if port, ok := serverConf["port"].(int64); ok {
				conf.Server.Port = int(port)
			}
		}
		if databaseConf, ok := viper.Get("database").(map[string]interface{}); ok {
			if username, ok := databaseConf["username"].(string); ok {
				conf.Database.Username = username
			}
			if password, ok := databaseConf["password"].(string); ok {
				conf.Database.Password = password
			}
			if name, ok := databaseConf["name"].(string); ok {
				conf.Database.Name = name
			}
			if host, ok := databaseConf["host"].(string); ok {
				conf.Database.Host = host
			}
			if port, ok := databaseConf["port"].(int64); ok {
				conf.Database.Port = int(port)
			}
			if maxConns, ok := databaseConf["max_conns"].(int64); ok {
				conf.Database.MaxConns = int(maxConns)
			}
			if minConns, ok := databaseConf["min_conns"].(int64); ok {
				conf.Database.MinConns = int(minConns)
			}
			if maxIdleTimeNS, ok := databaseConf["max_idle_time_ns"].(int64); ok {
				conf.Database.MaxIdleTimeNS = time.Duration(maxIdleTimeNS)
			}
		}
	}

	if err := viper.BindEnv("SERVER_PORT"); err == nil {
		viper.SetDefault("SERVER_PORT", conf.Server.Port)
		if port, ok := viper.Get("SERVER_PORT").(string); ok {
			if parsed, err := strconv.Atoi(port); err == nil {
				conf.Server.Port = parsed
			}
		}
	}
	if err := viper.BindEnv("DATABASE_NAME"); err == nil {
		viper.SetDefault("DATABASE_NAME", conf.Database.Name)
		if name, ok := viper.Get("DATABASE_NAME").(string); ok {
			conf.Database.Name = name
		}
	}
	if err := viper.BindEnv("DATABASE_USERNAME"); err == nil {
		viper.SetDefault("DATABASE_USERNAME", conf.Database.Username)
		if username, ok := viper.Get("DATABASE_USERNAME").(string); ok {
			conf.Database.Username = username
		}
	}
	if err := viper.BindEnv("DATABASE_PASSWORD"); err == nil {
		viper.SetDefault("DATABASE_PASSWORD", conf.Database.Password)
		if password, ok := viper.Get("DATABASE_PASSWORD").(string); ok {
			conf.Database.Password = password
		}
	}
	if err := viper.BindEnv("DATABASE_HOST"); err == nil {
		viper.SetDefault("DATABASE_HOST", conf.Database.Host)
		if host, ok := viper.Get("DATABASE_HOST").(string); ok {
			conf.Database.Host = host
		}
	}
	if err := viper.BindEnv("DATABASE_PORT"); err == nil {
		viper.SetDefault("DATABASE_PORT", conf.Database.Port)
		if port, ok := viper.Get("DATABASE_PORT").(string); ok {
			if parsed, err := strconv.Atoi(port); err == nil {
				conf.Database.Port = parsed
			}
		}
	}
	if err := viper.BindEnv("DATABASE_MAX_CONNS"); err == nil {
		viper.SetDefault("DATABASE_MAX_CONNS", conf.Database.MaxConns)
		if maxConns, ok := viper.Get("DATABASE_MAX_CONNS").(string); ok {
			if parsed, err := strconv.Atoi(maxConns); err == nil {
				conf.Database.MaxConns = parsed
			}
		}
	}
	if err := viper.BindEnv("DATABASE_MIN_CONNS"); err == nil {
		viper.SetDefault("DATABASE_MIN_CONNS", conf.Database.MinConns)
		if minConns, ok := viper.Get("DATABASE_MIN_CONNS").(string); ok {
			if parsed, err := strconv.Atoi(minConns); err == nil {
				conf.Database.MinConns = parsed
			}
		}
	}
	if err := viper.BindEnv("DATABASE_MAX_IDLE_TIME"); err == nil {
		viper.SetDefault("DATABASE_MAX_IDLE_TIME", conf.Database.MaxIdleTimeNS)
		if maxIdleTimeNS, ok := viper.Get("DATABASE_MAX_IDLE_TIME").(string); ok {
			if parsed, err := strconv.ParseInt(maxIdleTimeNS, 10, 64); err == nil {
				conf.Database.MaxIdleTimeNS = time.Duration(parsed)
			}
		}
	}

	return &conf, nil
}
