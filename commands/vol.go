package commands

import (
	"github.com/joshuarubin/marantz/client"
	"github.com/spf13/cobra"
)

var vol struct {
	cmd   *cobra.Command
	value int
}

func init() {
	vol.cmd = &cobra.Command{
		Use:   "vol",
		Short: "Set or get receiver power status",
		Long:  `Set or get receiver power status`,
		Run:   volMain,
	}

	vol.cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Increase volume",
		Long:  `Increase volume`,
		Run: func(*cobra.Command, []string) {
			initServerConfig()
			client.New(&srv).Vol(client.VolMsg{
				Mode: client.VolModeUp,
			})
		},
	})

	vol.cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Decrease volume",
		Long:  `Decrease volume`,
		Run: func(*cobra.Command, []string) {
			initServerConfig()
			client.New(&srv).Vol(client.VolMsg{
				Mode: client.VolModeDown,
			})
		},
	})

	vol.cmd.AddCommand(&cobra.Command{
		Use:   "upfast",
		Short: "Fast increase volume",
		Long:  `Fast increase volume`,
		Run: func(*cobra.Command, []string) {
			initServerConfig()
			client.New(&srv).Vol(client.VolMsg{
				Mode: client.VolModeUpFast,
			})
		},
	})

	vol.cmd.AddCommand(&cobra.Command{
		Use:   "downfast",
		Short: "Fast decrease volume",
		Long:  `Fast decrease volume`,
		Run: func(*cobra.Command, []string) {
			initServerConfig()
			client.New(&srv).Vol(client.VolMsg{
				Mode: client.VolModeDownFast,
			})
		},
	})

	vol.cmd.Flags().IntVarP(&vol.value, "value", "v", -1, "set specific value")

	marantzCmd.AddCommand(vol.cmd)
}

func volMain(*cobra.Command, []string) {
	initServerConfig()

	if vol.cmd.Flags().Lookup("value").Changed {
		client.New(&srv).Vol(client.VolMsg{
			Mode: client.VolModeSet,
			Val:  vol.value,
		})
		return
	}

	client.New(&srv).Vol(client.VolMsg{
		Mode: client.VolModeGet,
	})
}
