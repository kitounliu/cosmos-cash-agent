package main

import (
	"github.com/spf13/cobra"

	"github.com/allinbits/cosmos-cash-agent/cmd/elesto-agent/startcmd"
	"github.com/hyperledger/aries-framework-go/pkg/common/log"
)

// This is an application which starts Aries agent controller API on given port.
func main() {
	rootCmd := &cobra.Command{
		Use: "elesto-agent",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	logger := log.New("aries-framework/agent-rest")

	startCmd, err := startcmd.Cmd(&startcmd.HTTPServer{})
	if err != nil {
		logger.Fatalf(err.Error())
	}

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatalf("Failed to run elesto-rest: %s", err)
	}
}
