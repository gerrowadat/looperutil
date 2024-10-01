package cmd

import (
	"fmt"

	"github.com/gerrowadat/looperutil/database"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List information about memory slots",
	Run: func(cmd *cobra.Command, args []string) {
		doLs(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func doLs(_ *cobra.Command, args []string) {
	db, err := database.LoadMemoryFile(memoryFile)
	if err != nil {
		fmt.Println("Error loading memory file:", err)
		return
	}
	if len(args) == 0 {
		for _, mem := range db.Mem {
			fmt.Printf("%v: [%v]\n", mem.Number(), mem.Name.String())
		}
		return
	}
	if len(args[0]) != 2 {
		fmt.Println("Invalid memory slot number(01 to 99)")
		return
	}
	slot := db.GetMemorySlotByNumber(args[0])
	if slot == nil {
		fmt.Printf("Memory slot not found: %v\n", args[0])
		return
	}
	fmt.Print(slot.Describe())
}
