package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/CXeon/traefik-support/cmd"
	"github.com/CXeon/traefik-support/internal/config"
)

func main() {
	var (
		configMode  string
		env         string
		cluster     string
		company     string
		project     string
		serviceName string
		color       string
		httpPort    uint
	)

	flag.StringVarP(&configMode, "config-mode", "m", "local", "配置模式: local | remote")
	flag.StringVarP(&env, "env", "e", "dev", "环境名称")
	flag.StringVar(&cluster, "cluster", "local", "集群名称")
	flag.StringVar(&company, "company", "company", "公司标识")
	flag.StringVar(&project, "project", "project", "项目标识")
	flag.StringVarP(&serviceName, "service", "s", "traefik-support", "服务名称")
	flag.StringVar(&color, "color", "", "染色标记")
	flag.UintVarP(&httpPort, "port", "p", 18080, "HTTP 监听端口")
	flag.Parse()

	if serviceName == "" {
		fmt.Fprintln(os.Stderr, "error: --service is required")
		flag.Usage()
		os.Exit(1)
	}

	flags := config.Flags{
		ConfigMode:  config.Mode(configMode),
		Env:         env,
		Cluster:     cluster,
		Company:     company,
		Project:     project,
		ServiceName: serviceName,
		Color:       color,
		HttpPort:    httpPort,
	}

	app := cmd.New(flags)

	if err := app.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "init error: %v\n", err)
		os.Exit(1)
	}

	if err := app.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
		os.Exit(1)
	}
}
