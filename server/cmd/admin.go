/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent/user"
	"os"

	"github.com/spf13/cobra"
)

// adminCmd represents the admin command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin commands",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var makeAdminCmd = &cobra.Command{
	Use:   "make",
	Short: "Gives a user admin privileges",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		role := user.RoleAdmin
		if cmd.Flag("remove").Value.String() == "true" {
			role = user.RoleUser
		}
		println(cmd.Flag("remove").Value.String())
		user, err := database.NewDefaultDbManagers().UserManager.UpdateRole(context.Background(), args[0], role)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("User %s|%s|%s is now an `%v`\n", user.DiscordUsername, user.TelegramUsername, user.WalletAddress, user.Role)
	},
}

var listAdminsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all admins",
	Run: func(cmd *cobra.Command, args []string) {
		admins, err := database.NewDefaultDbManagers().UserManager.QueryAdmins(context.Background())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		text := "Admins:\n"
		for _, admin := range admins {
			text += fmt.Sprintf("%s|%s|%s\n", admin.DiscordUsername, admin.TelegramUsername, admin.WalletAddress)
		}
		fmt.Print(text)
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)
	adminCmd.AddCommand(makeAdminCmd)
	adminCmd.AddCommand(listAdminsCmd)

	makeAdminCmd.PersistentFlags().Bool("remove", false, "Remove admin privileges")
}
