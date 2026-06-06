/*
Copyright © 2025 Eden Phillips
*/
package menu

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Quit   key.Binding
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
}

var DefaultKeyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "navigate up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "navigate down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "select"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Select, k.Quit},
	}
}

type menuKeyMap struct {
	bindings []key.Binding
}

func (k menuKeyMap) ShortHelp() []key.Binding {
	return k.bindings
}

func (k menuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.bindings}
}

func (k KeyMap) MenuKeyMap() help.KeyMap {
	return menuKeyMap{bindings: []key.Binding{k.Up, k.Down, k.Select, k.Quit}}
}
