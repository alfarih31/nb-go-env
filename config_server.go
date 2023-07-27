package env

import (
	"encoding/json"
	"errors"
	"github.com/alfarih31/nb-go-env/internal"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ConfigServer interface {
	Get(key string) (string, bool)
	Dump() (string, error)
}

type defaultConfigServer struct {
	envs      map[string]string
	useDotEnv bool
}

func (dc defaultConfigServer) Dump() (string, error) {
	if !dc.useDotEnv {
		return "", errors.New("cannot dump env, you are using system-wide env")
	}

	j, e := json.Marshal(dc.envs)

	return string(j), e
}

var _ ConfigServer = new(defaultConfigServer)

func (dc defaultConfigServer) Get(k string) (string, bool) {
	if dc.useDotEnv {
		cfg, exist := dc.envs[k]

		return cfg, !internal.HasZeroValue(cfg) && exist
	}

	cfg := os.Getenv(k)
	return cfg, !internal.HasZeroValue(cfg)
}

func NewDefaultConfigServer(envPath string, fallbackToWide ...bool) (ConfigServer, error) {
	fBackToWide := false
	if len(fallbackToWide) > 0 {
		fBackToWide = fallbackToWide[0]
	}

	envs, err := godotenv.Read(envPath)
	if err != nil && !fBackToWide {
		return nil, err
	}

	if err != nil {
		if !fBackToWide {
			return nil, err
		}

		log.Printf("%s \n %s", "Warning! You are use System Wide Environment due to this error:", err.Error())
		return defaultConfigServer{
			envs:      envs,
			useDotEnv: false,
		}, nil
	}

	for key, val := range envs {
		err = os.Setenv(key, val)
	}

	return defaultConfigServer{
		envs:      envs,
		useDotEnv: err == nil,
	}, nil
}
