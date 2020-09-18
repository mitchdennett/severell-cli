package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"log"
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

	var b bytes.Buffer
	if Verbose {
		cmdCon.Stdout = &b
		cmdCon.Stderr = &b
	}


	err := cmdCon.Run()
	if err != nil {
		fmt.Println("")
		if Verbose {
			log.Fatalf("Unable to create project. \n%s", string(b.Bytes()))
		} else {
			log.Fatalf("Unable to create project. Run with -v to see underlying error.")
		}
	}
}