package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(load)
}

var load = &cobra.Command{
	Use:   "load",
	Short: "Load Commands From Project",
	Long:  `Load Commands From Project`,
	Run: func(cmd *cobra.Command, args []string) {

		loadCommands(nil, nil)
	},
}

func loadCommands(dir *string, packageName *string) {
	if packageName == nil {
		packageName = &config.Package
	}

	cmdCon := exec.Command("mvn","-q", "compile", "exec:java", `-Dexec.mainClass=` + *packageName + `.commands.Commander`, `-Dexec.args=load`)

	if dir != nil {
		cmdCon.Dir = *dir
	}

	cmdCon.Stdout = os.Stdout
	cmdCon.Stderr = os.Stderr

	err := cmdCon.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}