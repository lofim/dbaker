package config

type Config struct {
	Host     string
	Port     uint
	Database string
	Username string
	Password string
	SSLMode  string
	Tables   []string
	DataSize uint32
	IterFrom uint32
}
