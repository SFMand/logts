package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	fromPath string
	toPath   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "logts",
	Short: "logts - CLI file compression tool",
	Long:  `A CLI log archiving tool built with go, used to process file compression to .tar.gz`,
	RunE:  start,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Printf("Error occurred while running program: %v\n", err)
		os.Exit(1)
	}
}

func start(cmd *cobra.Command, args []string) error {
	return targzComp(&fromPath, &toPath)
}

func init() {
	rootCmd.PersistentFlags().StringVar(&fromPath, "from", "", "path of directory to compress")
	rootCmd.MarkPersistentFlagRequired("from")
	rootCmd.PersistentFlags().StringVar(&toPath, "to", "", "path of directory to send compressed folder (if omitted, will use same directory as \"from\" path)")
}
