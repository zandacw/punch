package main

import (
	"os"
)

func main() {

	switch os.Args[1] {
	case "client":
		Client()
	case "server":
		Server()
	}

}


