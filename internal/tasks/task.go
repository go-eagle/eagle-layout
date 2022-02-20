package tasks

import (
	"log"
	"sync"
	"time"

	"github.com/go-eagle/eagle/pkg/config"

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
}

func GetClient() *asynq.Client {
	once.Do(func() {
		c := config.New("config", config.WithEnv("local"))
		var cfg Config
		if err := c.Load("cron", &cfg); err != nil {
			panic(err)
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

func Example() {
	// ------------------------------------------------------
	// Enqueue task to be processed immediately.
	// Use (*Client).Enqueue method.
	// ------------------------------------------------------
	task, err := NewEmailWelcomeTask(1)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := GetClient().Enqueue(task, asynq.Queue(QueueDefault))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Schedule task to be processed in the future.
	// Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------
	task, err = NewEmailWelcomeTask(2)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = GetClient().Enqueue(task, asynq.ProcessIn(10*time.Second))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Set other options to tune task processing behavior.
	// Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------
	task, err = NewEmailWelcomeTask(3)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = GetClient().Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
