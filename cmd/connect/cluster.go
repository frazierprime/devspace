package connect

import (
	"github.com/devspace-cloud/devspace/pkg/devspace/cloud"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

type clusterCmd struct {
	Provider string
}

func newClusterCmd() *cobra.Command {
	cmd := &clusterCmd{}

	clusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "Connects an existing cluster to DevSpace Cloud",
		Long: `
#######################################################
############ devspace connect cluster #################
#######################################################
Connects an existing cluster to DevSpace Cloud.

Examples:
devspace connect cluster 
devspace connect cluster my-cluster
#######################################################
	`,
		Args: cobra.MaximumNArgs(1),
		Run:  cmd.RunConnectCluster,
	}

	clusterCmd.Flags().StringVar(&cmd.Provider, "provider", "", "The cloud provider to use")

	return clusterCmd
}

// RunConnectCluster executes the connect cluster command logic
func (cmd *clusterCmd) RunConnectCluster(cobraCmd *cobra.Command, args []string) {
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

	clusterName := ""
	if len(args) > 0 {
		clusterName = args[0]
	}

	// Connect cluster
	err = provider.ConnectCluster(clusterName)
	if err != nil {
		log.Fatal(err)
	}

	log.Donef("Successfully connected cluster to DevSpace Cloud. You can now run:\n- `%s` to create a new space\n- `%s` to open the ui and configure cluster access and users\n- `%s` to list all connected clusters", ansi.Color("devspace create space [NAME]", "white+b"), ansi.Color("devspace ui", "white+b"), ansi.Color("devspace list clusters", "white+b"))
}