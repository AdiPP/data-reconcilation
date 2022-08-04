package command

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AdiPP/reconciliation/pkg/exporting"
	"github.com/AdiPP/reconciliation/pkg/storage/memory"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "reconciliation",
	Short: "",
	Long:  ``,
}

func ExecuteConsole() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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

		var exporter exporting.Service

		s, err := memory.NewStorageFromFile(
			proxyfilepath,
			sourcefilepath,
		)

		if err != nil {
			panic(err.Error())
		}

		exporter = exporting.NewService(s)

		reportOpt, err := getReconcileOption(startdatestring, enddatestring)

		if err != nil {
			panic(err.Error())
		}

		err = exporter.GenerateReportFile(reportOpt, destinationdir)

		if err != nil {
			panic(err.Error())
		}

		err = exporter.GenerateSummaryReportFile(reportOpt, destinationdir)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("reconciliation run successfully")
	},
}

func init() {
	rootCmd.AddCommand(reconcileCmd)
}

func getReconcileOption(startdatestring string, enddatestring string) (exporting.ReportOption, error) {
	var startdate time.Time
	var enddate time.Time

	if startdatestring != "" {
		date, err := convertDateStringToDate(startdatestring)

		if err != nil {
			return exporting.ReportOption{}, err
		}

		startdate = date
	}

	if enddatestring != "" {
		date, err := convertDateStringToDate(enddatestring)

		if err != nil {
			return exporting.ReportOption{}, err
		}

		enddate = date
	}

	return exporting.ReportOption{
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
