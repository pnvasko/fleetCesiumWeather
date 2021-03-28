package appcontroler

type HttpConfig struct {
	Domain            string
	Port              int
	WsAllowOriginHost string
	StaticPath        string
	Sslcert           string
	Sslkey            string
}
