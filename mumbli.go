package main

import (
	"crypto/tls"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleopenal"
)

type Mumbli struct {
	Config *gumble.Config
	Client *gumble.Client

	Address   string
	TLSConfig tls.Config

	Stream *gumbleopenal.Stream
}
