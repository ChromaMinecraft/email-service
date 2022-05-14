package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host string `env:"APP_HOST" envDefault:"localhost"`
	Port int    `env:"APP_PORT" envDefault:"8080"`
	TZ   string `env:"APP_TZ" envDefault:"Asia/Jakarta"`

	LogLevel string `env:"LOG_LEVEL"`

	GracefulShutdownTimeout int `env:"GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"5"`

	MailHost     string `env:"MAIL_HOST" envDefault:"smtp.gmail.com"`
	MailPort     int    `env:"MAIL_PORT" envDefault:"587"`
	MailUser     string `env:"MAIL_USER"`
	MailPassword string `env:"MAIL_PASS"`

	WorkerCount int `env:"WORKER_COUNT" envDefault:"3"`

	AMQP
}

type AMQP struct {
	Host     string `env:"AMQP_HOST" envDefault:"localhost"`
	Port     int    `env:"AMQP_PORT" envDefault:"5672"`
	User     string `env:"AMQP_USER" envDefault:"guest"`
	Password string `env:"AMQP_PASSWORD" envDefault:"guest"`

	AMQPQoS
	AMQPQueue
}

type AMQPQoS struct {
	PrefetchCount int  `env:"AMQP_PREFETCH_COUNT" envDefault:"1"`
	PrefetchSize  int  `env:"AMQP_PREFETCH_SIZE" envDefault:"0"`
	IsGlobal      bool `env:"AMQP_IS_GLOBAL" envDefault:"false"`
}

type AMQPQueue struct {
	Name         string `env:"AMQP_QUEUE_NAME"`
	IsDurable    bool   `env:"AMQP_QUEUE_IS_DURABLE" envDefault:"true"`
	IsAutoDelete bool   `env:"AMQP_QUEUE_IS_AUTO_DELETE" envDefault:"false"`
	IsExclusive  bool   `env:"AMQP_QUEUE_IS_EXCLUSIVE" envDefault:"false"`
	IsNoWait     bool   `env:"AMQP_QUEUE_IS_NO_WAIT" envDefault:"false"`
}

func Get() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
