package tasks

import (
	"sync"
	"time"

	"github.com/go-eagle/eagle/pkg/config"

	"github.com/hibiken/asynq"
)

const (
	// QueueCritical queue priority high
	QueueCritical = "critical"
	// QueueDefault queue priority middle
	QueueDefault  = "default"
	// QueueLow queue priority low
	QueueLow      = "low"
)

var (
    // client define a async client
	client *asynq.Client
	// once define a lock for get async client
	once   sync.Once
)

// Config define a struct
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
}

// Task define a task struct
type Task struct {
	Name     string
	Schedule string
}

// GetClient get a async client
func GetClient() *asynq.Client {
	once.Do(func() {
		c := config.New("config")
		var cfg Config
		if err := c.Load("consumer", &cfg); err != nil {
			panic(err)
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
