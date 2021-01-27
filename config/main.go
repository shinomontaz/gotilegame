package config

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"

	_ "net/http/pprof"
)

type Config struct {
	Title      string
	TilesDir   string
	Fullscreen bool
	TestFlag   bool
}

type Env struct {
	Cfg      *Config
	loglevel log.Level
}

func NewEnv(path string) *Env {
	var cfg Config
	err := gonfig.GetConf(path+"/"+"conf.json", &cfg)
	if err != nil {
		err = gonfig.GetConf(path+"/"+"conf.tpl.json", &cfg)
		checkErr(err)
	}

	return &Env{
		Cfg:      &cfg,
		loglevel: log.WarnLevel,
	}
}

func (e *Env) InitLog() {
	if e.Cfg.TestFlag {
		e.loglevel = log.DebugLevel
	}

	log.SetLevel(e.loglevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func checkErr(err error) {
	if err != nil {
		_, filename, lineno, ok := runtime.Caller(1)
		message := ""
		if ok {
			message = fmt.Sprintf("%v:%v: %v\n", filename, lineno, err)
		}
		log.Panic(message, err)
	}
}
