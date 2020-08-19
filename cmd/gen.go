/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cage1016/rbacgen/gen"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate rbac sql",
	Long:  `generate rbac sql`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("out")
		if file == "" {
			logrus.Error("You must provide the output path")
			return
		}

		osFile, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		keys := viper.GetViper().AllKeys()
		if len(keys) > 0 {
			gen.Gen(osFile, viper.GetViper().AllSettings())
		} else {
			gen.Gen(osFile, nil)
		}

		fmt.Println("Successfuly wrote file to " + file)
		err = osFile.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	genCmd.Flags().StringP("out", "o", "", "Set output file path")
}
