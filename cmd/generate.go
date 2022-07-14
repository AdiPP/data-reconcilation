/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/AdiPP/recon-general/processor"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate {proxypath} {sourcepath} {destinationpath}",
	Short: "Generate Reconciliation Report",
	Long:  `Generate reconciliation report from proxy and source data.`,
	Run: func(cmd *cobra.Command, args []string) {
		proxypath := ""
		sourcepath := ""
		destinationpath := ""

		for i, arg := range args {
			if i == 0 {
				proxypath = arg
			} else if i == 1 {
				sourcepath = arg
			} else if i == 2 {
				destinationpath = arg
			}
		}

		err := validatePaths(proxypath, sourcepath, destinationpath)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = processor.ProcessCSV(sourcepath, proxypath, destinationpath)

		if err != nil {
			fmt.Printf("failed to generate report, with error: %s \n", err.Error())
			os.Exit(1)
		}

		fmt.Println("succesfully generate report")

		err = processor.ProcessSummary(sourcepath, proxypath, destinationpath)

		if err != nil {
			fmt.Printf("Failed to generate report summary, with error: %s \n", err.Error())
			os.Exit(1)
		}

		fmt.Println("succesfully generate report summary")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func validatePaths(proxypath string, sourcepath string, destionationpath string) error {
	if proxypath == "" {
		return errors.New("proxypath is required")
	}

	if sourcepath == "" {
		return errors.New("sourcepath is required")
	}

	if destionationpath == "" {
		return errors.New("destionationpath is required")
	}

	return nil
}
