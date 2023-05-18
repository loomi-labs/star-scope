/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/loomi-labs/star-scope/database"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "chain",
	Short: "Chain commands",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var listHandledMsgTypesCmd = &cobra.Command{
	Use:   "list-handled",
	Short: "List all handled message types",
	Run: func(cmd *cobra.Command, args []string) {
		chains := database.NewDefaultDbManagers().ChainManager.QueryAll(context.Background())
		text := "Infos:\n"
		for _, chain := range chains {
			text += fmt.Sprintf("%s:\n\t%s\n", chain.Name, strings.Join(strings.Split(chain.HandledMessageTypes, ","), "\n\t"))
		}
		fmt.Print(text)
	},
}

var listUnhandledMsgTypesCmd = &cobra.Command{
	Use:   "list-unhandled",
	Short: "List all unhandled message types",
	Run: func(cmd *cobra.Command, args []string) {
		chains := database.NewDefaultDbManagers().ChainManager.QueryAll(context.Background())
		text := "Infos:\n"
		for _, chain := range chains {
			text += fmt.Sprintf("%s:\n\t%s\n", chain.Name, strings.Join(strings.Split(chain.UnhandledMessageTypes, ","), "\n\t"))
		}
		fmt.Print(text)
	},
}

func enableChain(args []string, enable bool) {
	if len(args) == 0 {
		fmt.Println("Missing chain name")
		os.Exit(1)
	}
	chainManager := database.NewDefaultDbManagers().ChainManager
	chains := chainManager.QueryByName(context.Background(), args[0])
	if len(chains) == 0 {
		fmt.Println("Chain not found")
		os.Exit(1)
	}
	chain, err := chainManager.UpdateSetEnabled(context.Background(), chains[0], enable)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Chain %s is now %s\n", chain.Name, func() string {
		if enable {
			return "enabled"
		}
		return "disabled"
	}())
}

var enableChainCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable/disable chain",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		enable, err := cmd.Flags().GetBool("enable")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		enableChain(args, enable)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.AddCommand(listHandledMsgTypesCmd)
	infoCmd.AddCommand(listUnhandledMsgTypesCmd)
	infoCmd.AddCommand(enableChainCmd)
	enableChainCmd.Flags().BoolP("enable", "e", true, "Enable chain")
}
