package terminal

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type TextMsg string

type Model struct {
	viewport viewport.Model
	content  string
	width    int
	height   int
}

func New() Model {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.FillHeight = true

	return Model{
		viewport: vp,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.viewport.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height)

		m.updateContent()

	case TextMsg:
		text := string(msg)

		switch text {
		case "\b":
			if len(m.content) > 0 {
				m.content = m.content[:len(m.content)-1]
			}
		default:
			m.content += text
		}

		m.updateContent()
		m.viewport.GotoBottom()
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

func (m *Model) updateContent() {
	if m.width == 0 {
		return
	}

	lines := strings.Split(m.content, "\n")
	centered := make([]string, 0, len(lines))

	style := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)

	for _, line := range lines {
		centered = append(centered, style.Render(line))
	}

	m.viewport.SetContent(strings.Join(centered, "\n"))
}

func (m Model) View() tea.View {
	v := tea.NewView(m.viewport.View())
	v.AltScreen = true
	return v
}
