/*
Copyright © 2024 Dmitrii Petrov <dpetroff@gmail.com>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/can3p/slapper/pkg/runner"
	"github.com/spf13/cobra"
)

type arrayFlags []string

func (i *arrayFlags) Type() string {
	return "stringSlice"
}

func (i *arrayFlags) String() string {
	return fmt.Sprintf("+%v", *i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func runCommand() *cobra.Command {
	var headerFlags arrayFlags

	out := &cobra.Command{
		Use:   "run",
		Short: "Run a test",
		Long:  `Run a load test`,
		RunE: func(cmd *cobra.Command, args []string) error {
			workers, err := cmd.Flags().GetUint("workers")

			if err != nil {
				return err
			}

			timeout, err := cmd.Flags().GetDuration("timeout")

			if err != nil {
				return err
			}

			targets, err := cmd.Flags().GetString("targets")

			if err != nil {
				return err
			}

			base64body, err := cmd.Flags().GetBool("base64body")

			if err != nil {
				return err
			}

			rate, err := cmd.Flags().GetUint64("rate")

			if err != nil {
				return err
			}

			minY, err := cmd.Flags().GetDuration("minY")

			if err != nil {
				return err
			}

			maxY, err := cmd.Flags().GetDuration("maxY")

			if err != nil {
				return err
			}

			return runner.Run(workers, timeout, targets, base64body, rate, minY, maxY, headerFlags)
		},
	}

	out.Flags().Uint("workers", 8, "Number of workers")
	out.Flags().Duration("timeout", 30*time.Second, "Requests timeout")
	out.Flags().String("targets", "", "Targets file")
	out.Flags().Bool("base64body", false, "Bodies in targets file are base64-encoded")
	out.Flags().Uint64("rate", 50, "Requests per second")
	out.Flags().Duration("minY", 0, "min on Y axe (default 0ms)")
	out.Flags().Duration("maxY", 100*time.Millisecond, "max on Y axe")
	out.Flags().VarP(&headerFlags, "header", "H", "HTTP header 'key: value' set on all requests. Repeat for more than one header.")

	return out
}

func init() {
	rootCmd.AddCommand(runCommand())
}
