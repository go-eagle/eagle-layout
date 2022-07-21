package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-eagle/eagle-layout/internal/model"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
	v "github.com/go-eagle/eagle/pkg/version"
	"github.com/spf13/pflag"

	"github.com/hibiken/asynq"
)

var (
	cfgDir  = pflag.StringP("config dir", "c", "config", "config path.")
	env     = pflag.StringP("env name", "e", "", "env var name.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func init() {
	pflag.Parse()
	if *version {
		ver := v.Get()
		marshaled, err := json.MarshalIndent(&ver, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return
	}

	// init config
	c := config.New(*cfgDir, config.WithEnv(*env))
	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	// -------------- init resource -------------
	logger.Init()
	// init db
	model.Init()
	// init redis
	redis.Init()
}

func main() {
	// load config
	c := config.New(*cfgDir, config.WithEnv(*env))
	var cfg tasks.Config
	if err := c.Load("cronjob", &cfg); err != nil {
		panic(err)
	}

	// -------------- Run worker server ------------
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Addr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: cfg.Concurrency,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				tasks.QueueCritical: 6,
				tasks.QueueDefault:  3,
				tasks.QueueLow:      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	// register handlers...
	mux.HandleFunc(tasks.TypeEmailWelcome, tasks.HandleEmailWelcomeTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
