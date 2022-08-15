package device

import (
	"github.com/koreyoehlers/ciscossh/ssh"
)

type IOSDevice struct {
	Name     string
	IP       string
	Username string
	Password string
	Enable   string
	prompt   string
	mode     string
	conn     ssh.Connection
}

func (d *IOSDevice) Connect() error {

	d.conn = ssh.NewConnection(d.IP, d.Username, d.Password)
	err := d.conn.Connect()

	if err != nil {
		return err
	}

	err = d.sessionPrep()

	if err != nil {
		return err
	}

	return nil
}

func (d *IOSDevice) Disconnect() error {

	return nil
}

func (d *IOSDevice) SendCommand(command string) string {

}

func (d *IOSDevice) SendConfig(commands []string) {

}

func (d *IOSDevice) sessionPrep() error {

}

func (d *IOSDevice) enableMode() {

}

func (d *IOSDevice) setPaging() {

}
