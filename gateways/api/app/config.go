package app

import "time"

type Config struct {
	Application struct {
		Name            string        `env:"NAME" envDefault:"URL Shortener"`
		Address         string        `env:"ADDR" envDefault:"localhost"`
		Port            int           `env:"PORT" envDefault:"8080"`
		ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"2s"`
		WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"2s"`
		Timeout         time.Duration `env:"TIMEOUT" envDefault:"2s"`
		ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	} `envPrefix:"APP_"`
	Internal struct {
		Address      string        `env:"ADDR" envDefault:"localhost"`
		Port         int           `env:"PORT" envDefault:"8081"`
		ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"2s"`
		WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"2s"`
	} `envPrefix:"INTERNAL_"`
	Tracing struct {
		Host               string        `env:"HOST" envDefault:"localhost"`
		Port               int           `env:"PORT" envDefault:"4317"`
		Enabled            bool          `env:"ENABLED" envDefault:"true"`
		BatchScheduleDelay time.Duration `env:"BATCH_SCHEDULE_DELAY" envDefault:"5s"`
		SamplingRatio      float64       `env:"SAMPLING_RATIO" envDefault:"1.0"`
		MaxExportBatchSize int           `env:"MAX_EXPORT_BATCH_SIZE" envDefault:"256"`
		KeepAliveTime      time.Duration `env:"KEEP_ALIVE_TIME" envDefault:"20s"`
		KeepAliveTimeout   time.Duration `env:"KEEP_ALIVE_TIMEOUT" envDefault:"5s"`
	} `envPrefix:"TRACING_"`
}
