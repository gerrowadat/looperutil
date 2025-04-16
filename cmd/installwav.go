package cmd

import (
	"fmt"

	"github.com/gerrowadat/looperutil/database"
	"github.com/spf13/cobra"
)

// installwavCmd represents the installwav command
var installwavCmd = &cobra.Command{
	Use:   "installwav",
	Short: "Convert and install an audio file to a slot on the looper",
	Run: func(cmd *cobra.Command, args []string) {
		doInstallWav(cmd, args)

	},
}

var clearExisting bool
var looperDataDir string

func init() {
	rootCmd.AddCommand(installwavCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installwavCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	installwavCmd.PersistentFlags().BoolVar(&clearExisting, "clear-existing-wav", false, "Clear existing wav files on the looper first.")
	installwavCmd.PersistentFlags().StringVar(&looperDataDir, "looper-data-dir", "/mnt/BOSS RC-5", "Directory where the looper data is stored.")
}

func doInstallWav(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Printf("Usage: %v <slot> <audio file>\n", cmd.UseLine())
		return
	}
	slot := args[0]
	audioFile := args[1]
	fmt.Printf("Installing %v to slot %v\n", audioFile, slot)

	db, err := database.LoadMemoryFile(looperDataDir + "/ROLAND/DATA/MEMORY1.RC0")
	if err != nil {
		fmt.Println(err)
		return
	}
	mem := db.GetMemorySlotByNumber(slot)
	if mem == nil {
		fmt.Printf("Memory slot not found: %v\n", slot)
		return
	}

	err = database.InstallWav(audioFile, mem, looperDataDir, clearExisting)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Installed %v to slot %v\n", audioFile, slot)
}
