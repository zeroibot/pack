// Package root contains root application functions
package root

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/lang"
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
type CmdHandler[A any] = func(A, []string)

type CmdConfig[A any] struct {
	Command   string
	MinParams int
	Docs      string
	Handler   CmdHandler[A]
}

type cmdTask[A any] interface {
	CmdHandler() CmdHandler[A]
}

// NewCommand creates a new CmdConfig
func NewCommand[A any](command string, minParams int, docs string, handler CmdHandler[A]) *CmdConfig[A] {
	return new(CmdConfig[A]{command, minParams, docs, handler})
}

// NewCommandTask creates a new CmdConfig
func NewCommandTask[A any](command string, minParams int, docs string, task cmdTask[A]) *CmdConfig[A] {
	return NewCommand[A](command, minParams, docs, task.CmdHandler())
}

// NewCommandMap creates a new map of command to CmdConfigs
func NewCommandMap[A any](cfgs ...*CmdConfig[A]) map[string]*CmdConfig[A] {
	commands := make(map[string]*CmdConfig[A])
	for _, cfg := range cfgs {
		commands[cfg.Command] = cfg
	}
	return commands
}

// MainLoop is the main REPL of the root application
func MainLoop[A any](app A, cmdMap map[string]*CmdConfig[A], onExit func()) {
	var err error
	var line, command string
	var params []string

	fmt.Println("Commands:", len(cmdMap))
	fmt.Printf("Root: type `%s` for list of commands, `%s` to close\n", cmdHelp, cmdExit)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n> ")
		line, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		command, params = getCommandParams(cmdMap, line)
		if command == "" {
			continue
		}
		switch command {
		case cmdExit:
			if onExit != nil {
				onExit()
			}
			return
		case cmdHelp:
			if len(params) == 0 {
				command = allCommands
			} else {
				command = params[0]
			}
			displayHelp(cmdMap, command)
		case cmdSearch:
			keyword := params[0]
			searchCommand(cmdMap, keyword)
		default:
			cmdMap[command].Handler(app, params)
		}
	}
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

// validateCommandParams checks if the command exists and the parameters meet the min parameter count
func validateCommandParams[A any](cmdMap map[string]*CmdConfig[A], command string, params []string) error {
	if command == cmdExit || command == cmdHelp {
		return nil
	}
	if command == cmdSearch {
		return lang.Ternary(len(params) < 1, errInvalidParamCount, nil)
	}
	cfg, ok := cmdMap[command]
	if !ok {
		return errInvalidCommand
	}
	if len(params) < cfg.MinParams {
		return errInvalidParamCount
	}
	return nil
}

// getCommandParams gets the command and parameters from the line
func getCommandParams[A any](cmdMap map[string]*CmdConfig[A], line string) (string, []string) {
	if strings.TrimSpace(line) == "" {
		fmt.Println(getHelp)
		return "", nil
	}
	args := str.SpaceSplit(line)
	command, params := strings.ToLower(args[0]), args[1:]
	err := validateCommandParams(cmdMap, command, params)
	if err != nil {
		fmt.Println("Error:", err)
		if errors.Is(err, errInvalidCommand) {
			fmt.Println(getHelp)
		} else if errors.Is(err, errInvalidParamCount) {
			displayHelp(cmdMap, command)
		}
		return "", nil
	}
	return command, params
}

// displayHelp displays the help list
func displayHelp[A any](cmdMap map[string]*CmdConfig[A], targetCommand string) {
	targetCommand = strings.ToLower(targetCommand)
	if _, ok := cmdMap[targetCommand]; !ok && targetCommand != allCommands && !slices.Contains(helpSkipCommands, targetCommand) {
		fmt.Println("Error: unknown command: ", targetCommand)
		fmt.Println(getHelp)
		return
	}
	fmt.Println("Usage: <command> <params>")
	fmt.Println("\nCommands and params:")

	for _, command := range dict.SortedKeys(cmdMap) {
		if slices.Contains(helpSkipCommands, command) {
			continue
		}
		cfg := cmdMap[command]
		if targetCommand == allCommands || targetCommand == command {
			fmt.Printf("%-30s\t%s\n", command, cfg.Docs)
		}
	}
}

// searchCommand searches for command keyword
func searchCommand[A any](cmdMap map[string]*CmdConfig[A], keyword string) {
	keyword = strings.ToLower(keyword)
	commands := dict.SortedKeys(cmdMap)
	if keyword == allCommands {
		stems := ds.NewSet[string]()
		for _, command := range commands {
			if slices.Contains(helpSkipCommands, command) {
				continue
			}
			stem := str.CleanSplit(command, cmdGlue)[0]
			stems.Add(stem)
		}
		heads := stems.Items()
		slices.Sort(heads)
		for _, head := range heads {
			fmt.Println(head)
		}
	} else {
		for _, command := range commands {
			if slices.Contains(helpSkipCommands, command) {
				continue
			}
			if strings.Contains(command, keyword) {
				cfg := cmdMap[command]
				fmt.Printf("%-30s\t%s\n", command, cfg.Docs)
			}
		}
	}
}
