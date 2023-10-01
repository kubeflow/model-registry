package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "model-registry",
	Short: "A go server for ml-metadata",
	Long: `The model-registry is a gRPC server that stores metadata 
for ML applications. 

It's based on the ml-metadata project that provides a python client library 
for ML applications to record metadata about metadata such as Artifacts, 
Executions and Contexts. 
This go server is an alternative to the CPP gRPC service provided by the 
ml-metadata project. It's meant to provide extra features such as loading 
custom metadata libraries, exposing a higher level GraphQL API, RBAC, etc.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		glog.Exitf("error: %v", err)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.model-registry.yaml)")

	// default to logging to stderr
	_ = flag.Set("logtostderr", "true")
	// also add standard glog flags
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		// Search config in home directory with name ".model-registry" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".model-registry")
	}

	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		glog.Info("using config file: ", viper.ConfigFileUsed())
	} else {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		ok := errors.As(err, &configFileNotFoundError)
		// ignore if it's a file not found error for default config file
		if !(cfgFile == "" && ok) {
			return fmt.Errorf("reading config %s: %v", viper.ConfigFileUsed(), err)
		}
	}

	// bind flags to config
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	var err error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		name := f.Name
		if err == nil && !f.Changed && viper.IsSet(name) {
			value := viper.Get(name)
			err = cmd.Flags().Set(name, fmt.Sprintf("%v", value))
		}
	})

	return err
}
