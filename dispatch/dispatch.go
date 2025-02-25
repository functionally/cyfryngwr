package dispatch

import (
    "bytes"
    "fmt"
    "log"

    "github.com/google/shlex"
    "github.com/spf13/cobra"
)

var (
    version   = "dev"
    gitCommit = "none"
)

func Run(input string) (string, error) {

    var result bytes.Buffer

    var rootCmd = &cobra.Command{
        Use:   "",
        Short: "Cyfryngwr agent",
        Run: func(cmd *cobra.Command, args []string) {
	    result.WriteString("Send `--help` for usage.")
        },
    }
    rootCmd.SetOut(&result)

    var versionCmd = &cobra.Command{
        Use:   "version",
        Short: "Reply with version information",
        Run: func(cmd *cobra.Command, args []string) {
	  result.WriteString(fmt.Sprintf("Cyfryngwr %s (%s)", version, gitCommit))
        },
    }
    rootCmd.AddCommand(versionCmd)

    // Define the 'greet' subcommand.
    var name string
    var greetCmd = &cobra.Command{
        Use:   "greet",
        Short: "Greet the user",
        Run: func(cmd *cobra.Command, args []string) {
            if name == "" {
                name = "World"
            }
            result.WriteString(fmt.Sprintf("Hello, %s!\n", name))
        },
    }
    greetCmd.Flags().StringVarP(&name, "name", "n", "", "Name to greet")
    rootCmd.AddCommand(greetCmd)

    // Example input string containing quotes and escapes.
//  input := `greet --name "Alice Bob"`

    // Use shlex to parse the input string into an array of arguments.
    args, err := shlex.Split(input)
    if err != nil {
        log.Fatalf("Failed to split input: %v", err)
    }

    // Set the parsed arguments for Cobra instead of using os.Args.
    rootCmd.SetArgs(args)

    // Execute the command.
    if err := rootCmd.Execute(); err != nil {
        fmt.Println("Error:", err)
    }
    
    return result.String(), nil
}
