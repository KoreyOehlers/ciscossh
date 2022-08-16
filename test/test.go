package main

import "github.com/koreyoehlers/ciscossh"

func main() {

	username, password, _ := ciscossh.GetCredentials()
	testsw := ciscossh.NewDevice("test", "10.")
}
