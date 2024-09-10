/*
Copyright © 2020 Romber Li <romber2001@gmail.com>

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

	"github.com/romberli/go-util/constant"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/romberli/go-mod/config"
	"github.com/romberli/go-mod/module/mod"
	"github.com/romberli/go-mod/pkg/message"

	msgMod "github.com/romberli/go-mod/pkg/message/mod"
)

// parentCmd represents the parent command
var parentCmd = &cobra.Command{
	Use:   "parent",
	Short: "parent command",
	Long:  `print parent of the package.`,
	Run: func(cmd *cobra.Command, args []string) {
		// init config
		err := initConfig()
		if err != nil {
			fmt.Println(fmt.Sprintf(constant.LogWithStackString, message.NewMessage(message.ErrInitConfig, err)))
			os.Exit(constant.DefaultAbnormalExitCode)
		}

		modDir := viper.GetString(config.ModDirKey)
		modName := viper.GetString(config.ModNameKey)
		modVersion := viper.GetString(config.ModVersionKey)

		c := mod.NewController(modDir)
		err = c.PrintParentChain(modName, modVersion)
		if err != nil {
			fmt.Println(fmt.Sprintf(constant.LogWithStackString, message.NewMessage(
				msgMod.ErrModParentPrintParentChain, err, modDir, modName, modVersion)))
		}

		os.Exit(constant.DefaultNormalExitCode)
	},
}

func init() {
	rootCmd.AddCommand(parentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
