/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/gerrowadat/looperutil/database"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Individual attributes from your database",
	Run: func(cmd *cobra.Command, args []string) {
		doGet(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func doGet(_ *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("Usage: get <slot> <attribute>")
		return
	}
	db, err := database.LoadMemoryFile(memoryFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	mem := db.GetMemorySlotByNumber(args[0])
	if mem == nil {
		fmt.Printf("Memory slot not found: %v\n", args[0])
		return
	}

	if args[1] == "Name" {
		fmt.Println(mem.Name.String())
	} else {
		val, err := mem.GetAttributeByName(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(val)
	}

}
