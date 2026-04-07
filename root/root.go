package root

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/str"
	"golang.org/x/term"
)

const (
	allCommands string = "*"
	cmdHelp     string = "help"
	cmdExit     string = "exit"
	cmdSearch   string = "cmd"
	cmdGlue     string = "/"
)

var (
	errInvalidCommand    = fmt.Errorf("invalid command")
	errInvalidParamCount = fmt.Errorf("invalid param count")
	getHelp              = fmt.Sprintf("Type `%s` for list of commands, `%s <keyword>` to search for command", cmdHelp, cmdSearch)
	helpSkipCommands     = []string{cmdHelp, cmdExit, cmdSearch}
)

// CmdHandler takes in a list of string parameters
type CmdHandler = func([]string)

type CmdConfig struct {
	Command   string
	MinParams int
	Docs      string
	Handler   CmdHandler
}

// NewCommand creates a new CmdConfig
func NewCommand(command string, minParams int, docs string, handler CmdHandler) *CmdConfig {
	return new(CmdConfig{command, minParams, docs, handler})
}

// NewCommandMap creates a new map of command to CmdConfigs
func NewCommandMap(cfgs ...*CmdConfig) map[string]*CmdConfig {
	commands := make(map[string]*CmdConfig)
	for _, cfg := range cfgs {
		commands[cfg.Command] = cfg
	}
	return commands
}

// ParamsMap gets the key=value map from the parameters list
func ParamsMap(params []string, required []string, optional []string) (dict.Strings, error) {
	if required == nil {
		required = make([]string, 0)
	}
	if optional == nil {
		optional = make([]string, 0)
	}
	paramsMap := make(dict.Strings)
	for _, param := range params {
		parts := str.CleanSplitN(param, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		if !slices.Contains(required, key) && !slices.Contains(optional, key) {
			continue
		}
		paramsMap[key] = value
	}
	for _, key := range required {
		if _, ok := paramsMap[key]; !ok {
			return nil, fail.MissingParams
		}
	}
	return paramsMap, nil
}

// Authenticate performs authentication for the Root account in the command-line app
func Authenticate(authFn func(string) error) error {
	fmt.Print("Enter password: ")
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()
	password := strings.TrimSpace(string(pwd))
	err = authFn(password)
	if err != nil {
		return fmt.Errorf("root authentication failed: %w", err)
	}
	return nil
}
