package reset

import "github.com/spf13/cobra"

// NewResetCmd creates a new cobra command
func NewResetCmd() *cobra.Command {
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Resets an cluster token",
		Long: `
#######################################################
################## devspace reset #####################
#######################################################
	`,
		Args: cobra.NoArgs,
	}

	resetCmd.AddCommand(newKeyCmd())
	resetCmd.AddCommand(newVarsCmd())

	return resetCmd
}
