package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/abunasyr7/go-final-project/pkg/api"
)

const (
	defaultPort = 7540
	webDir = "./web"
)

func Run() error {
	port, err := getPort()

	if err != nil {
		return err
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init()

	addr := fmt.Sprintf(":%d", port)
	log.Printf("server listening in %s", addr)

	return http.ListenAndServe(addr, nil)
}

func getPort() (int, error) {
	env := os.Getenv("TODO_PORT")

	if env == "" {
		return defaultPort, nil
	}

	port, err := strconv.Atoi(env)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid TODO_PORT %q: must be a number from 1 to 65535", env)
	}
	return port, nil
}