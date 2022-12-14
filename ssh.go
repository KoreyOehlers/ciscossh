package ciscossh

import (
	"errors"
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

var ciphers = []string{"aes256-cbc", "aes192-cbc", "aes128-cbc", "aes128-ctr", "3des-cbc"}

type SSHConn struct {
	Host     string
	Username string
	Password string
	client   *ssh.Client
	reader   io.Reader
	writer   io.WriteCloser
}

func NewConnection(host string, username string, password string) SSHConn {
	return SSHConn{host, username, password, nil, nil, nil}
}

func (c *SSHConn) Connect() error {

	sshConfig := &ssh.ClientConfig{
		User:            c.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         6 * time.Second,
	}
	sshConfig.Ciphers = append(sshConfig.Ciphers, ciphers...)

	conn, err := ssh.Dial("tcp", c.Host+":22", sshConfig)
	if err != nil {
		return errors.New("failed to connect to device: " + err.Error())
	}

	session, err := conn.NewSession()
	if err != nil {
		return errors.New("failed to Start a new session: " + err.Error())
	}

	reader, _ := session.StdoutPipe()
	writer, _ := session.StdinPipe()

	c.client = conn
	c.reader = reader
	c.writer = writer

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	err = session.RequestPty("vt100", 0, 200, modes)

	if err != nil {
		return errors.New("failed to request Pty: " + err.Error())
	}

	err = session.Shell()
	if err != nil {
		return errors.New("failed to invoke shell: " + err.Error())
	}

	return nil
}

func (c *SSHConn) Disconnect() error {

	err := c.client.Close()

	return err
}

func (c *SSHConn) Read() (string, error) {

	buff := make([]byte, 2048)
	n, err := c.reader.Read(buff)

	return string(buff[:n]), err
}

func (c *SSHConn) Write(cmd string) error {

	commandBytes := []byte(cmd)
	_, err := c.writer.Write(commandBytes)

	return err
}
