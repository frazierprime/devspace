package cmd

import (
	"fmt"
	"strings"

	"github.com/devspace-cloud/devspace/cmd/add"
	"github.com/devspace-cloud/devspace/cmd/cleanup"
	"github.com/devspace-cloud/devspace/cmd/connect"
	"github.com/devspace-cloud/devspace/cmd/create"
	"github.com/devspace-cloud/devspace/cmd/list"
	"github.com/devspace-cloud/devspace/cmd/remove"
	"github.com/devspace-cloud/devspace/cmd/reset"
	"github.com/devspace-cloud/devspace/cmd/set"
	"github.com/devspace-cloud/devspace/cmd/status"
	"github.com/devspace-cloud/devspace/cmd/update"
	"github.com/devspace-cloud/devspace/cmd/use"
	"github.com/devspace-cloud/devspace/pkg/devspace/upgrade"
	"github.com/devspace-cloud/devspace/pkg/util/analytics/cloudanalytics"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devspace",
	Short: "Welcome to the DevSpace CLI!",
	Long: `DevSpace accelerates developing, deploying and debugging applications with Docker and Kubernetes. Get started by running the init command in one of your projects:

	devspace init`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	version := upgrade.GetVersion()
	defer cloudanalytics.ReportPanics()

	if version != "" {
		rootCmd.Version = upgrade.GetVersion()

		if strings.Contains(upgrade.GetVersion(), "-alpha") == false && strings.Contains(upgrade.GetVersion(), "-beta") == false {
			newerVersion, err := upgrade.CheckForNewerVersion()
			if err == nil && newerVersion != "" {
				log.Warnf("There is a newer version of DevSpace CLI v%s. Run `devspace upgrade` to update the CLI.\n", newerVersion)
			}
		}
	}

	// Execute command
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}

	cloudanalytics.SendCommandEvent(err)
}

func init() {
	// Add sub commands
	rootCmd.AddCommand(add.NewAddCmd())
	rootCmd.AddCommand(cleanup.NewCleanupCmd())
	rootCmd.AddCommand(connect.NewConnectCmd())
	rootCmd.AddCommand(create.NewCreateCmd())
	rootCmd.AddCommand(list.NewListCmd())
	rootCmd.AddCommand(remove.NewRemoveCmd())
	rootCmd.AddCommand(reset.NewResetCmd())
	rootCmd.AddCommand(set.NewSetCmd())
	rootCmd.AddCommand(status.NewStatusCmd())
	rootCmd.AddCommand(use.NewUseCmd())
	rootCmd.AddCommand(update.NewUpdateCmd())

	// Add main commands
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewDevCmd())
	rootCmd.AddCommand(NewBuildCmd())
	rootCmd.AddCommand(NewSyncCmd())
	rootCmd.AddCommand(NewInstallCmd())
	rootCmd.AddCommand(NewPurgeCmd())
	rootCmd.AddCommand(NewUpgradeCmd())
	rootCmd.AddCommand(NewDeployCmd())
	rootCmd.AddCommand(NewEnterCmd())
	rootCmd.AddCommand(NewLoginCmd())
	rootCmd.AddCommand(NewAnalyzeCmd())
	rootCmd.AddCommand(NewLogsCmd())
	rootCmd.AddCommand(NewOpenCmd())
	rootCmd.AddCommand(NewUICmd())
	rootCmd.AddCommand(NewContainerizeCmd())

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Panic(err)
		}

		// Search config in home directory with name ".devspace" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".devspace")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
}
