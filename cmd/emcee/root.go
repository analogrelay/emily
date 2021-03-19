package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type emcee struct {
}

func init() {
	var (
		emcee = &emcee{}
		cmd   = &cobra.Command{
			Use:   "emcee",
			Short: "Emily Compiler (emcee)",
			Long:  `Compile programs in the Emily language.`,
			Args:  cobra.MinimumNArgs(1),
			RunE:  emcee.Run,
		}
	)
	rootCmd = cmd
}

func (e *emcee) Run(cmd *cobra.Command, args []string) error {
	panic("not implemented")
}
