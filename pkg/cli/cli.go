package cli

import (
	"flag"
	"fmt"
	"os"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Action      func(args []string) error
	Subcommands map[string]*Command
	Flags       *flag.FlagSet
}

// CLI represents the main CLI application
type CLI struct {
	Name        string
	Version     string
	Description string
	Commands    map[string]*Command
}

// NewCLI creates a new CLI instance
func NewCLI(name, version, description string) *CLI {
	return &CLI{
		Name:        name,
		Version:     version,
		Description: description,
		Commands:    make(map[string]*Command),
	}
}

// AddCommand adds a command to the CLI
func (c *CLI) AddCommand(cmd *Command) {
	c.Commands[cmd.Name] = cmd
}

// Run executes the CLI with the given arguments
func (c *CLI) Run(args []string) error {
	if len(args) < 1 {
		c.PrintUsage()
		return nil
	}

	cmdName := args[0]

	// Handle global flags
	if cmdName == "help" || cmdName == "--help" || cmdName == "-h" {
		if len(args) > 1 {
			// Help for specific command
			if cmd, ok := c.Commands[args[1]]; ok {
				cmd.PrintUsage()
				return nil
			}
		}
		c.PrintUsage()
		return nil
	}

	if cmdName == "version" || cmdName == "--version" || cmdName == "-v" {
		fmt.Printf("%s version %s\n", c.Name, c.Version)
		return nil
	}

	// Find and execute command
	cmd, ok := c.Commands[cmdName]
	if !ok {
		return fmt.Errorf("unknown command: %s\nRun '%s help' for usage", cmdName, c.Name)
	}

	// Check for subcommands
	if len(args) > 1 && len(cmd.Subcommands) > 0 {
		subCmdName := args[1]

		// Handle subcommand help
		if subCmdName == "help" || subCmdName == "--help" || subCmdName == "-h" {
			cmd.PrintUsage()
			return nil
		}

		if subCmd, ok := cmd.Subcommands[subCmdName]; ok {
			return subCmd.Action(args[2:])
		}
		return fmt.Errorf("unknown subcommand: %s %s\nRun '%s %s help' for usage", cmdName, subCmdName, c.Name, cmdName)
	}

	// Execute command action
	if cmd.Action != nil {
		return cmd.Action(args[1:])
	}

	// If no action and has subcommands, show usage
	if len(cmd.Subcommands) > 0 {
		cmd.PrintUsage()
		return nil
	}

	return fmt.Errorf("command '%s' has no action defined", cmdName)
}

// PrintUsage prints the CLI usage information
func (c *CLI) PrintUsage() {
	fmt.Printf("%s - %s\n\n", c.Name, c.Description)
	fmt.Println("Usage:")
	fmt.Printf("  %s <resource> <action> [options]\n\n", c.Name)
	fmt.Println("Resources:")

	// Print commands in order
	commandOrder := []string{"http", "env", "context", "demo", "version", "help"}
	for _, name := range commandOrder {
		if cmd, ok := c.Commands[name]; ok {
			fmt.Printf("  %-15s %s\n", name, cmd.Description)
		}
	}

	fmt.Println("\nGlobal Options:")
	fmt.Println("  --help, -h      Show help information")
	fmt.Println("  --version, -v   Show version information")
	fmt.Println("\nExamples:")
	fmt.Printf("  %s http run requests.http --env production\n", c.Name)
	fmt.Printf("  %s env list\n", c.Name)
	fmt.Printf("  %s env show development\n", c.Name)
	fmt.Printf("  %s context set --http-file requests.http --env development\n", c.Name)
	fmt.Printf("  %s http run requests.http --save-responses\n", c.Name)
	fmt.Printf("\nRun '%s <resource> help' for more information on a resource.\n", c.Name)
}

// PrintUsage prints the command usage information
func (cmd *Command) PrintUsage() {
	fmt.Printf("%s - %s\n\n", cmd.Name, cmd.Description)

	if len(cmd.Subcommands) > 0 {
		fmt.Println("Available actions:")
		for name, subCmd := range cmd.Subcommands {
			fmt.Printf("  %-15s %s\n", name, subCmd.Description)
		}
		fmt.Printf("\nRun 'postie %s <action> --help' for more information on an action.\n", cmd.Name)
	} else if cmd.Flags != nil {
		fmt.Println("Usage:")
		fmt.Printf("  postie %s [options]\n\n", cmd.Name)
		fmt.Println("Options:")
		cmd.Flags.PrintDefaults()
	}
}

// StringFlag represents a string flag with short and long names
type StringFlag struct {
	Name      string
	ShortName string
	Value     string
	Usage     string
	Required  bool
}

// BoolFlag represents a boolean flag
type BoolFlag struct {
	Name      string
	ShortName string
	Value     bool
	Usage     string
}

// ParseFlags is a helper to parse flags with short and long names
func ParseFlags(args []string, stringFlags []*StringFlag, boolFlags []*BoolFlag) (*flag.FlagSet, error) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	// Define string flags
	for _, sf := range stringFlags {
		fs.StringVar(&sf.Value, sf.Name, "", sf.Usage)
		if sf.ShortName != "" {
			fs.StringVar(&sf.Value, sf.ShortName, "", sf.Usage)
		}
	}

	// Define bool flags
	for _, bf := range boolFlags {
		fs.BoolVar(&bf.Value, bf.Name, false, bf.Usage)
		if bf.ShortName != "" {
			fs.BoolVar(&bf.Value, bf.ShortName, false, bf.Usage)
		}
	}

	// Parse
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// Check required flags
	for _, sf := range stringFlags {
		if sf.Required && sf.Value == "" {
			return nil, fmt.Errorf("required flag --%s not provided", sf.Name)
		}
	}

	return fs, nil
}
