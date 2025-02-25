package gen

import (
	"errors"
	"gitee.com/conero/uymas/v2/cli"
	"gitee.com/conero/uymas/v2/data/convert"
	"gitee.com/conero/uymas/v2/data/input"
	"gitee.com/conero/uymas/v2/rock"
	"gitee.com/conero/uymas/v2/str"
	"reflect"
	"regexp"
	"strings"
)

// ArgsTagName The struct tag is named cmd. If this parameter is not set,
// the struct tag is divided into lowercase cases, such as file_name, make_up.
//
// json tag is also supported when cmd is not set.
//
// syntax rules of name: "cmd>>json>>FileName".
const ArgsTagName = "cmd"

// ArgsCmdRequired the option of command is required
const ArgsCmdRequired = "required"

// ArgsCmdHelp the help information of the option
const ArgsCmdHelp = "help"

// ArgsCmdDefault to set option a default value
const ArgsCmdDefault = "default"

// ArgsTagOmit ignore command line properties, that is, do not parse
const ArgsTagOmit = "-"

// ArgsTagData as a mapping identifier for the data in the command
const ArgsTagData = "isdata"

// ArgsTagNext as a command output data specific index, default index=1
const ArgsTagNext = "next"

// ArgsTagMark the option document is identified as the value of the option when it is generated
const ArgsTagMark = "mark" // tag option input values name
// ArgsTagOwner a tag in a structure command that is an option to the command it belongs to
const ArgsTagOwner = "owner"

// ArgsGlobalOwner as an identifier for global options
const ArgsGlobalOwner = "globalOwner"

// ArgsOptionNoValid the command does not verify option data
const ArgsOptionNoValid = "notValid"

// ArgsCmdStructGen marking it requires parsing of substructure documents and values
const ArgsCmdStructGen = "structGen"

func argsValueCheck(ref reflect.Value) (reflect.Value, error) {
	structMark := ref.Kind() == reflect.Struct
	isPtr := false
	if ref.Kind() == reflect.Ptr {
		structMark = ref.Elem().Kind() == reflect.Struct
		isPtr = true
	}

	if !structMark {
		return reflect.Value{}, errors.New("data: the param of ArgsDress only support struct")
	}

	rValue := ref
	if isPtr {
		rValue = ref.Elem()
	}

	return rValue, nil
}

func setValueByOption(vField reflect.Value, option *cli.Option, args cli.ArgsParser, keys []string) {
	if option == nil {
		return
	}
	// data
	if option.IsData {
		valueStr := rock.ListGetOr(args.CommandList(), option.Next-1, args.SubCommand())
		convert.SetByStr(vField, valueStr)
		return
	}

	// set struct gen
	if option.StructGen {
		if vField.Kind() != reflect.Struct {
			return
		}
		if len(option.StructItems) < 1 {
			option.StructItems = StructDress(vField)
		}
		name := option.GetName()
		cfg := cli.ConfigWith()
		for _, item := range option.StructItems {
			structName := item.GetName()
			var key = name + cfg.StructGenSep + structName
			var valueStr = args.Get(key)
			if valueStr == "" {
				valueStr = item.DefValue
			}

			structName = str.Str(structName).Ucfirst()
			childFld := vField.FieldByName(structName)
			if !childFld.IsValid() {
				continue
			}
			setValueByStr(childFld, []string{key}, args, valueStr)
		}
		return
	}

	value := args.Get(keys...)
	if value == "" {
		value = option.DefValue
	}

	setValueByStr(vField, keys, args, value)
}

func setValueByStr(vField reflect.Value, keys []string, args cli.ArgsParser, defStrs ...string) {
	vfKind := vField.Kind()
	if vfKind == reflect.Bool && args.Switch(keys...) {
		vField.SetBool(true)
		return
	}

	vSlice := args.List(keys...)
	if vfKind == reflect.Slice && vSlice != nil {
		convert.SetByStrSlice(vField, vSlice)
		return
	}

	value := args.Get(keys...)
	value = rock.Param(value, defStrs...)
	convert.SetByStr(vField, value)
}

func setToStruct(tgt reflect.Value, args cli.ArgsParser) {
	if tgt.Kind() != reflect.Struct {
		return
	}
	rtp := tgt.Type()
	for i := 0; i < tgt.NumField(); i++ {
		fieldType := rtp.Field(i)
		name := fieldType.Tag.Get(ArgsTagName)
		tagValue := name
		if name == "" {
			name = str.JsonTagName(fieldType.Tag.Get("json"))
		}
		if name == "" {
			name = str.Str(fieldType.Name).LowerStyle()
		}

		if name == ArgsTagOmit {
			continue
		}

		// field inherit by parent struct.
		if fieldType.Anonymous {
			setToStruct(tgt.Field(i), args)
			continue
		}

		keys := getNameByTag(name)
		if len(keys) == 0 {
			continue
		}
		option := OptionTagParse(tagValue)
		setValueByOption(tgt.Field(i), option, args, keys)
	}
}

