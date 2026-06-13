package main

import (
	"log"
	"my-go-project/config"
	"my-go-project/pkg/logger"
	"my-go-project/pkg/utils"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func main() {
	log.Println("Starting api server")

	configName := utils.GetConfigName(os.Getenv("config"))

	configFile, err := config.LoadConfig(configName)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}
	_ = configFile
	log.Println("Config file loaded successfully")

	cfg, err := config.ParseConfig(configFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.DPanic()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

}
