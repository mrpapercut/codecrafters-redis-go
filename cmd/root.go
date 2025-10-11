package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/codecrafters-io/redis-starter-go/server"
	"github.com/spf13/cobra"
)

var runLocal bool

var rootCmd = &cobra.Command{
	Use:   "app <commands>",
	Short: "Codecrafters Redis challenge",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.Flags().BoolVar(&runLocal, "local", false, "Run on localhost")
}

func startServer() {
	s := server.GetInstance()
	defer s.Close()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		for range signalChannel {
			s.Close()

			os.Exit(0)
		}
	}()

	address := "0.0.0.0:6379"
	if runLocal {
		address = "127.0.0.1:6379"
	}

	go s.StartListening(address)

	if runLocal {
		fmt.Println("started listening")
	}

	select {}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("error running rootcmd: %v", err)
	}
}
