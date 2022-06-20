/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"

	"github.com/savingoyal/systems-assignment/pkg/client"
)

var (
	host string
	port int
	key  []string
	json bool
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a value from the api server",
	Long: `Use this command to fetch a value for a key from the api server.

	Usage:
		client get --host <host> --port <port> --key <key> [--json]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if host == "" || port == 0 || key == nil {
			fmt.Println("host, port and key are required")
			os.Exit(1)
		}
		run()
	},
}

func run() {
	for _, k := range key {
		result, err := client.Get(host, port, k, json)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(result)
	}
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().StringVar(&host, "host", "localhost", "The host of the api server")
	getCmd.PersistentFlags().IntVar(&port, "port", 8080, "The port of the api server")
	getCmd.PersistentFlags().StringArrayVarP(&key, "key", "k", nil, "The key(uuid) to get")
	getCmd.PersistentFlags().BoolVar(&json, "json", false, "Return the value as json")
}
