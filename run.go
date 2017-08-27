package main

import (
	"os"
	"strings"

	"github.com/kasheemlew/xperiMoby/cgroups"
	"github.com/kasheemlew/xperiMoby/cgroups/subsystems"
	"github.com/kasheemlew/xperiMoby/container"
	"github.com/sirupsen/logrus"
)

// Run envokes the command
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	cgroupManager := cgroups.NewCgroupManager("xperiMoby-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	// add container processes to cgroup
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(comArray, writePipe)
	if tty {
		parent.Wait()
		mntURL := "/root/mnt/"
		rootURL := "/root/"
		container.DeleteWorkSpace(rootURL, mntURL, volume)
	}
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
