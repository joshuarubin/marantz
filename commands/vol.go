package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joshuarubin/marantz/client"
	"github.com/spf13/cobra"
)

var volCmd *cobra.Command

func init() {
	volCmd = &cobra.Command{
		Use:   "vol",
		Short: "Set or get receiver power status",
		Long:  `Set or get receiver power status`,
		Run:   volMain,
	}

	marantzCmd.AddCommand(volCmd)
}

func volMain(*cobra.Command, []string) {
	initServerConfig()

	if volCmd.Flags().NArg() > 0 {
		val, err := strconv.Atoi(volCmd.Flags().Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
		client.New(&srv).Vol(val)
	} else {
		client.New(&srv).Vol()
	}
}
