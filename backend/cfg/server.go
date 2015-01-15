package cfg

import (
	"fmt"
	"os"
)

type server struct {
	Host string
	Port string

	Frontend struct {
		Path string
	}

	API struct {
		Prefix  string
		PrivKey string
		PubKey  string
	}
}

var Server server

func init() {
	Server.Host = os.Getenv("SERVER_HOST")
	Server.Port = os.Getenv("SERVER_PORT")
	Server.Frontend.Path = os.Getenv("SERVER_FRONTEND_PATH")
	Server.API.Prefix = os.Getenv("SERVER_API_PREFIX")
	Server.API.PrivKey = os.Getenv("SERVER_PRIVKEY")
	Server.API.PubKey = os.Getenv("SERVER_PUBKEY")
}

func (s *server) BasePath() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
