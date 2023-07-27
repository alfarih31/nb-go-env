package env

import (
	"errors"
	"fmt"
	"github.com/alfarih31/nb-go-parser"
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
	cs ConfigServer
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
	cfg, exist := c.cs.Get(k)
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
	cfg, exist := c.cs.Get(k)
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
	cfg, exist := c.cs.Get(k)
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

func (c env) GetStringArr(k string, def ...[]string) ([]string, error) {
	cfg, exist := c.cs.Get(k)
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
	cfg, exist := c.cs.Get(k)
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
	return c.cs.Dump()
}

func LoadWithConfigServer(configServer ConfigServer) (Env, error) {
	if configServer == nil {
		return nil, errors.New("config server can't be nill")
	}

	return env{
		cs: configServer,
	}, nil
}

func LoadEnv(envPath string, fallbackToWide ...bool) (Env, error) {
	cs, err := NewDefaultConfigServer(envPath, fallbackToWide...)
	if err != nil {
		return nil, err
	}

	return LoadWithConfigServer(cs)
}
