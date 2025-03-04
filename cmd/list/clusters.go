package list

import (
	"strconv"

	cloudpkg "github.com/devspace-cloud/devspace/pkg/devspace/cloud"
	"github.com/devspace-cloud/devspace/pkg/util/log"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

type clustersCmd struct {
	Provider string
	All      bool
}

func newClustersCmd() *cobra.Command {
	cmd := &clustersCmd{}

	clustersCmd := &cobra.Command{
		Use:   "clusters",
		Short: "Lists all connected clusters",
		Long: `
#######################################################
############## devspace list clusters #################
#######################################################
List all connected user clusters

Example:
devspace list clusters
#######################################################
	`,
		Args: cobra.NoArgs,
		Run:  cmd.RunListClusters,
	}

	clustersCmd.Flags().StringVar(&cmd.Provider, "provider", "", "Cloud Provider to use")
	clustersCmd.Flags().BoolVar(&cmd.All, "all", false, "Show all available clusters including hosted DevSpace cloud clusters")

	return clustersCmd
}

// RunListClusters executes the "devspace list clusters" functionality
func (cmd *clustersCmd) RunListClusters(cobraCmd *cobra.Command, args []string) {
	// Check if user has specified a certain provider
	var cloudProvider *string
	if cmd.Provider != "" {
		cloudProvider = &cmd.Provider
	}

	// Get provider
	provider, err := cloudpkg.GetProvider(cloudProvider, log.GetInstance())
	if err != nil {
		log.Fatalf("Error getting cloud context: %v", err)
	}

	log.StartWait("Retrieving clusters")
	clusters, err := provider.GetClusters()
	if err != nil {
		log.Fatalf("Error retrieving clusters: %v", err)
	}
	log.StopWait()

	headerColumnNames := []string{
		"ID",
		"Name",
		"Owner",
		"Created",
	}

	values := [][]string{}

	for _, cluster := range clusters {
		owner := ""
		createdAt := ""
		if cluster.Owner != nil {
			owner = cluster.Owner.Name

			if cluster.CreatedAt != nil {
				createdAt = *cluster.CreatedAt
			}
		} else if cmd.All == false {
			continue
		}

		values = append(values, []string{
			strconv.Itoa(cluster.ClusterID),
			cluster.Name,
			owner,
			createdAt,
		})
	}

	if len(values) > 0 {
		log.PrintTable(log.GetInstance(), headerColumnNames, values)
	} else {
		log.Infof("No clusters found. You can connect a cluster with `%s`", ansi.Color("devspace connect cluster", "white+b"))
	}
}
