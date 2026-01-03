/*
Copyright © 2025 Eden Phillips
*/
package session

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Quit        key.Binding
	stopSession key.Binding
	Continue    key.Binding
	Skip        key.Binding
	Save        key.Binding
	Back        key.Binding
	Exit        key.Binding
}

var DefaultKeyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	stopSession: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "stop session"),
	),
	Continue: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "continue"),
	),
	Skip: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "skip to notes"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s", "ctrl+enter"),
		key.WithHelp("ctrl+s/ctrl+enter", "save"),
	),
	Back: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "go back"),
	),
	Exit: key.NewBinding(
		key.WithKeys("enter", "q"),
		key.WithHelp("enter/q", "exit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.stopSession},
		{k.Continue, k.Skip},
		{k.Save, k.Back},
		{k.Exit},
	}
}

func (k KeyMap) sessionHelp() []key.Binding {
	return []key.Binding{k.stopSession, k.Quit}
}

func (k KeyMap) MilestoneHelp() []key.Binding {
	return []key.Binding{k.Continue, k.Quit}
}

func (k KeyMap) FocusQualityHelp() []key.Binding {
	return []key.Binding{k.Continue, k.Skip, k.Quit}
}

func (k KeyMap) InterruptionsHelp() []key.Binding {
	return []key.Binding{k.Continue, k.Skip, k.Back, k.Quit}
}

func (k KeyMap) ReflectionHelp() []key.Binding {
	return []key.Binding{k.Save, k.Back, k.Quit}
}

func (k KeyMap) DoneHelp() []key.Binding {
	return []key.Binding{k.Exit}
}

type stateKeyMap struct {
	bindings []key.Binding
}

func (k stateKeyMap) ShortHelp() []key.Binding {
	return k.bindings
}

func (k stateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.bindings}
}

func (k KeyMap) sessionKeyMap() help.KeyMap {
	return stateKeyMap{bindings: k.sessionHelp()}
}

func (k KeyMap) MilestoneKeyMap() help.KeyMap {
	return stateKeyMap{bindings: k.MilestoneHelp()}
}

func (k KeyMap) FocusQualityKeyMap() help.KeyMap {
	return stateKeyMap{bindings: k.FocusQualityHelp()}
}

func (k KeyMap) InterruptionsKeyMap() help.KeyMap {
	return stateKeyMap{bindings: k.InterruptionsHelp()}
}

func (k KeyMap) ReflectionKeyMap() help.KeyMap {
	return stateKeyMap{bindings: k.ReflectionHelp()}
}

func (k KeyMap) DoneKeyMap() help.KeyMap {
	return stateKeyMap{bindings: k.DoneHelp()}
}

func (k KeyMap) SavingKeyMap() help.KeyMap {
	return stateKeyMap{bindings: []key.Binding{k.Quit}}
}
