package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	util "github.com/fubarhouse/ansible-role-tester/util"
	"github.com/spf13/cobra"
)

// runCmd represents the dockerRun command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts a container",
	Long: `Start a container from a specified image.

Volume mount locations image and id are all configurable.
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := util.AnsibleConfig{
			HostPath:         source,
			RemotePath:       destination,
			RequirementsFile: "",
			PlaybookFile:     "",
			Verbose:          verbose,
		}

		dist, e := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		if e != nil {
			log.Fatalln("Incompatible distribution was inputted.")
		}

		dist.CID = containerID

		if !config.IsAnsibleRole() {
			log.Fatalf("Path %v is not recognized as an Ansible role.", config.HostPath)
		}
		if !dist.DockerCheck() {
			dist.DockerRun(&config)
		} else {
			log.Warnf("Container %v is already running", dist.CID)
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	pwd, _ := os.Getwd()
	runCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	runCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	runCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")

	runCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	runCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	runCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")

	runCmd.MarkFlagRequired("name")
}
