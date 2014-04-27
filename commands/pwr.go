package commands

import (
	"github.com/joshuarubin/marantz/client"
	"github.com/joshuarubin/marantz/msg"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "pwr",
		Short: "Set or get receiver power status",
		Long:  `Set or get receiver power status`,
		Run: func(*cobra.Command, []string) {
			initServerConfig()
			client.SendCmd(srv.Config.String(), &msg.Cmd{
				Cmd: msg.Cmd_CMD_PWR.Enum(),
			})
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "on",
			Short: "Turn receiver on",
			Long:  "Turn receiver on",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_PWR.Enum(),
					Pwr: msg.Cmd_PWR_ON.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "off",
			Short: "Turn receiver off",
			Long:  "Turn receiver off",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_PWR.Enum(),
					Pwr: msg.Cmd_PWR_OFF.Enum(),
				})
			},
		})

	marantzCmd.AddCommand(cmd)
}
