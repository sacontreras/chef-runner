// Package ssh implements the driver.Driver interface. The implementation is
// just a simple wrapper around openssh and rsync.
package ssh

import (
	"fmt"
	"strings"

	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/mlafeldt/chef-runner/rsync"
)

type Driver struct {
	host        string
	sshClient   *openssh.Client
	rsyncClient *rsync.Client
}

func NewDriver(host string) (*Driver, error) {
	sshClient, err := openssh.NewClient(host)
	if err != nil {
		return nil, err
	}

	sshCmd := sshClient.Command("")
	remoteShell := strings.Join(sshCmd[:len(sshCmd)-1], " ")
	rsyncClient := &rsync.Client{
		Archive:     true,
		Delete:      true,
		Verbose:     true,
		RemoteHost:  sshClient.Host,
		RemoteShell: remoteShell,
	}

	return &Driver{host, sshClient, rsyncClient}, nil
}

func (drv Driver) String() string {
	return fmt.Sprintf("SSH driver (host: %s)", drv.sshClient.Host)
}

func (drv Driver) RunCommand(command string) error {
	return drv.sshClient.RunCommand(command)
}

func (drv Driver) Upload(dst string, src ...string) error {
	return drv.rsyncClient.Copy(dst, src...)
}
