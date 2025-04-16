/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/gerrowadat/looperutil/database"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var (
	setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set an attribute on a memory slot",
		Long: `This command takes a memory slot number and produces an updated xml spec
that makes the changes specified.`,
		Run: func(cmd *cobra.Command, args []string) {
			doSet(cmd, args)
		},
	}
	xmlOutput string
)

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.PersistentFlags().StringVar(&xmlOutput, "xml-output", "", "Destination for the updated xml spec (default is stdout)")
}

func doSet(_ *cobra.Command, args []string) {
	if len(args) != 3 {
		fmt.Println("Usage: set <slot> <attribute> <value>")
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

	// 'Name' is a special case, because it's a string.
	if args[1] == "Name" {
		fmt.Printf("Name: [%v] -> [%v]\n", mem.Name.String(), args[2])
		err := mem.SetNameFromString(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err = mem.SetAttributeByName(args[1], args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	// db is now updated to our detination state.
	xml, err := db.ToXML()
	if err != nil {
		fmt.Println(err)
		return
	}
	if xmlOutput == "" {
		fmt.Println(xml)
	} else {
		err = database.WriteXML(xmlOutput, xml)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
