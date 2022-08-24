package ciscossh

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func NewDevice(
	name string,
	ip string,
	username string,
	password string,
	args ...string, //args used for enable password
) IOSDevice {

	enable := ""

	if len(args) > 0 {
		enable = args[0]
	}

	iosdevice := IOSDevice{
		Name:     name,
		IP:       ip,
		Username: username,
		Password: password,
		Enable:   enable,
	}

	return iosdevice
}

func GetCredentials() (string, string, error) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println("")

	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func GetEnable() (string, error) {

	fmt.Print("Enter Enable Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println("")

	password := string(bytePassword)

	return strings.TrimSpace(password), nil
}
