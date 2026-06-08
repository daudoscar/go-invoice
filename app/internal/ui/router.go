package ui

import (
	"fmt"
)

type AppController interface {
	RunLocalProcess()
	RunSyncHyperlink()
}

type Router struct {
	Menu *Menu
	App  AppController
}

func NewRouter(m *Menu, a AppController) *Router {
	return &Router{
		Menu: m,
		App:  a,
	}
}

// Start adalah satu-satunya fungsi yang dipanggil oleh main
func (r *Router) Start() {
	r.Menu.ShowHeader()

	for {
		choice := r.Menu.GetChoice()

		switch choice {
		case "1":
			r.Menu.StartProcess("SCAN & PROCESS PDF LOKAL")
			r.App.RunLocalProcess()
			r.Menu.EndProcess()

		case "2":
			r.Menu.StartProcess("SYNC HYPERLINK")
			r.App.RunSyncHyperlink()
			r.Menu.EndProcess()

		case "0":
			r.Menu.ShowExitMessage()
			return

		default:
			fmt.Println("\n[!] Pilihan tidak tersedia, silakan coba lagi.")
		}
	}
}
