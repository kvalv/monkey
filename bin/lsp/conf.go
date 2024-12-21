package main

import (
	"flag"
	"io"
	"log/slog"
	"os"
)

type Config struct {
	r io.Reader
	w io.Writer
}

var pipe = flag.String("pipe", "", "pipe to use for communication")

func NewConfig(log *slog.Logger) *Config {
	flag.Parse()
	if *pipe != "" {
		log.Info("using pipe", "pipe", *pipe)
		// path to file
		f, err := os.OpenFile(*pipe, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		return &Config{
			r: f,
			w: f,
		}
	}
	return &Config{
		r: os.Stdin,
		w: os.Stdout,
	}
}
