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
	Use:   "info",
	Short: "Info commands",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.AddCommand(listUnhandledMsgTypesCmd)
}
