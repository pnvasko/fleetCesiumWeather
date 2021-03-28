package appcontroler

type DbConfig struct {
	Name           string
	Adapter        string
	Host           string
	Port           int
	Database       string
	User           string
	Password       string
	MaxConnections int
}
