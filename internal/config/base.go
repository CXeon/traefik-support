package config

import (
	"context"

	configCli "github.com/CXeon/tiles/config"
	"github.com/CXeon/tiles/config/apollo"
	"github.com/CXeon/tiles/config/viper"
)

const baseLocalFileName = "base"

type base struct {
	JWT struct {
		Secret     string `json:"secret" mapstructure:"secret"`
		AccessTTL  int    `json:"access_ttl" mapstructure:"access_ttl"`   // seconds
		RefreshTTL int    `json:"refresh_ttl" mapstructure:"refresh_ttl"` // seconds
	} `json:"jwt" mapstructure:"jwt"`
	Log struct {
		Filename     string `json:"filename" mapstructure:"filename"`
		Level        string `json:"level" mapstructure:"level"`
		MaxSize      int    `json:"max_size" mapstructure:"max_size"`
		MaxBackups   int    `json:"max_backups" mapstructure:"max_backups"`
		MaxAge       int    `json:"max_age" mapstructure:"max_age"`
		Compress     bool   `json:"compress" mapstructure:"compress"`
		EnableStdout bool   `json:"enable_stdout" mapstructure:"enable_stdout"`
	} `json:"log" mapstructure:"log"`
	Gateway struct {
		Enabled  bool `json:"enabled" mapstructure:"enabled"`
		TTL      int  `json:"ttl" mapstructure:"ttl"`
		Weight   int  `json:"weight" mapstructure:"weight"`
		Provider struct {
			Endpoints   []string `json:"endpoints" mapstructure:"endpoints"`
			Username    string   `json:"username" mapstructure:"username"`
			Password    string   `json:"password" mapstructure:"password"`
			DialTimeout int      `json:"dial_timeout" mapstructure:"dial_timeout"` // seconds
			Namespace   string   `json:"namespace" mapstructure:"namespace"`
		} `json:"provider" mapstructure:"provider"`
	} `json:"gateway" mapstructure:"gateway"`
	Registry struct {
		Enabled bool `json:"enabled" mapstructure:"enabled"`
		Weight  int  `json:"weight" mapstructure:"weight"`
		Watch   struct {
			Services []string `json:"services" mapstructure:"services"`
			ComProj  []struct {
				Company  string   `json:"company" mapstructure:"company"`
				Projects []string `json:"projects" mapstructure:"projects"`
			} `json:"com_proj" mapstructure:"com_proj"`
		} `json:"watch" mapstructure:"watch"`
		Provider struct {
			Endpoints            []string `json:"endpoints" mapstructure:"endpoints"`
			Username             string   `json:"username" mapstructure:"username"`
			Password             string   `json:"password" mapstructure:"password"`
			DialTimeout          int      `json:"dial_timeout" mapstructure:"dial_timeout"` // seconds
			LoadBalancerStrategy uint8    `json:"load_balancer_strategy" mapstructure:"load_balancer_strategy"`
		} `json:"provider" mapstructure:"provider"`
	} `json:"registry" mapstructure:"registry"`
}

type BaseLoader struct {
	cli configCli.Config
}

func NewBaseLoader(opt Option) *BaseLoader {
	var cli configCli.Config
	if opt.f.ConfigMode == Remote {
		cfg := apollo.Config{
			AppID:          opt.f.ServiceName,
			Cluster:        opt.f.Cluster,
			IP:             opt.remote.IP,
			NamespaceName:  opt.remote.NamespaceName,
			Secret:         opt.remote.Secret,
			IsBackupConfig: true,
		}
		cli = apollo.New(cfg)
	}
	if opt.f.ConfigMode != Remote {
		cfg := viper.Config{
			ConfigPaths: opt.local.paths,
			ConfigName:  baseLocalFileName,
			ConfigType:  opt.local.typ,
			EnvPrefix:   opt.local.envPrefix,
			AutoEnv:     opt.local.autoEnv,
		}
		cli = viper.New(cfg)
	}
	return &BaseLoader{cli: cli}
}

func (l *BaseLoader) Load() (*base, error) {
	err := l.cli.Load()
	if err != nil {
		return nil, err
	}
	base, err := l.Unmarshal()
	if err != nil {
		return nil, err
	}

	return base, nil
}

func (l *BaseLoader) Unmarshal() (*base, error) {
	var cfg struct {
		Base base `mapstructure:"base"`
	}
	if err := l.cli.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg.Base, nil
}

func (l *BaseLoader) Watch(handler configCli.ChangeHandler) error {
	return l.cli.Watch(handler)
}

func (l *BaseLoader) Close() error {
	return l.cli.Close(context.Background())
}
