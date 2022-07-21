package tasks

import (
	"sync"
	"time"

	"github.com/hibiken/asynq"
)

const (
	// queue name
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

var (
	client *asynq.Client
	once   sync.Once
)

type Config struct {
	Addr         string
	Password     string
	DB           int
	MinIdleConn  int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	PoolTimeout  time.Duration
	Concurrency  int //并发数
	Jobs         []Job
}

type Job struct {
	Name     string
	Schedule string
}

func NewClient(cfg *Config) *asynq.Client {
	once.Do(func() {
		if cfg == nil {
			panic("tasks client is nil")
		}
		client = asynq.NewClient(asynq.RedisClientOpt{
			Addr:         cfg.Addr,
			Password:     cfg.Password,
			DB:           cfg.DB,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			PoolSize:     cfg.PoolSize,
		})
	})
	return client
}
