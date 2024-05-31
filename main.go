package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Project struct {
	Name  string
	Tasks []Task
}

type Task struct {
	Name string
	Done bool
}

type keymap struct {
	up, down, left, right, quit, enter, new key.Binding
}

type model struct {
	width        int
	height       int
	projects     []Project
	focus        int
	projectIndex int
	taskIndex    int
	keymap       keymap
	help         help.Model
}

func initModel() model {
	m := model{
		help: help.New(),
		keymap: keymap{
			up:    key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/↑", "Up")),
			down:  key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/↓", "Down")),
			left:  key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("h/←", "Left")),
			right: key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l/→", "Right")),
			quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "Quit")),
			enter: key.NewBinding(key.WithKeys("enter", "space"), key.WithHelp("enter", "Select")),
			new:   key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "New")),
		},
		projects: []Project{
			{
				Name: "Project 1", Tasks: []Task{
					{Name: "Task 1", Done: false},
					{Name: "Task 2", Done: false},
				},
			},
			{
				Name: "Project 2", Tasks: []Task{
					{Name: "Task 4", Done: false},
					{Name: "Task 5", Done: false},
				},
			},
		},
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.up):
			if m.focus == 0 {
				if m.projectIndex > 0 {
					m.projectIndex--
				}
			} else {
				if m.taskIndex > 0 {
					m.taskIndex--
				}
			}
		case key.Matches(msg, m.keymap.down):
			if m.focus == 0 {
				if m.projectIndex < len(m.projects)-1 {
					m.projectIndex++
				}
			} else {
				if m.taskIndex < len(m.projects[m.projectIndex].Tasks)-1 {
					m.taskIndex++
				}
			}
		case key.Matches(msg, m.keymap.left):
			m.focus = 0
			m.taskIndex = 0
		case key.Matches(msg, m.keymap.right):
			m.focus = 1
		case key.Matches(msg, m.keymap.enter):
			if m.focus == 0 {
				m.focus = 1
			} else {
				m.projects[m.projectIndex].Tasks[m.taskIndex].Done = !m.projects[m.projectIndex].Tasks[m.taskIndex].Done
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.up,
		m.keymap.down,
		m.keymap.left,
		m.keymap.right,
		m.keymap.quit,
		m.keymap.enter,
		m.keymap.new,
	})

	projectView := ""
	for i, p := range m.projects {
		if i == m.projectIndex {
			projectView += lipgloss.NewStyle().Background(lipgloss.Color("205")).Render(p.Name) + "\n"
		} else {
			projectView += p.Name + "\n"
		}
	}

	taskView := ""
	taskView += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("146")).PaddingBottom(1).Render(m.projects[m.projectIndex].Name) + "\n"
	for i, t := range m.projects[m.projectIndex].Tasks {
		checker := "[ ]"
		if t.Done {
			checker = "[x]"
			// strikethrough
			t.Name = lipgloss.NewStyle().Strikethrough(true).Render(t.Name)
		}

		cursor := " "
		if i == m.taskIndex && m.focus == 1 {
			cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(">")
		}

		taskView += lipgloss.NewStyle().Render(lipgloss.JoinHorizontal(lipgloss.Left, cursor, checker, " ", t.Name)) + "\n"
	}

	projectViewStyle := lipgloss.NewStyle().
		Width(m.width/4).
		Height(m.height-lipgloss.Height(help)-1).
		Padding(1, 2).
		Render(projectView)
	taskViewStyle := lipgloss.NewStyle().Width(3*(m.width/4)).Height(m.height-lipgloss.Height(help)-1).Padding(1, 2).Border(lipgloss.NormalBorder(), false, false, false, true).Render(taskView)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		projectViewStyle,
		taskViewStyle,
	) + "\n\n" + help
}

func main() {
	if _, err := tea.NewProgram(initModel(), tea.WithAltScreen()).Run(); err != nil {
		panic(err)
	}
}
