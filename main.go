package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"layeh.com/gumble/gumble"
	_ "layeh.com/gumble/opus"
)

func main() {
	// Command line flags
	server := flag.String("server", "", "the server to connect to")
	username := flag.String("username", "", "the username of the client")
	password := flag.String("password", "", "the password of the server")
	insecure := flag.Bool("insecure", false, "skip server certificate verification")
	certificate := flag.String("certificate", "", "PEM encoded certificate and private key")
	muted := flag.Bool("muted", false, "listen only")

	flag.Parse()

	m := Mumbli{
		Config:  gumble.NewConfig(),
		Address: *server,
	}

	m.Config.Username = *username
	m.Config.Password = *password

	if *insecure {
		m.TLSConfig.InsecureSkipVerify = true
	}

	if *certificate != "" {
		cert, err := tls.LoadX509KeyPair(*certificate, *certificate)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		m.TLSConfig.Certificates = append(m.TLSConfig.Certificates, cert)
	}

	m.start(*muted)

	// Disconnect when ctrl-c is pressed
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Printf("\nDisconnecting from server...\n")
	m.Client.Disconnect()
	os.Exit(0)
}
