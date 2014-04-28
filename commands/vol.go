package commands

import (
	"code.google.com/p/gogoprotobuf/proto"
	"github.com/joshuarubin/marantz/client"
	"github.com/joshuarubin/marantz/msg"
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
		Run: func(*cobra.Command, []string) {
			volCmd(nil)
		},
	}

	vol.cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Increase volume",
		Long:  `Increase volume`,
		Run: func(*cobra.Command, []string) {
			volCmd(msg.Cmd_VOL_UP.Enum())
		},
	})

	vol.cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Decrease volume",
		Long:  `Decrease volume`,
		Run: func(*cobra.Command, []string) {
			volCmd(msg.Cmd_VOL_DOWN.Enum())
		},
	})

	vol.cmd.AddCommand(&cobra.Command{
		Use:   "upfast",
		Short: "Fast increase volume",
		Long:  `Fast increase volume`,
		Run: func(*cobra.Command, []string) {
			volCmd(msg.Cmd_VOL_UP_FAST.Enum())
		},
	})

	vol.cmd.AddCommand(&cobra.Command{
		Use:   "downfast",
		Short: "Fast decrease volume",
		Long:  `Fast decrease volume`,
		Run: func(*cobra.Command, []string) {
			volCmd(msg.Cmd_VOL_DOWN_FAST.Enum())
		},
	})

	vol.cmd.Flags().IntVarP(&vol.value, "value", "v", -1, "set specific value")

	marantzCmd.AddCommand(vol.cmd)
}

func volCmd(value *msg.Cmd_VolValue) {
	initServerConfig()

	cmd := &msg.Cmd{
		Cmd: msg.Cmd_CMD_VOL.Enum(),
	}

	switch value {
	case nil:
		if vol.cmd.Flags().Lookup("value").Changed {
			cmd.IntValue = proto.Int32(int32(vol.value))
		}
	default:
		cmd.Vol = value
	}

	client.SendCmd(srv.Config.String(), cmd)
}
