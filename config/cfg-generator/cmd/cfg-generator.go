package cmd

import "github.com/spf13/cobra"

func Execute() {
	if err := cfgGeneratorCmd.Execute(); err != nil {
		panic(err)
	}
}

var cfgGeneratorCmd = &cobra.Command{
	Use:   "cfg-generator",
	Short: "Generate configuration files for the NFs",
	Run:   cfgGeneratorRun,
}

func init() {

}

func cfgGeneratorRun(cmd *cobra.Command, args []string) {

}
