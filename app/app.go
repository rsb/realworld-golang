// Package app is responsible for application specific concerns including
// start, stop, or any other admin/cli tasks, configuration, dependency
// injection and all entry points api, cli etc...
package app

import (
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

const (
	ServiceName = "lola-service"
)

type Dependencies struct {
	ServiceName     string
	Build           string
	Host            string
	DebugHost       string
	ReadTimout      time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	Kubernetes      KubeInfo
	Shutdown        chan os.Signal
	Logger          *zap.SugaredLogger
}

type KubeInfo struct {
	Pod       string
	PodIP     string
	Node      string
	Namespace string
}

// RootDir is designed to return the absolute path of the directory of
// this project
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
