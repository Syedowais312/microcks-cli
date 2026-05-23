package cmd

import (
	"fmt"
	"os"

	"github.com/microcks/microcks-cli/pkg/config"
	"github.com/microcks/microcks-cli/pkg/connectors"
	"github.com/microcks/microcks-cli/pkg/errors"
	"github.com/spf13/cobra"
)

func NewLogoutCommand(globalClientOpts *connectors.ClientOptions) *cobra.Command {

	logoutCmd := &cobra.Command{
		Use:   "logout CONTEXT",
		Short: "Log out from Microcks",
		Long:  "Log out from Microcks",
		Example: `# Log out from a Microcks server URL
microcks logout http://localhost:8080

# Log out from a named context
microcks logout dev-context`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			context := args[0]
			err := logoutContext(context, globalClientOpts.ConfigPath)
			if err != nil {
				fmt.Printf("Error logging out from '%s': %v\n", context, err)
				os.Exit(1)
			}
			fmt.Printf("Logged out from '%s'\n", context)
		},
	}

	return logoutCmd
}

func logoutContext(context, configPath string) error {
	localCfg, err := config.ReadLocalConfig(configPath)
	errors.CheckError(err)
	if localCfg == nil {
		return fmt.Errorf("Nothing to logout from")
	}

	tokenOwner := context
	resolvedCtx, err := localCfg.ResolveContext(context)
	if err == nil {
		tokenOwner = resolvedCtx.User.Name
	}

	ok := localCfg.RemoveToken(tokenOwner)
	if !ok {
		return fmt.Errorf("Context %s does not exist", context)
	}

	err = config.ValidateLocalConfig(*localCfg)
	if err != nil {
		return fmt.Errorf("Error in loging out: %s", err)
	}

	return config.WriteLocalConfig(*localCfg, configPath)
}
