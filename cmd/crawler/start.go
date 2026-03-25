package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the web crawler",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Crawler Starting...")
		return nil
	},
}
