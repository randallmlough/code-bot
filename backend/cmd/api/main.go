package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"sync"

	"github.com/randallmlough/code-bot/internal/env"
	"github.com/randallmlough/code-bot/internal/leveledlog"
	"github.com/randallmlough/code-bot/internal/openai"
	"github.com/randallmlough/code-bot/internal/version"
)

func main() {
	logger := leveledlog.NewLogger(os.Stdout, leveledlog.LevelAll, true)
	err := run(logger)
	if err != nil {
		logger.Fatal(err, debug.Stack())
	}
}

type config struct {
	baseURL     string
	httpPort    string
	openAIToken string
}

type application struct {
	config config
	logger *leveledlog.Logger
	wg     sync.WaitGroup
	gpt    *openai.OpenAI
}

func run(logger *leveledlog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:8888")
	cfg.httpPort = env.GetString("PORT", "8888")
	cfg.openAIToken = env.GetString("OPENAI_TOKEN", os.Getenv("OPENAI_TOKEN"))

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	gClient := openai.New(cfg.openAIToken)
	app := &application{
		config: cfg,
		logger: logger,
		gpt:    gClient,
	}

	return app.serveHTTP()
}
