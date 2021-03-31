package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/kennygrant/sanitize"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleopenal"
	"layeh.com/gumble/gumbleutil"
)

func esc(str string) string {
	return sanitize.HTML(str)
}

func (m *Mumbli) start(muted bool) {
	m.Config.Attach(gumbleutil.AutoBitrate)
	m.Config.Attach(m)

	// Retry until connecting to mumble server was successful
	for {
		var err error
		_, err = gumble.DialWithDialer(new(net.Dialer), m.Address, m.Config, &m.TLSConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while connecting to mumble server: %s\n", err)
			fmt.Printf("Retrying in 10 seconds...\n")
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	// Audio
	if os.Getenv("ALSOFT_LOGLEVEL") == "" {
		os.Setenv("ALSOFT_LOGLEVEL", "0")
	}

	// Create new client
	if stream, err := gumbleopenal.New(m.Client); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		m.Stream = stream
	}

	// Start audio stream
	if !muted {
		m.Stream.StartSource() // To toggle voice: m.Stream.StopSource()
	}
}

func (m *Mumbli) OnConnect(e *gumble.ConnectEvent) {
	m.Client = e.Client

	fmt.Printf("To: %s\n", e.Client.Self.Channel.Name)
	fmt.Printf("Connected to %s\n", m.Client.Conn.RemoteAddr())
	if e.WelcomeMessage != nil {
		fmt.Printf("Welcome message: %s\n", esc(*e.WelcomeMessage))
	}
}

func (m *Mumbli) OnDisconnect(e *gumble.DisconnectEvent) {
	var reason string
	switch e.Type {
	case gumble.DisconnectError:
		reason = "connection error"
	}
	if reason == "" {
		fmt.Println("Disconnected")
	} else {
		fmt.Println("Disconnected: " + reason)
	}
}

func (m *Mumbli) OnTextMessage(e *gumble.TextMessageEvent) {
	fmt.Println(e.Sender, e.Message)
}

func (m *Mumbli) OnUserChange(e *gumble.UserChangeEvent) {
	if e.Type.Has(gumble.UserChangeChannel) && e.User == m.Client.Self {
		fmt.Printf("To: %s\n", e.User.Channel.Name)
	}
}

func (m *Mumbli) OnChannelChange(e *gumble.ChannelChangeEvent) {
}

func (m *Mumbli) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	var info string
	switch e.Type {
	case gumble.PermissionDeniedOther:
		info = e.String
	case gumble.PermissionDeniedPermission:
		info = "insufficient permissions"
	case gumble.PermissionDeniedSuperUser:
		info = "cannot modify SuperUser"
	case gumble.PermissionDeniedInvalidChannelName:
		info = "invalid channel name"
	case gumble.PermissionDeniedTextTooLong:
		info = "text too long"
	case gumble.PermissionDeniedTemporaryChannel:
		info = "temporary channel"
	case gumble.PermissionDeniedMissingCertificate:
		info = "missing certificate"
	case gumble.PermissionDeniedInvalidUserName:
		info = "invalid user name"
	case gumble.PermissionDeniedChannelFull:
		info = "channel full"
	case gumble.PermissionDeniedNestingLimit:
		info = "nesting limit"
	}
	fmt.Printf("Permission denied: %s\n", info)
}

func (m *Mumbli) OnUserList(e *gumble.UserListEvent) {
}

func (m *Mumbli) OnACL(e *gumble.ACLEvent) {
}

func (m *Mumbli) OnBanList(e *gumble.BanListEvent) {
}

func (m *Mumbli) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
}

func (m *Mumbli) OnServerConfig(e *gumble.ServerConfigEvent) {
}
