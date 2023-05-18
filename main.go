package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/PondWader/GoPractice/config"
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/server"
	"github.com/PondWader/GoPractice/utils"
)

var version string = "0.0.1"

func main() {
	rand.Seed(time.Now().UnixNano())

	printName()

	utils.Info("Loading server config...")
	cfg := config.LoadConfig()
	utils.Info("Config loaded.")

	server := server.New(cfg, version)

	utils.Info("Server starting...")
	listener, err := net.Listen("tcp", ":"+fmt.Sprint(server.Config.Port))
	if err != nil {
		utils.Error("Error listening on port "+fmt.Sprint(server.Config.Port)+":", err)
		return
	}
	defer listener.Close()

	utils.Info("Server listening on port", server.Config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.Error("Error accepting connection:", err)
			os.Exit(1)
		}

		utils.Info("Incoming connection from", conn.RemoteAddr())
		go handleClient(conn, server)
	}
}

func handleClient(conn net.Conn, s *server.Server) {
	client := protocol.NewClient(conn, s)
	if client.Ended == false {
		server.NewPlayer(client, s)
	}
}

func printName() {
	fmt.Println(utils.Cyan(`  ____       ____                 _   _
 / ___| ___ |  _ \ _ __ __ _  ___| |_(_) ___ ___
| |  _ / _ \| |_) | '__/ _` + "`" + ` |/ __| __| |/ __/ _ \
| |_| | (_) |  __/| | | (_| | (__| |_| | (_|  __/
 \____|\___/|_|   |_|  \__,_|\___|\__|_|\___\___|
	`))
}
