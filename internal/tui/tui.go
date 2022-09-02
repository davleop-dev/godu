package tui

import (
	"fmt"
	. "internal/du"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("FFFDF5")).
			Background(lipgloss.Color("25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type listMsg bool

func updateList() tea.Cmd {
	return func() tea.Msg {
		return listMsg(true)
	}
}

type Order int64

const (
	Undefined Order = iota
	Name
	Size
	ModTime
)

type Model struct {
	// This section is for maintaining the `du` content
	CurrentFolder Folder
	Root          Folder

	// other options
	ListOrder      Order
	Descending     bool
	ShowHidden     bool
	DirectoryFirst bool

	// the rest is for actually maintaining the TUI display
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	Version      string
}

func (o Order) String() string {
	switch o {
	case Name:
		return "name"
	case Size:
		return "size"
	case ModTime:
		return "modify"
	}
	return "unknown"
}

type item struct {
	title       string
	description string
	bck         *Model
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) Bck() *Model         { return i.bck }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

func (m Model) updateCurrentFiles(folder Folder) []list.Item {
	if m.DirectoryFirst {
		/*switch m.ListOrder {
		case Name:
		case Size:
		case ModTime:
		}*/
		fileCount := len(folder.Files)
		folderCount := len(folder.Folders)
		totalCount := fileCount + folderCount
		items := make([]list.Item, totalCount)

		for i := 0; i < folderCount; i++ {
			title := m.formatFolderItemTitle(m.CurrentFolder.Folders[i])
			items[i] = item{title: title, bck: &m}
		}
		j := 0
		for i := folderCount; i < totalCount; i++ {
			title := m.formatFileItemTitle(m.CurrentFolder.Files[j])
			items[i] = item{title: title, bck: &m}
			j++
		}
		return items
	} else {
		items := make([]list.Item, 1)
		items[0] = item{title: "test", bck: &m}
		return items
	}
	// TODO(david): add filter checks here...
	/*for _, file := range m.Files {
		if !m.ShowHidden && strings.HasPrefix(file.Name, ".") {
			continue
		}
		if file.HighDir == m.CurrentDirectory && file.Name != m.CurrentDirectory {
			if file.IsDir {
				m.currentDirectories = append(m.currentDirectories, file)
			} else {
				m.currentFiles = append(m.currentFiles, file)
			}
		}
	}

	// Calculate usage of currently displayed items

	// TODO(david): there's probably a prettier way to do this, but it'll  work
	if !m.DirectoryFirst {
		m.currentFiles = append(m.currentFiles, m.currentDirectories...)

		switch m.ListOrder {
		case Name:
			sort.Sort(NameSorter(m.currentFiles))
		case Size:
			sort.SliceStable(m.currentFiles, func(i, j int) bool {
				return m.currentFiles[i].Size > m.currentFiles[j].Size
			})
		case ModTime:
			sort.Sort(TimeSorter(m.currentFiles))
		}

		fileCount := len(m.currentFiles)
		items := make([]list.Item, fileCount)
		for i := 0; i < fileCount; i++ {
			title := m.formatItemTitle(m.currentFiles[i])
			items[i] = item{title: title}
		}
		return items
	} else {
		switch m.ListOrder {
		case Name:
			sort.Sort(NameSorter(m.currentDirectories))
			sort.Sort(NameSorter(m.currentFiles))
		case Size:
			sort.SliceStable(m.currentDirectories, func(i, j int) bool {
				return m.currentDirectories[i].Size > m.currentDirectories[j].Size
			})
			sort.SliceStable(m.currentFiles, func(i, j int) bool {
				return m.currentFiles[i].Size > m.currentFiles[j].Size
			})
		case ModTime:
			sort.Sort(TimeSorter(m.currentDirectories))
			sort.Sort(TimeSorter(m.currentFiles))
		}

		// Create list view
		directoryCount := len(m.currentDirectories)
		fileCount := len(m.currentFiles)
		totalCount := directoryCount + fileCount
		items := make([]list.Item, totalCount)

		if m.DirectoryFirst {
			for i := 0; i < directoryCount; i++ {
				title := m.formatItemTitle(m.currentDirectories[i])
				items[i] = item{title: title}
			}
			j := 0
			for i := directoryCount; i < totalCount; i++ {
				title := m.formatItemTitle(m.currentFiles[j])
				items[i] = item{title: title}
				j++
			}
			return items
		} else {
			for i := 0; i < fileCount; i++ {
				title := m.formatItemTitle(m.currentFiles[i])
				items[i] = item{title: title}
			}
			j := 0
			for i := fileCount; i < totalCount; i++ {
				title := m.formatItemTitle(m.currentDirectories[j])
				items[i] = item{title: title}
				j++
			}
			return items
		}
	}*/

}

func (m Model) formatFileItemTitle(file File) string {
	// this should formatted eventually as so:
	// F SSS.S UUU [BBBBBBBBB] filename -->
	// TODO(david): calculate sizes later
	prog := progress.New(progress.WithScaledGradient("#00FF00", "#FF0000"))
	prog.Width = 11
	n := 0.0
	humanSize := file.HumanSize
	n = float64(file.Size) / float64(m.Root.Size)
	graph := prog.ViewAs(n)

	// setting `F` here
	mode := " "
	if !file.Mode.IsRegular() {
		mode = "@"
	}
	return fmt.Sprintf("%-2s %8s %-9s   %s", mode, humanSize, graph, file.Name)
}

func (m Model) formatFolderItemTitle(file Folder) string {
	// this should formatted eventually as so:
	// F SSS.S UUU [BBBBBBBBB] filename -->
	// TODO(david): calculate sizes later
	prog := progress.New(progress.WithScaledGradient("#00FF00", "#FF0000"))
	prog.Width = 11
	n := 0.0
	humanSize := file.HumanSize
	n = float64(file.Size) / float64(m.Root.Size)
	humanSize = file.HumanSize
	graph := prog.ViewAs(n)

	// setting `F` here
	mode := " "

	return fmt.Sprintf("%-2s %8s %-9s   %s/", mode, humanSize, graph, file.Name)
}
func NewModel(m Model) Model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	items := m.updateCurrentFiles(m.CurrentFolder)

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	delegate.ShowDescription = false
	currentFiles := list.New(items, delegate, 0, 0)
	title := fmt.Sprintf("godu-%s | Total: %s | %s", m.Version, PrettyPrintSize(m.Root.Size), m.CurrentFolder.Path)
	currentFiles.Title = title
	currentFiles.Styles.Title = titleStyle
	currentFiles.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	m.list = currentFiles
	m.keys = listKeys
	m.delegateKeys = delegateKeys

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.insertItem):
			m.delegateKeys.remove.SetEnabled(true)
			return m, nil
		}

	case listMsg:
		//
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return appStyle.Render(m.list.View())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
