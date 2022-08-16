package main

import (
	"fmt"

	"github.com/koreyoehlers/ciscossh"
)

func main() {

	username, password, _ := ciscossh.GetCredentials()
	testsw := ciscossh.NewDevice("test", "10.44.1.10", username, password)

	err := testsw.Connect()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer testsw.Disconnect()

}
