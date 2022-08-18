# ciscossh

Go methods for connecting to Cisco devices over SSH

### Usage: 

```
import (
	"fmt"

	"github.com/koreyoehlers/ciscossh"
)

func main() {

	username, password, _ := ciscossh.GetCredentials()
	testswitch := ciscossh.NewDevice("test", "10.1.1.1", username, password)

	err := testswitch.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer testswitch.Disconnect()
	
	results, err := testswitch.SendCommand("show run")
	fmt.Println(results)
}
```

### If you need to use an enable password:

```
	username, password, _ := ciscossh.GetCredentials()
	enable, _ := ciscossh.GetEnable()

	testswitch := ciscossh.NewDevice("test", "10.1.1.1", username, password, enable)

```