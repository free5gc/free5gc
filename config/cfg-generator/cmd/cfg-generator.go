// Package cmd implements the cfg-generator command line.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/free5gc/config/cfg-generator/generator"
)

// exitCodesHelp documents the exit codes in --help.
const exitCodesHelp = `Exit codes:
  0  success (global values that could not be applied are reported as
     "[not applied]" on stderr)
  1  input: arguments, file read, YAML syntax
  2  global block validation
  3  applying the global block to NF configs
  4  NF config validation
  5  writing or publishing output
  6  --strict and at least one global value was not applied
On exit codes 1-4 and 6 nothing in the output directory is modified.
Exit 5 during publishing may leave a partial update; the error lists the
files already published.`

// Execute runs the CLI with the process arguments and exits with the
// stage-specific exit code (0 on success; see exitCodesHelp).
func Execute() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

// Run runs the CLI with args and returns the exit code. It never panics on
// user errors, so run.sh can rely on the exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	var (
		opts   generator.Options
		runErr error
	)

	root := &cobra.Command{
		Use:           "cfg-generator",
		Short:         "Generate configuration files for the NFs",
		Long:          "Generate configuration files for the NFs.\n\n" + exitCodesHelp,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			res, err := generator.Run(opts)
			if res != nil {
				// Also set on a --strict failure, so the user sees why.
				reportIssues(stderr, res)
			}
			if err != nil {
				runErr = err
				return err
			}
			changes := len(res.Changes)
			notApplied := len(res.NotApplied)
			for _, c := range res.Changes {
				_, _ = fmt.Fprintln(stdout, c.String())
			}
			if opts.DryRun {
				_, _ = fmt.Fprintf(stdout, "dry run: %d change(s), %d not applied, would write %v to %s\n",
					changes, notApplied, res.Files, opts.OutDir)
			} else {
				_, _ = fmt.Fprintf(stdout, "%d change(s), %d not applied, wrote %v to %s\n",
					changes, notApplied, res.Files, opts.OutDir)
			}
			return nil
		},
	}
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.Flags().StringVarP(&opts.InputPath, "input", "i", "", "integrated input YAML (required)")
	root.Flags().StringVarP(&opts.OutDir, "output", "o", "./config", "output directory")
	root.Flags().BoolVar(&opts.DryRun, "dry-run", false, "validate and report changes without writing files")
	root.Flags().BoolVar(&opts.Strict, "strict", false,
		"fail with exit code 6, writing nothing, if a global value cannot be applied to an NF")
	_ = root.MarkFlagRequired("input")

	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		if runErr == nil {
			// Flag parsing or argument errors from cobra.
			return generator.StageInput.ExitCode()
		}
		return generator.ExitCode(runErr)
	}
	return 0
}

// reportIssues prints warnings and not-applied global values to stderr.
func reportIssues(stderr io.Writer, res *generator.Result) {
	for _, w := range res.Warnings {
		_, _ = fmt.Fprintln(stderr, "[warn] "+w.String())
	}
	for _, n := range res.NotApplied {
		_, _ = fmt.Fprintln(stderr, "[not applied] "+n.String())
	}
}
