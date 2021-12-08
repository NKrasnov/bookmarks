/*
config package provides basic app configuration capabilities
through command line parameters and environment variables

App configuration struct
 tags description
 cmd     - command line parameter name
 env     - environment variable name
 default - default value for the parameter.
           command line parameters take precedence over environment variables
 usage   - short description

 Example:

	cfg := struct{
		Host 	string `param:"cmd=host,env=APP_HOST,default=127.0.0.1,usage=hostname or IP"`
		Port 	int	   `param:"cmd=port,env=APP_PORT,default=8080,usage=posrt the server will listen to"`
		Timeout int    `param:"cmd=timeout,env=APP_TIMEOUT,default=10,usage=server read timeout"`
	}{}

	// Parse function parses command line parameters into the cfg struct
	cfg, err := config.Parse(&cfg, os.Args)


	---------------------------

	TODO
	- ability to specify units of time in timeouts. For example 10s, 20ms, 1h, 2m, instead of plain digits 10, 20, 30
*/
package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
)

var (
	ErrNoParametersSpecified = errors.New("no command line parameters specified")
	ErrMalformedParameter    = errors.New("malformed command line parameter")
	ErrWrongParameterUsage   = errors.New("wrong parameter usage: -parameterName=value")
	ErrNoParamTagSpecified   = errors.New("`param` is missing")
	ErrTagWithoutValue       = errors.New("struct tag without value")
	ErrTagWithTooManyValues  = errors.New("tag with too many values")
	ErrNoCommandDefined      = errors.New("no command defined")
	ErrUnknownOrNoValue      = errors.New("unknown parameter or no value provided")
	ErrHelpNeeded            = errors.New("help requested")
	ErrInvalidConfigStruct   = errors.New("config sruct expected to be a pointer")
)

// CMDParam describes configuration parameter
type CMDParam struct {
	FieldName string
	FieldType string
	EnvName   string
	Value     string
	Usage     string
}

var commands map[string]CMDParam

func PrintUsage() {

	tw := new(tabwriter.Writer)
	tw.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintln(tw)
	fmt.Fprintf(tw, "%s\t%s\t%s\t %s\n", "Parameter", "Ddata type", "Default value", "Description")
	fmt.Fprintln(tw, "------------------------------------------------------------------------------------------")
	for c, v := range commands {
		fmt.Fprintf(tw, "-%s   \t<%s>   \t%s\t %s\n", c, v.FieldType, v.Value, v.Usage)
	}
	tw.Flush()
}

// Parse parses command line and populates Conf struct fields according to the param tags
func Parse(cfgStruct interface{}, params []string) error {

	tagName := "param"

	commands = map[string]CMDParam{}

	t := reflect.ValueOf(cfgStruct)

	if t.Kind() != reflect.Ptr {
		return ErrInvalidConfigStruct
	}

	t = t.Elem()

	for i := 0; i < t.Type().NumField(); i++ {
		field := t.Type().Field(i)
		tv := field.Tag.Get(tagName)
		if tv == "" {
			//ignore fields without param tag
			continue
		}
		if cmds, ok := parseTagValues(tv); ok {
			cmd, cmdOk := cmds["cmd"]
			env, envOk := cmds["env"]
			if !cmdOk && !envOk {
				return ErrNoCommandDefined
			}
			val := func() string {
				if envOk {
					envVal, ok := os.LookupEnv(env)
					if ok {
						return envVal
					}
				}
				return cmds["default"]
			}

			commands[cmd] = CMDParam{
				FieldName: field.Name,
				Value:     val(),
				EnvName:   env,
				Usage:     cmds["usage"],
				FieldType: field.Type.Kind().String(),
			}
		}
	}
	pmap, err := parseCommandLineParameters(params)
	if err != nil {
		return err
	}
	// assign populate config struct
	// command line parameter takes precedence over environment variable
	// 1) command line parameter is used if specified
	// 2) If there is no command line parameter specified, corresponding environment variable's gonna be used instead
	// 3) if neither command line parameter nor environment vriable is specified defaults are used
	//tv := t.Elem()
	for c, v := range commands {
		f := t.FieldByName(commands[c].FieldName)
		cmdVal, cmdOK := pmap[c]

		switch f.Kind() {
		case reflect.String:
			if cmdOK {
				f.SetString(cmdVal)
			} else {
				f.SetString(v.Value)
			}
		case reflect.Int, reflect.Int64:
			//parse int first
			var val string
			if cmdOK {
				val = cmdVal
			} else {
				val = v.Value
			}
			r, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("%s is not a numeric value", val)
			}
			f.SetInt(int64(r))
		}

	}
	return nil
}

func parseTagValues(tv string) (map[string]string, bool) {
	res := map[string]string{}

	cmdpairs := strings.Split(strings.Trim(tv, ","), ",")
	for _, v := range cmdpairs {
		if cmd, val, ok := splitCommand(v); ok {
			res[cmd] = val
		}
	}
	return res, len(res) != 0
}

func splitCommand(cmd string) (string, string, bool) {
	tokens := strings.Split(strings.Trim(cmd, "="), "=")
	if len(tokens) < 2 {
		return "", "", false
	}
	return tokens[0], tokens[1], true
}

func parseCommandLineParameters(p []string) (map[string]string, error) {

	res := map[string]string{}
	plen := len(p)
	if plen < 2 {
		return nil, nil
	}
	for _, pv := range p[1:] {
		if pv == "-h" || pv == "--help" {
			return nil, ErrHelpNeeded
		}
		if pv[0] != '-' {
			return nil, ErrMalformedParameter
		}
		if c, v, ok := splitCommand(pv); ok {
			res[strings.Trim(c, "-")] = v
		} else {
			return nil, ErrUnknownOrNoValue
		}
	}
	return res, nil
}