// ArgsDress Dress the command argument up on the specified data entity (struct)
func ArgsDress(args cli.ArgsParser, data any) error {
	ref := reflect.ValueOf(data)
	realValue, err := argsValueCheck(ref)
	if err != nil {
		return err
	}

	setToStruct(realValue, args)
	return nil
}

// ArgsDecompose Decompose the structure into an option list
func ArgsDecompose(data any, excludes ...string) ([]cli.Option, error) {
	ref := reflect.ValueOf(data)
	realValue, err := argsValueCheck(ref)
	if err != nil {
		return nil, err
	}
	return StructDress(realValue, excludes...), nil
}

// StructDress dress up the struct property value (which supports composition/inheritance) on `cli.Option`
func StructDress(vStruct reflect.Value, excludes ...string) (inheritOpts []cli.Option) {
	if vStruct.Kind() != reflect.Struct {
		return
	}
	vType := vStruct.Type()
	num := vStruct.NumField()
	for i := 0; i < num; i++ {
		field := vStruct.Field(i)
		sField := vType.Field(i)
		cmdTag := sField.Tag.Get(ArgsTagName)
		if cmdTag == ArgsTagOmit {
			continue
		}
		// field inherit by parent struct.
		if sField.Anonymous {
			inheritOpts = append(inheritOpts, StructDress(field, excludes...)...)
			continue
		}
		option := OptionTagParse(cmdTag)
		if option == nil {
			continue
		}
		option.FieldName = sField.Name
		//@todo optimize code follow business
		var name string
		if option == nil {
			if name == "" {
				name = str.JsonTagName(sField.Tag.Get("json"))
			}
			if name == "" {
				name = str.Str(sField.Name).LowerStyle()
			}
			if rock.InList(excludes, name) {
				continue
			}
			option = &cli.Option{
				Alias: []string{name},
			}
		} else if rock.InList(excludes, option.Name) {
			continue
		}
		// parse child items
		if option.StructGen {
			option.StructItems = StructDress(field)
		}
		inheritOpts = append(inheritOpts, *option)
	}
	return
}

func ArgsDecomposeMust(data any, excludes ...string) []cli.Option {
	opts, _ := ArgsDecompose(data, excludes...)
	return opts
}

// OptionTagParse Resolves the value of the tag into an option object
//
// syntax rules of tag: `"name,n required default:111 help:help msg"`.
//
// When using command data instead of options, you can specify `next` or default `subCommand`.
//
// `"input isdata next:2"`
func OptionTagParse(vTag string) *cli.Option {
	if vTag == "" {
		return nil
	}
	spaceList := regexp.MustCompile(`\s{2,}`)
	vTag = spaceList.ReplaceAllString(vTag, " ")
	vTag = strings.TrimSpace(vTag)
	if vTag == "" {
		return nil
	}

	option := &cli.Option{}
	toSetOptionByKeyFn := func(key string) {
		switch key {
		case ArgsCmdStructGen:
			option.StructGen = true
		case ArgsTagData:
			option.IsData = true
		}
	}
	for i, s := range strings.Split(vTag, " ") {
		if !strings.Contains(s, ":") {
			name := str.Str(s).ClearSpace()
			option.List = append(option.List, name)
			if i == 0 {
				option.Alias = strings.Split(name, ",")
			}
			toSetOptionByKeyFn(s)
			continue
		}
		toSetOptionByKeyFn(s)
		if s == ArgsCmdRequired {
			continue
		}
		if s == ArgsTagNext {
			option.Next = 1
			continue
		}
		if s == ArgsTagMark {
			option.Mark = option.GetName()
			continue
		}
		idx := strings.Index(s, ":")
		if idx > 0 {
			key := s[:idx]
			value := s[idx+1:]
			switch key {
			case ArgsCmdHelp:
				option.Help = str.Str(value).Unescape()
			case ArgsCmdDefault:
				option.DefValue = str.Str(value).Unescape()
			case ArgsTagNext:
				option.Next = input.Stringer(value).Int()
			case ArgsTagMark:
				option.Mark = value
			case ArgsTagOwner:
				option.Owner = value
			}
		}
	}

	return option
}

// Get the name by parsing the tag of struct, format like `cmd:"name,n"`
func getNameByTag(tag string) []string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}
	var name string
	for _, vs := range strings.Split(tag, " ") {
		vs = strings.TrimSpace(vs)
		if vs == "" {
			continue
		}
		if strings.Contains(vs, ":") {
			continue
		}

		name = vs
		break

	}

	if len(name) > 0 {
		return strings.Split(str.Str(name).ClearSpace(), ",")
	}
	return []string{tag}
}
