package main

import (
	"flag"
	"os/exec"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

type Server struct{}

var (
	output = flag.String("output", "eDP-1", "The output display.")
)

func (s Server) Screenshot(includeCursor bool, flash bool, filename string) (success bool, filenameUsed string, err *dbus.Error) {
	RunCommand(includeCursor, filename)
	return true, filename, nil
}

func RunCommand(include_cursor bool, filename string) {
	var cursor string
	if include_cursor {
		cursor = "-c"
	} else {
		cursor = ""
	}
	cmd := exec.Command("/usr/bin/grim", "-o", *output, cursor, filename)

	cmd.Run()
}

func requestName(conn *dbus.Conn, name string) {
	reply, err := conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		panic("Name " + name + " already taken")
	}
}

func main() {
	flag.Parse()
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	log.Warn(`Starting dbus screen sharing daemon for output: `, *output)

	requestName(conn, "org.gnome.SessionManager")
	requestName(conn, "org.freedesktop.PowerManagement.Inhibit")
	requestName(conn, "org.freedesktop.ScreenSaver")
	requestName(conn, "org.gnome.Shell")
	requestName(conn, "org.gnome.Shell.Screenshot")

	s := Server{}

	err = conn.Export(s, "/org/gnome/Shell/Screenshot", "org.gnome.Shell.Screenshot")
	if err != nil {
		panic(err)
	}

	select {}
}
