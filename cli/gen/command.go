package gen

import (
	"gitee.com/conero/uymas/v2/cli"
	"gitee.com/conero/uymas/v2/cli/evolve"
	"gitee.com/conero/uymas/v2/str"
	"reflect"
)

type Runnable interface {
	cli.Fn | func()
}

type StructCmd struct {
	indexRn     any
	lostRn      any
	initRn      any
	commandList map[string]any
	contextVal  reflect.Value
}

func (s *StructCmd) Index() any {
	return s.indexRn
}

func (s *StructCmd) Lost() any {
	return s.lostRn
}

func (s *StructCmd) Init() any {
	return s.initRn
}

func ParseStruct(vStruct any) *StructCmd {
	if vStruct == nil {
		return nil
	}
	rv := reflect.ValueOf(vStruct)
	realVal := rv
	if rv.Kind() == reflect.Ptr {
		realVal = rv.Elem()
	}

	if realVal.Kind() != reflect.Struct {
		return nil
	}

	num := rv.NumMethod()
	rType := rv.Type()
	sc := &StructCmd{
		commandList: map[string]any{},
	}
	sc.contextVal = rv
	for i := 0; i < num; i++ {
		method := rType.Method(i)
		methodValue := rv.Method(i)
		if !methodValue.CanInterface() {
			continue
		}
		name := method.Name
		vAny := methodValue.Interface()
		switch name {
		case evolve.CmdMtdIndex:
			sc.indexRn = vAny
		case evolve.CmdMtdLost:
			sc.lostRn = vAny
		case evolve.CmdMtdInit:
			sc.initRn = vAny
		default:
			sc.commandList[str.Str(name).Lcfirst()] = vAny
		}
	}

	return sc
}

func (s *StructCmd) SetArgs(args cli.ArgsParser) {
	if s.contextVal.IsNil() || !s.contextVal.IsValid() || s.contextVal.IsZero() {
		return
	}
	value := s.contextVal
	if s.contextVal.Kind() == reflect.Ptr {
		value = s.contextVal.Elem()
	}
	field := value.FieldByName(evolve.CmdFidArgs)
	if !field.CanSet() {
		return
	}
	field.Set(reflect.ValueOf(args))
}

func AsCommand(vStruct any, cfgs ...cli.Config) cli.Application[any] {
	pCmd := ParseStruct(vStruct)
	if pCmd == nil {
		panic("vStruct is not struct, and parse fail")
	}
	evl := evolve.NewEvolve(cfgs...)
	evl.Lost(pCmd.Lost())
	evl.Index(pCmd.Index())
	evl.RouterBefore(func(args cli.ArgsParser) {
		pCmd.SetArgs(args)
	})

	for vCmd, runnable := range pCmd.commandList {
		evl.Command(runnable, vCmd)
	}
	return evl
}
