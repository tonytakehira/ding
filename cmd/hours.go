/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"strconv"
	"fmt"

	"github.com/spf13/cobra"

)

// minutesCmd represents the minutes command
var hoursCmd = &cobra.Command{
	Use:   "hours",
	Short: "sets a timer using hours",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		numberFlt, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Println("Could not read value for timer.")
		}
		fmt.Println("153762", int(numberFlt * 60 * 60))
		go startCountDown(int(numberFlt * 60 * 60))
		go waitForExit()
		openGame()

	},
}


func init() {
	rootCmd.AddCommand(hoursCmd)
}
