package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var reconcileCmd = &cobra.Command{
	Use:   "reconcile {proxyfilepath} {sourcefilepath} {destinationdir}",
	Short: "Reconcile proxy file and source file",
	Long:  `Reconcile and generate reconciliation report from proxy file and source file.`,
	Run: func(cmd *cobra.Command, args []string) {
		proxyfilepath := ""
		sourcefilepath := ""
		destinationdir := ""

		for i, arg := range args {
			if i == 0 {
				proxyfilepath = arg
			} else if i == 1 {
				sourcefilepath = arg
			} else if i == 2 {
				destinationdir = arg
			}
		}

		if err := validateArguments(proxyfilepath, sourcefilepath, destinationdir); err != nil {
			panic(err.Error())
		}

		var err error

		db, err = NewStorageFile(
			proxyfilepath,
			sourcefilepath,
		)

		if err != nil {
			log.Fatal(err)
		}

		err = ReconcileCSVFile(destinationdir)

		if err != nil {
			panic(err.Error())
		}

		err = ReconcileSummaryFile(destinationdir)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("succesfully")
	},
}

func init() {
	rootCmd.AddCommand(reconcileCmd)
}

func validateArguments(proxypath string, sourcepath string, destionationpath string) error {
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
