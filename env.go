package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alfarih31/nb-go-parser"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
)

type Err struct {
	e error
}

// HasZeroValue Check a variable has Zero Value
func hasZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}

	t := reflect.TypeOf(v)
	if t == nil {
		return true
	}

	switch t.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return false
	}

	return v == reflect.Zero(t).Interface()
}

func (e *Err) Errorf(f string, s ...interface{}) *Err {
	return &Err{
		e: errors.New(fmt.Sprintf(fmt.Sprintf("this %s: %s", e.e.Error(), f), s...)),
	}
}

func (e *Err) Error() string {
	return e.e.Error()
}

func NewEnvErr(s string) *Err {
	return &Err{
		e: errors.New(s),
	}
}

var ErrorVarNotExist = NewEnvErr("variable not exist")

type env struct {
	envs      map[string]string
	useDotEnv bool
}

type Env interface {
	GetInt(k string, def ...int) (int, error)
	MustGetInt(k string, def ...int) int
	GetString(k string, def ...string) (string, error)
	MustGetString(k string, def ...string) string
	GetBool(k string, def ...bool) (bool, error)
	MustGetBool(k string, def ...bool) bool
	GetStringArr(k string, def ...[]string) ([]string, error)
	MustGetStringArr(k string, def ...[]string) []string
	GetIntArr(k string, def ...[]int) ([]int, error)
	MustGetIntArr(k string, def ...[]int) []int
	Dump() (string, error)
}

func (c env) GetInt(k string, def ...int) (int, error) {
	cfg, exist := c.get(k)
	if !exist {
		if len(def) == 0 {
			return 0, ErrorVarNotExist.Errorf("%s", k)
		}
		return def[0], nil
	}

	i, e := parser.String(cfg).ToInt()

	if e != nil {
		if len(def) == 0 {
			return 0, e
		}

		return def[0], nil
	}

	return i, e
}

func (c env) MustGetInt(k string, def ...int) int {
	v, err := c.GetInt(k, def...)
	if err != nil {
		panic(err)
	}

	return v
}

func (c env) GetString(k string, def ...string) (string, error) {
	cfg, exist := c.get(k)
	if !exist {
		if len(def) == 0 {
			return "", ErrorVarNotExist.Errorf("%s", k)
		}
		return def[0], nil
	}

	return cfg, nil
}

func (c env) MustGetString(k string, def ...string) string {
	v, err := c.GetString(k, def...)
	if err != nil {
		panic(err)
	}

	return v
}

func (c env) GetBool(k string, def ...bool) (bool, error) {
	cfg, exist := c.get(k)
	if !exist {
		if len(def) == 0 {
			return false, ErrorVarNotExist.Errorf("%s", k)
		}
		return def[0], nil
	}

	b, e := parser.String(cfg).ToBool()
	if e != nil {
		if len(def) == 0 {
			return false, e
		}

		return def[0], nil
	}

	return b, e
}

func (c env) MustGetBool(k string, def ...bool) bool {
	v, err := c.GetBool(k, def...)
	if err != nil {
		panic(err)
	}

	return v
}

func (c env) get(k string) (string, bool) {
	if c.useDotEnv {
		cfg, exist := c.envs[k]

		return cfg, !hasZeroValue(cfg) && exist
	}

	cfg := os.Getenv(k)
	return cfg, !hasZeroValue(cfg)
}

func (c env) GetStringArr(k string, def ...[]string) ([]string, error) {
	cfg, exist := c.get(k)
	if !exist {
		if len(def) == 0 {
			return nil, ErrorVarNotExist.Errorf("%s", k)
		}

		return def[0], nil
	}

	return parser.String(cfg).ToStringArr()
}

func (c env) MustGetStringArr(k string, def ...[]string) []string {
	v, err := c.GetStringArr(k, def...)
	if err != nil {
		panic(err)
	}

	return v
}

func (c env) GetIntArr(k string, def ...[]int) ([]int, error) {
	cfg, exist := c.get(k)
	if !exist {
		if len(def) == 0 {
			return nil, ErrorVarNotExist.Errorf("%s", k)
		}
		return def[0], nil
	}

	is, e := parser.String(cfg).ToIntArr()

	if e != nil {
		if len(def) == 0 {
			return nil, e
		}

		return def[0], nil
	}

	return is, e
}

func (c env) MustGetIntArr(k string, def ...[]int) []int {
	v, err := c.GetIntArr(k, def...)
	if err != nil {
		panic(err)
	}

	return v
}

func (c env) Dump() (string, error) {
	if !c.useDotEnv {
		return "", errors.New("cannot dump env, you are using system-wide env")
	}

	j, e := json.Marshal(c.envs)

	return string(j), e
}

func LoadEnv(envPath string, fallbackToWide ...bool) (Env, error) {
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
		return env{
			envs:      envs,
			useDotEnv: false,
		}, nil
	}

	for key, val := range envs {
		err = os.Setenv(key, val)
	}

	return env{
		envs:      envs,
		useDotEnv: err == nil,
	}, nil
}
