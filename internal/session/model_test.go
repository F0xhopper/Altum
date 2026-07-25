/*
Copyright © 2025 Eden Phillips
*/
package session

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func keyRune(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func quits(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

// Typing "q" in a review step must go into the focused input, not quit the
// program and lose the session.
func TestTypingQInInputStatesDoesNotQuit(t *testing.T) {
	states := []sessionState{stateMilestone, stateFocusQuality, stateInterruptions, stateReflection}

	for _, state := range states {
		m := InitialModel(t.TempDir() + "/altum.db")
		m.state = state

		updated, cmd := m.Update(keyRune('q'))
		if quits(cmd) {
			t.Fatalf("state %v: typing 'q' quit the program", state)
		}

		um := updated.(model)
		var got string
		switch state {
		case stateMilestone:
			um.milestoneInput.Focus()
			updated, _ = um.Update(keyRune('q'))
			got = updated.(model).milestoneInput.Value()
		case stateFocusQuality:
			um.focusQualityInput.Focus()
			updated, _ = um.Update(keyRune('q'))
			got = updated.(model).focusQualityInput.Value()
		case stateInterruptions:
			um.interruptionsInput.Focus()
			updated, _ = um.Update(keyRune('q'))
			got = updated.(model).interruptionsInput.Value()
		case stateReflection:
			um.reflectionInput.Focus()
			updated, _ = um.Update(keyRune('q'))
			got = updated.(model).reflectionInput.Value()
		}
		if got == "" {
			t.Fatalf("state %v: 'q' was not typed into the focused input", state)
		}
	}
}

// ctrl+c must still quit from every review step.
func TestCtrlCStillQuitsInInputStates(t *testing.T) {
	states := []sessionState{stateMilestone, stateFocusQuality, stateInterruptions, stateReflection}

	for _, state := range states {
		m := InitialModel(t.TempDir() + "/altum.db")
		m.state = state

		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		if !quits(cmd) {
			t.Fatalf("state %v: ctrl+c did not quit", state)
		}
	}
}
