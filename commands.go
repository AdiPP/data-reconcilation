package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var reconcileCmd = &cobra.Command{
	Use:   "reconcile {proxyfilepath} {sourcefilepath} {destinationdir} {startdate} {enddate}",
	Short: "Reconcile proxy file and source file",
	Long:  `Reconcile and generate reconciliation report from proxy file and source file.`,
	Run: func(cmd *cobra.Command, args []string) {
		var proxyfilepath string
		var sourcefilepath string
		var destinationdir string
		var startdatestring string
		var enddatestring string

		for i, v := range args {
			switch i {
			case 0:
				proxyfilepath = v
			case 1:
				sourcefilepath = v
			case 2:
				destinationdir = v
			case 3:
				startdatestring = v
			case 4:
				enddatestring = v
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
			panic(err.Error())
		}

		reconOpt, err := getReconcileOption(startdatestring, enddatestring)

		if err != nil {
			panic(err.Error())
		}

		err = ReconcileReportFile(reconOpt, destinationdir)

		if err != nil {
			panic(err.Error())
		}

		err = ReconcileReportSummaryFile(reconOpt, destinationdir)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("succesfully")
	},
}

func init() {
	rootCmd.AddCommand(reconcileCmd)
}

func getReconcileOption(startdatestring string, enddatestring string) (ReconcileOption, error) {
	var startdate time.Time
	var enddate time.Time

	if startdatestring != "" {
		date, err := convertDateStringToDate(startdatestring)

		if err != nil {
			return ReconcileOption{}, err
		}

		startdate = date
	}

	if enddatestring != "" {
		date, err := convertDateStringToDate(enddatestring)

		if err != nil {
			return ReconcileOption{}, err
		}

		enddate = date
	}

	return ReconcileOption{
		StartDate: startdate,
		EndDate:   enddate,
	}, nil
}

func convertDateStringToDate(dateString string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateString)

	if err != nil {
		return time.Time{}, err
	}

	return date, nil
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
