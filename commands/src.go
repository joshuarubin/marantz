package commands

import (
	"github.com/joshuarubin/marantz/client"
	"github.com/joshuarubin/marantz/msg"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "src",
		Short: "Select or get receiver source",
		Long:  `Select or get receiver source`,
		Run: func(*cobra.Command, []string) {
			initServerConfig()
			client.SendCmd(srv.Config.String(), &msg.Cmd{
				Cmd: msg.Cmd_CMD_SRC.Enum(),
			})
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "tv",
			Short: "Select TV source",
			Long:  "Select TV source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_TV.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "dvd",
			Short: "Select DVD source",
			Long:  "Select DVD source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_DVD.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "vcr",
			Short: "Select VCR1 source",
			Long:  "Select VCR1 source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_VCR1.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "dss",
			Short: "Select DSS/VCR2 source",
			Long:  "Select DSS/VCR2 source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_DSS_VCR2.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "aux1",
			Short: "Select AUX1 source",
			Long:  "Select AUX1 source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_AUX1.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "aux2",
			Short: "Select AUX2 source",
			Long:  "Select AUX2 source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_AUX2.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "cd",
			Short: "Select CD/CD-R source",
			Long:  "Select CD/CD-R source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_CD_CDR.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "tape",
			Short: "Select Tape source",
			Long:  "Select Tape source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_TAPE.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "tuner",
			Short: "Select Tuner source",
			Long:  "Select Tuner source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_TUNER1.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "fm",
			Short: "Select FM1 source",
			Long:  "Select FM1 source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_FM1.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "am",
			Short: "Select AM1 source",
			Long:  "Select AM1 source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_AM1.Enum(),
				})
			},
		})

	cmd.AddCommand(
		&cobra.Command{
			Use:   "xm",
			Short: "Select XM1 source",
			Long:  "Select XM1 source",
			Run: func(*cobra.Command, []string) {
				initServerConfig()
				client.SendCmd(srv.Config.String(), &msg.Cmd{
					Cmd: msg.Cmd_CMD_SRC.Enum(),
					Src: msg.Cmd_SRC_XM1.Enum(),
				})
			},
		})

	marantzCmd.AddCommand(cmd)
}
