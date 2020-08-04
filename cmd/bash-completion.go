package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var bashCompletionCmd = &cobra.Command{
	Use:   "bash-completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(sigsci bash-completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(goldap bash-completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(bashCompletionCmd)
}
