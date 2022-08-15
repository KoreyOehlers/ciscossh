package ssh

type Connection interface {
	Connect() error
	Disconnect() error
	Read() (string, error)
	Write(cmd string) int
}

func NewConnection(ip string, username string, password string) Connection {
	return NewSSHConn(ip, username, password)
}
