package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"layeh.com/gumble/gumble"
	_ "layeh.com/gumble/opus"
)

func main() {
	// Command line flags
	server := flag.String("server", "localhost:64738", "the server to connect to")
	username := flag.String("username", "", "the username of the client")
	password := flag.String("password", "", "the password of the server")
	insecure := flag.Bool("insecure", false, "skip server certificate verification")
	certificate := flag.String("certificate", "", "PEM encoded certificate and private key")

	flag.Parse()

	r := Mumbli{
		Config:  gumble.NewConfig(),
		Address: *server,
	}

	r.Config.Username = *username
	r.Config.Password = *password

	if *insecure {
		r.TLSConfig.InsecureSkipVerify = true
	}

	if *certificate != "" {
		cert, err := tls.LoadX509KeyPair(*certificate, *certificate)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		r.TLSConfig.Certificates = append(r.TLSConfig.Certificates, cert)
	}

	r.start()

	for {
		// If ctrl-c is pressed: r.Client.Disconnect()
	}
}
