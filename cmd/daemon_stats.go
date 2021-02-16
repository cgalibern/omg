/*
Copyright © 2021 OPENSVC SAS <contact@opensvc.com>

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

// Package cmd defines the opensvc command line actions and options.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/client"
)

// daemonStatsCmd represents the daemonStats command
var daemonStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Print the daemon threads and objects statistics.",
	Long: `Print the daemon threads and objects statistics.
	
The daemon gathers statistics from the selected nodes.
Statistics include CPU, memory and number of processes or threads.`,
	Run: func(cmd *cobra.Command, args []string) {
		daemonStats()
	},
}

func init() {
	daemonCmd.AddCommand(daemonStatsCmd)
}

func daemonStats() {
	api := client.New(client.Config{
		URL: "raw://opt/opensvc/var/lsnr/lsnr.sock",
	})
	//requester := client.New(client.Config{
	//	URL: "https://127.0.0.1:1215"
	//	InsecureSkipVerify: true, // get from config
	//})
	c := client.NewDaemonStatsCmdConfig()
	data, err := api.DaemonStats(*c)
	if err != nil {
		return
	}
	fmt.Println(data)
}