package reset

import (
	"github.com/devspace-cloud/devspace/pkg/devspace/cloud"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/spf13/cobra"
)

type keyCmd struct {
	Provider string
}

func newKeyCmd() *cobra.Command {
	cmd := &keyCmd{}

	keyCmd := &cobra.Command{
		Use:   "key",
		Short: "Resets a cluster key",
		Long: `
#######################################################
############### devspace reset key ####################
#######################################################
Resets a key for a given cluster. Useful if the key 
cannot be obtained anymore. Needs cluster access scope

Examples:
devspace reset key my-cluster
#######################################################
	`,
		Args: cobra.ExactArgs(1),
		Run:  cmd.RunResetkey,
	}

	keyCmd.Flags().StringVar(&cmd.Provider, "provider", "", "The cloud provider to use")

	return keyCmd
}

// RunResetkey executes the reset key command logic
func (cmd *keyCmd) RunResetkey(cobraCmd *cobra.Command, args []string) {
	// Check if user has specified a certain provider
	var cloudProvider *string
	if cmd.Provider != "" {
		cloudProvider = &cmd.Provider
	}

	// Get provider
	provider, err := cloud.GetProvider(cloudProvider, log.GetInstance())
	if err != nil {
		log.Fatal(err)
	}

	// Reset the key
	err = provider.ResetKey(args[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Donef("Successfully reseted key for cluster %s", args[0])
}
