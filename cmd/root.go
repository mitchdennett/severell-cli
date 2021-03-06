package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "severell-cli",
		Short: "The Severell Framework CLI Tool",
		Long: `The Severell CLI tool empowers developers to quickly scaffold applications`,
	}
	config	foo
)

var Verbose bool

// Execute executes the root command.
func Execute() error {
	file, _ := ioutil.ReadFile("severell.json")
	_ = json.Unmarshal([]byte(file), &config)

	for _, comm := range config.Commands {
		clazz := comm
		var args cobra.PositionalArgs
		if comm.NumArgs > -1 {
			args = cobra.ExactArgs(comm.NumArgs)
		}
		var commandImported = &cobra.Command{
			Use:   comm.Command,
			Short: comm.Description,
			Long:  comm.Description,
			Args: 	args,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Compiling Command...")
				flagSet := cmd.Flags()
				var argsSlice []string
				for _, flag := range clazz.Flags {
					acFlag := flagSet.Lookup(flag.Flag)
					if acFlag.Changed {
						s := fmt.Sprintf("-%s=%s", acFlag.Name, acFlag.Value)
						argsSlice = append(argsSlice, s)
					}
				}

				argsToPass := fmt.Sprintf("-Dexec.args=%s args=%s flags=%s", clazz.Class,strings.Join(args, ","), strings.Join(argsSlice, ","))
				cmdCon := exec.Command("mvn","-q", "-T", "1C", "compile", "exec:java", `-Dexec.mainClass=` + config.Package + `.commands.Commander`, argsToPass)

				var b bytes.Buffer
				if Verbose {
					cmdCon.Stdout = &b
					cmdCon.Stderr = &b
				}

				err := cmdCon.Run()
				if err != nil {
					fmt.Println("")
					if Verbose {
						log.Fatalf("Unable to run command. \n%s", string(b.Bytes()))
					} else {
						log.Fatalf("Unable to run command. Run with -v to see underlying error.")
					}
				}
			},
		}

		for _, flag := range clazz.Flags {
			commandImported.Flags().StringVarP(&flag.Value, flag.Flag, flag.Flag, "", flag.Description)
		}

		rootCmd.AddCommand(commandImported)
	}
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	cobra.OnInitialize(initConfig)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}