package tui

import (
	"fmt"
	. "internal/du"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.SetSpacing(0)
	d.SetHeight(0)

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string
		var bck *Model

		if i, ok := m.SelectedItem().(item); ok {
			bck = i.Bck()
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				if title == "                          .." {
					stackLength := len(bck.Stack)
					if stackLength > 0 {
						n := stackLength - 1
						bck.CurrentFolder = bck.Stack[n]
						bck.Stack = bck.Stack[:n]
						m.SetItems(bck.updateCurrentFiles(bck.CurrentFolder))
					}
					newTitle := fmt.Sprintf("godu-%s | Total: %s | %s", bck.Version, PrettyPrintSize(bck.CurrentFolder.Size), bck.CurrentFolder.Path)
					m.Title = newTitle
				} else {
					for _, folder := range bck.CurrentFolder.Folders {
						if strings.Contains(title, folder.Name) {
							bck.Stack = append(bck.Stack, bck.CurrentFolder)
							bck.CurrentFolder = folder
							newTitle := fmt.Sprintf("godu-%s | Total: %s | %s", bck.Version, PrettyPrintSize(bck.CurrentFolder.Size), bck.CurrentFolder.Path)
							m.Title = newTitle

							m.SetItems(bck.updateCurrentFiles(bck.CurrentFolder))
							return updateList()
						}
					}
				}

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle("Deleted " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}
