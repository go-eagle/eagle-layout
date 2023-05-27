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
	Redis struct {
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

	}
	Tasks []Task
}

type Task struct {
	Name     string
	Schedule string
}

func NewClient(cfg *Config) *asynq.Client {
	once.Do(func() {
		if cfg == nil {
			panic("tasks config is nil")
		}
		client = asynq.NewClient(asynq.RedisClientOpt{
			Addr:         cfg.Redis.Addr,
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
			DialTimeout:  cfg.Redis.DialTimeout,
			ReadTimeout:  cfg.Redis.ReadTimeout,
			WriteTimeout: cfg.Redis.WriteTimeout,
			PoolSize:     cfg.Redis.PoolSize,
		})
	})
	return client
}
