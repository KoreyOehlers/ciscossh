package ciscossh

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type IOSDevice struct {
	Name     string
	IP       string
	Username string
	Password string
	Enable   string
	prompt   string
	mode     string
	conn     SSHConn
}

func (d *IOSDevice) Connect() error {

	d.conn = NewConnection(d.IP, d.Username, d.Password)
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

	if (SSHConn{}) == d.conn {
		return "", errors.New("failed to send command. Run Connect() before" +
			" SendCommand()")
	}

	command = command + "\n"

	err := d.conn.Write(command)
	if err != nil {
		return "", err
	}

	results, err := d.readSSH(d.prompt)
	if err != nil {
		return "", err
	}

	// Remove the prompt and the command from results
	command = strings.Replace(command, "\n", "", -1)
	results = strings.Replace(results, command, "", -1)
	results = strings.Replace(results, "\n"+d.prompt, "", -1)

	return results, nil
}

func (d *IOSDevice) SendConfig(commands []string) error {

	_, err := d.SendCommand("config t")
	if err != nil {
		return errors.New("could not enter config mode")
	}

	for _, cmd := range commands {
		results, err := d.SendCommand(cmd)
		if err != nil {
			return errors.New("could not send command " + cmd + " " + results)
		}
	}

	_, err = d.SendCommand("end")
	if err != nil {
		return errors.New("could not exit config mode")
	}

	return nil
}

func (d *IOSDevice) SaveConfig() error {

	_, err := d.SendCommand("write mem")

	if err != nil {
		return errors.New("failed to save config: " + err.Error())
	}

	return nil
}

func (d *IOSDevice) sessionPrep() error {

	regex := `(\w.*)[#>]`
	pattern := "#|>"
	r, _ := regexp.Compile(regex)

	out, err := d.readSSH(pattern)
	if err != nil {
		return errors.New("failed session prep: " + err.Error())
	}

	if !r.MatchString(out) {
		return errors.New("failed to find prompt pattern: " + pattern +
			", output: " + out)
	}

	stringmatch := r.FindStringSubmatch(out)

	d.prompt = stringmatch[1]
	d.mode = stringmatch[0][len(stringmatch[0])-1:]

	if d.prompt == "" || d.mode == "" {
		return errors.New("failed to get prompt or mode")
	}

	if d.mode != "#" {

		if d.Enable == "" {
			return errors.New("failed to enter enable mode: enter enable " +
				"password after the user password when creating NewDevice")
		}

		err = d.enableMode()
		if err != nil {
			return errors.New("failed to enter enable mode: " + err.Error())
		}
	}

	d.prompt = d.prompt + d.mode

	err = d.setPaging()
	if err != nil {
		return errors.New("failed to set paging: " + err.Error())
	}

	return nil
}

func (d *IOSDevice) enableMode() error {

	if d.mode != ">" {
		return errors.New("> not found in mode string")
	}

	err := d.conn.Write("enable\n")
	if err != nil {
		return errors.New("error sending enable command")
	}

	err = d.conn.Write(d.Enable + "\n")
	if err != nil {
		return errors.New("error sending enable password")
	}

	_, err = d.readSSH(d.prompt + "#")
	if err != nil {
		return errors.New("incorrect enable password or other issue at enable")
	}

	d.mode = "#"

	return nil
}

func (d *IOSDevice) setPaging() error {

	command := "terminal length 0"

	_, err := d.SendCommand(command)

	if err != nil {
		return errors.New("could not send terminal length command " +
			err.Error())
	}

	return nil
}

func (d *IOSDevice) readSSH(pattern string) (string, error) {

	outChan := make(chan string)
	errChan := make(chan error)

	go func(pattern string) {

		r, err := regexp.Compile(pattern)
		if err != nil {
			err = errors.New(err.Error() + " " + pattern)
			errChan <- err
			return
		}

		result, err := d.conn.Read()
		if err != nil {
			errChan <- err
			return
		}

		if r.MatchString(result) {
			outChan <- result
		}

		for (err == nil) && (!r.MatchString(result)) {
			out, err := d.conn.Read()

			if err != nil {
				errChan <- err
				return
			}

			result += out
		}

		outChan <- result
	}(pattern)

	select {

	case recv := <-outChan:
		return recv, nil

	case recv := <-errChan:
		return "", recv

	case <-time.After(6 * time.Second):
		err := errors.New("timeout while reading, read pattern not found" +
			" pattern: " + pattern)
		return "", err
	}
}
