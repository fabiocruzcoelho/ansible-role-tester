package util

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func (dist *Distribution) IdempotenceTest(config *AnsibleConfig) {

	// Test role idempotence.
	log.Infoln("Testing role idempotence...")
	out, _ := DockerExec([]string{
		"exec",
		"--tty",
		dist.CID,
		"ansible-playbook",
		fmt.Sprintf("%v/tests/%v", config.RemotePath, config.PlaybookFile),
	}, true)

	idempotence := IdempotenceResult(out)
	if idempotence {
		log.Infoln("Idempotence test: PASS")
	} else {
		log.Errorln("Idempotence test: FAIL")
	}
}

func IdempotenceResult(output string) bool {

	lines := strings.Split(output, "\n")

	changed := ""
	failed := ""

	for _, line := range lines {
		if strings.Contains(line, "=") {
			f := strings.Split(line, "=")
			if strings.Contains(line, "changed=") {
				changed = strings.Split(f[2], " ")[0]
			}
			if strings.Contains(line, "failed=") {
				failed = strings.Split(f[4], " ")[0]
			}
		}
	}

	if failed != "0" {
		return false
	}

	if changed != "0" {
		return false
	}

	return true
}