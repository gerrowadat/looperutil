package cmd

import (
	"fmt"

	"github.com/gerrowadat/looperutil/database"
	"github.com/spf13/cobra"
)

var xmlnameCmd = &cobra.Command{
	Use:   "xmlname",
	Short: "Print the <NAME> tag for a given memory slot name",
	Run: func(cmd *cobra.Command, args []string) {
		doXmlName(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(xmlnameCmd)
}

func doXmlName(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("PLease specify one argument")
		return
	}
	slotname, err := database.SlotNameFromString(args[0])
	if err != nil {
		fmt.Println("Error parsing slot name:", err)
		return
	}
	fmt.Println(slotname.XmlString())
}
