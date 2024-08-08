/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/454270186/mydocker/cmd/handler"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Create a container with namespace and cgroups limit",
	Run: handler.RunCmdHandler,
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().BoolVarP(&handler.IsTTY, "tty", "t", false, "enable tty")

	runCmd.Flags().StringVarP(&handler.MemoryLimit, "mem", "m", "", "memory limit")
	runCmd.Flags().StringVarP(&handler.CpusetLimit, "cpuset", "c", "", "cpuset limit")
	runCmd.Flags().IntVarP(&handler.CpuLimit, "cpu", "C", 0, "cpu limit")

	runCmd.Flags().StringVarP(&handler.Volume, "volume", "v", "", "volume,e.g.: -v /ect/conf:/etc/conf")
}
