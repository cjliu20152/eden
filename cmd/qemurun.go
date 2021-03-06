package cmd

import (
	"fmt"
	"github.com/lf-edge/eden/pkg/utils"
	"github.com/spf13/cobra"
	"log"
	"runtime"
	"strings"
)

var (
	qemuARCH         string
	qemuOS           string
	qemuAccel        bool
	qemuSMBIOSSerial string
	qemuConfigFile   string
)

var qemuRunCmd = &cobra.Command{
	Use:   "qemurun",
	Short: "run qemu-system with eve",
	Long:  `Run qemu-system with eve.`,
	Run: func(cmd *cobra.Command, args []string) {
		qemuCommand := ""
		qemuOptions := "-display none -serial mon:stdio -nodefaults -no-user-config "
		if qemuSMBIOSSerial != "" {
			qemuOptions += fmt.Sprintf("-smbios type=1,serial=%s ", qemuSMBIOSSerial)
		}
		if qemuOS == "" {
			qemuOS = runtime.GOOS
		} else {
			qemuOS = strings.ToLower(qemuOS)
		}
		if qemuOS != "linux" && qemuOS != "darwin" {
			log.Fatalf("OS not supported: %s", qemuOS)
		}
		if qemuARCH == "" {
			qemuARCH = runtime.GOARCH
		} else {
			qemuARCH = strings.ToLower(qemuARCH)
		}
		switch qemuARCH {
		case "amd64":
			qemuCommand = "qemu-system-x86_64"
			if qemuAccel {
				if qemuOS == "darwin" {
					qemuOptions += "-M accel=hvf --cpu host "
				} else {
					qemuOptions += "-enable-kvm --cpu host "
				}
			} else {
				qemuOptions += "--cpu SandyBridge "
			}
		case "arm64":
			qemuCommand = "qemu-system-aarch64"
			qemuOptions += "-machine virt,gic_version=3 -machine virtualization=true -cpu cortex-a57 -machine type=virt "
		default:
			log.Fatalf("Arch not supported: %s", runtime.GOARCH)
		}
		if qemuConfigFile != "" {
			qemuOptions += fmt.Sprintf("-readconfig %s ", qemuConfigFile)
		}
		err := utils.RunCommandForeground(qemuCommand, strings.Fields(qemuOptions)...)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func qemuRunInit() {
	qemuRunCmd.Flags().StringVarP(&qemuARCH, "arch", "", "", "arch of system")
	qemuRunCmd.Flags().StringVarP(&qemuOS, "os", "", "", "os to run on")
	qemuRunCmd.Flags().BoolVarP(&qemuAccel, "accel", "", true, "use acceleration")
	qemuRunCmd.Flags().StringVarP(&qemuSMBIOSSerial, "serial", "", "", "SMBIOS serial")
	qemuRunCmd.Flags().StringVarP(&qemuConfigFile, "config", "", "", "config file to use")
}
