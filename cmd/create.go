package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type command struct {
	Class string
	Description string
	Command string
	NumArgs int
	Flags []flag
}

var ARCHETYPE_GROUP = "com.severell"
var ARCHETYPE_ARTIFACT = "severell-archetype"
var ARCHETYPE_VERSION = "0.0.1-SNAPSHOT"

type flag struct {
	Flag string
	Description string
	Value string
}

type foo struct {
	AppName string
	Package string
	BasePackage string
	GroupId string
	Commands []command
}

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new Severell project",
	Long:  `Create a new Severell project`,
	Args: 	cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		//Getting Base Package
		fmt.Print("Base Package (Group ID): ")
		basePackage, _ := reader.ReadString('\n')
		basePackage = strings.Replace(basePackage, string('\n'), "", 1)

		var appName string
		if len(args) > 0 {
			appName = args[0]
		} else {
			//Getting Artifact Id
			fmt.Print("Application Name (Artifact ID): ")
			appName, _ = reader.ReadString('\n')
			appName = strings.Replace(appName, string('\n'), "", 1)
		}

		st := "create"
		var wg sync.WaitGroup

		wg.Add(1)
		go writeToConsole(&st, &wg)

		cmdCon := exec.Command("mvn","-B", "archetype:generate",
			fmt.Sprintf("-DarchetypeGroupId=%s", ARCHETYPE_GROUP),
			fmt.Sprintf("-DarchetypeArtifactId=%s", ARCHETYPE_ARTIFACT),
			fmt.Sprintf("-DarchetypeVersion=%s", ARCHETYPE_VERSION),
			fmt.Sprintf("-DgroupId=%s", basePackage),
			fmt.Sprintf("-DartifactId=%s", appName),
			"-Dversion=1.0-SNAPSHOT")

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


		p := basePackage + "." + appName
		loadCommands(&appName, &p)

		st = "stop"

		wg.Wait()
	},
}





func writeToConsole(st *string, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	fmt.Println("")
	fmt.Println("Creating new project")
	time.Sleep(1 * time.Second)
	for{
		if *st == "stop" {
			fmt.Println("")
			fmt.Println("**********************************")
			fmt.Println("Successfully Created Project")
			fmt.Println("**********************************")
			fmt.Println("")
			break
		}

		switch count {
			case 0:
				fmt.Println("Setting Group ID")
			case 1:
				fmt.Println("Setting Artifact ID")
			case 2:
				fmt.Println("Setting Version")
			case 3:
				fmt.Println("Setting Up Project")
			case 4:
				fmt.Println("Loading Commands")
			default:
				if *st == "create" {
					fmt.Print(".")
				}
		}

		count = count + 1
		time.Sleep(1 * time.Second)
	}

}