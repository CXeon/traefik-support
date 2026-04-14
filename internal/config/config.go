package config

import (
	"errors"
	"sync"

	configCli "github.com/CXeon/tiles/config"
)

var Config *config
var once sync.Once

type config struct {
	// 记录启动参数作为配置，生命周期内无法修改
	Flags
	lextra *localExtra
	rextra *remoteExtra

	l          sync.RWMutex
	baseLoader *BaseLoader

	Base *base `json:"base"`
}

// Flags 包装程序启动参数，在程序的生命周期内不允许修改
type Flags struct {
	ConfigMode  Mode   // 配置模式 local 加载本地配置文件，remote 从远程配置中心加载
	Env         string // 环境名称
	Cluster     string // 集群名称
	Company     string
	Project     string
	ServiceName string // 服务名称
	Color       string // 染色
	HttpPort    uint   // 服务监听端口
}

type Mode string

const (
	Local  Mode = "local"
	Remote Mode = "remote"
)

// 本地配置补充
type localExtra struct {
	paths     []string
	typ       string
	envPrefix string
	autoEnv   bool
}

// 远程配置补充
type remoteExtra struct {
	// AppID         string `json:"app_id"`
	Cluster       string `json:"cluster"`
	Dev           string `json:"dev"`
	IP            string `json:"ip"`
	NamespaceName string `json:"namespace_name"`
	Secret        string `json:"secret"`
	// IsBackup      bool   `json:"is_backup"`
}

type Option struct {
	f      Flags
	local  *localExtra
	remote *remoteExtra
}

func Load(f Flags) error {
	// 判断配置模式
	if f.ConfigMode != Local && f.ConfigMode != Remote {
		return errors.New("invalid config mode")
	}

	once.Do(func() {
		Config = &config{
			Flags: f,
			l:     sync.RWMutex{},
		}

		if f.ConfigMode == Local {
			// 本地模式启动。补充配置参数
			Config.lextra = &localExtra{
				paths:     []string{"./configs/"},
				typ:       "yaml",
				envPrefix: "",
				autoEnv:   true,
			}
		}

		if f.ConfigMode == Remote {
			tmp, err := getRemoteExtra(f.Env, f.Cluster)
			if err != nil {
				panic(err)
			}
			Config.rextra = &tmp
		}

		opt := Option{
			f:      f,
			local:  Config.lextra,
			remote: Config.rextra,
		}

		// 加载base配置
		baseLoader := NewBaseLoader(opt)
		base, err := baseLoader.Load()
		if err != nil {
			panic(err)
		}
		Config.baseLoader = baseLoader
		Config.Base = base
	})

	return nil
}

func Watch(handler configCli.ChangeHandler) error {

	if Config == nil || Config.baseLoader == nil {
		return errors.New("config not loaded, call Load() first")
	}

	var err error
	err = Config.baseLoader.Watch(func(event *configCli.ChangeEvent) {
		// 1. 自动同步全局 Config.Base
		Config.l.Lock()
		newBase, err := Config.baseLoader.Unmarshal()
		if err == nil {
			Config.Base = newBase
		}
		Config.l.Unlock()

		// 2. 将变更事件传给调用方，由调用方决定副作用
		if handler != nil {
			handler(event)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	if Config == nil || Config.baseLoader == nil {
		return nil
	}
	var errs []error
	if err := Config.baseLoader.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// TODO 根据环境和集群获取远程配置中心IP等信息
func getRemoteExtra(env, cluster string) (remoteExtra, error) {
	panic("fix me")
}
