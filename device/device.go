package device

import (
	"errors"
	"regexp"
	"time"

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
	conn     ssh.SSHConn
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

	err := d.conn.Disconnect()
	return err
}

func (d *IOSDevice) SendCommand(command string) (string, error) {

}

func (d *IOSDevice) SendConfig(commands []string) error {

}

func (d *IOSDevice) SaveConfig() error {

}

func (d *IOSDevice) sessionPrep() error {

	regex := "\r?(.*)[#>]"
	pattern := "#|>"

	r, _ := regexp.Compile(regex)

	out, err := d.readSSH(pattern)
	if err != nil {
		return err
	}

	if !r.MatchString(out) {
		return errors.New("failed to find prompt, pattern: " + pattern + " , output: " + out)
	}

	stringmatch := r.FindStringSubmatch(out)

	d.prompt = stringmatch[1]
	d.mode = stringmatch[0][len(stringmatch[0])-1:]
}

func (d *IOSDevice) enableMode() {

}

func (d *IOSDevice) setPaging() {

}

func (d *IOSDevice) readSSH(pattern string) (string, error) {

	outChan := make(chan string)
	errChan := make(chan error)

	go func(pattern string) {
		var out string

		r, err := regexp.Compile(pattern)

		if err != nil {
			errChan <- err
			return
		}

		result, err := d.conn.Read()

		for (err == nil) && (!r.MatchString(result)) {
			out, _ := d.conn.Read()
			result += out

		}

		outChan <- out
	}(pattern)

	select {
	case <-outChan:
		return <-outChan, nil

	case <-errChan:
		return "", <-errChan

	case <-time.After(6 * time.Second):
		err := errors.New("timeout while reading, read pattern not found pattern: " + pattern)
		return "", err
	}
}
