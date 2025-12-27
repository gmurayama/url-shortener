package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/ardanlabs/conf/v3"
)

type Config struct {
	Application struct {
		Name string `conf:"default:URL Shortener"`
	}
	Server struct {
		Address         string        `conf:"default:0.0.0.0:8080"`
		ReadTimeout     time.Duration `conf:"default:2s"`
		WriteTimeout    time.Duration `conf:"default:2s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}
	Tracing struct {
		Host               string        `conf:"default:localhost"`
		Port               int           `conf:"default:4317"`
		Enabled            bool          `conf:"default:true"`
		BatchScheduleDelay time.Duration `conf:"default:5s"`
		SamplingRatio      float64       `conf:"default:1.0"`
		MaxExportBatchSize int           `conf:"default:256"`
		KeepAliveTime      time.Duration `conf:"default:20s"`
		KeepAliveTimeout   time.Duration `conf:"default:5s"`
	}
	InternalServer struct {
		Address      string        `conf:"default:0.0.0.0:9090"`
		ReadTimeout  time.Duration `conf:"default:5s"`
		WriteTimeout time.Duration `conf:"default:5s"`
		EnablePprof  bool          `conf:"default:false"`
	}
}

func New() (Config, error) {
	prefix := "APP"
	var cfg Config
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return cfg, err
		}
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, err
}
