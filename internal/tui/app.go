package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/aungmyatmoe/zencode/internal/config"
	"github.com/aungmyatmoe/zencode/internal/ui"
)

// Msg types
type (
	assistantChunkMsg string
	toolResultMsg     string
)

// Model is the root bubbletea model for Zencode TUI.
type Model struct {
	cfg      *config.Config
	viewport viewport.Model
	textarea textarea.Model
	entries  []string // transcript lines (we'll make richer entries later)
	width    int
	height   int
	ready    bool

	// Fake "agent is working" state for the scaffold
	working bool
	mode    string // Normal | Plan | YOLO
}

func New(cfg *config.Config) Model {
	ta := textarea.New()
	ta.Placeholder = "Type a message...  (Enter=send  Esc=focus scrollback  Ctrl+C=cancel  ? = help)"
	ta.Focus()
	ta.CharLimit = 0
	ta.SetWidth(80)
	ta.SetHeight(3)
	ta.ShowLineNumbers = false

	vp := viewport.New(80, 20)
	vp.SetContent("Welcome to Zencode (early TUI scaffold).\n\n" +
		"• This is a real bubbletea app.\n" +
		"• LLM + real tools land in the next passes.\n" +
		"• Try typing and pressing Enter (it will fake a response + todo).\n\n" +
		"Keyboard hints: j/k scroll (later), Tab focus, Ctrl+O = yolo toggle, Shift+Tab = plan mode.\n")

	m := Model{
		cfg:      cfg,
		textarea: ta,
		viewport: vp,
		entries:  []string{},
		mode:     "Normal",
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+o":
			// Toggle YOLO (fake for now)
			if m.mode == "YOLO" {
				m.mode = "Normal"
			} else {
				m.mode = "YOLO"
			}
			return m, nil

		case "shift+tab":
			// Cycle mode (Normal -> Plan -> YOLO -> Normal)
			switch m.mode {
			case "Normal":
				m.mode = "Plan"
			case "Plan":
				m.mode = "YOLO"
			default:
				m.mode = "Normal"
			}
			return m, nil

		case "?":
			m.appendEntry(ui.Dim.Render("Help (stub): Enter sends • Ctrl+O yolo • Shift+Tab mode • Ctrl+C quit"))
			m.viewport.GotoBottom()
			return m, nil

		case "esc":
			// Focus scrollback area (simple: just blur prompt)
			m.textarea.Blur()
			return m, nil

		case "tab":
			m.textarea.Focus()
			return m, nil

		case "enter":
			if !m.textarea.Focused() {
				// If scrollback focused, Enter could do something later
				m.textarea.Focus()
				return m, nil
			}
			val := strings.TrimSpace(m.textarea.Value())
			if val == "" {
				return m, nil
			}
			m.appendEntry(ui.UserHeader.Render("You: ") + val)
			m.textarea.Reset()
			m.viewport.GotoBottom()

			// Simulate agent turn (real agent loop will replace this)
			m.working = true
			return m, tea.Batch(
				m.simulateAssistant(val),
				textarea.Blink,
			)
		}

	case assistantChunkMsg:
		m.appendEntry(string(msg))
		m.viewport.GotoBottom()
		m.working = false
		return m, nil

	case toolResultMsg:
		m.appendEntry(ui.ToolHeader.Render("[tool] ") + string(msg))
		m.viewport.GotoBottom()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			// First size event — set sensible sizes
			m.viewport = viewport.New(msg.Width, msg.Height-6)
			m.textarea.SetWidth(msg.Width)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 6
			m.textarea.SetWidth(msg.Width)
		}
	}

	// Route remaining input to sub components
	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(taCmd, vpCmd)
}

func (m *Model) appendEntry(s string) {
	m.entries = append(m.entries, s)
	content := strings.Join(m.entries, "\n")
	m.viewport.SetContent(content)
}

func (m Model) simulateAssistant(userMsg string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(120 * time.Millisecond)

		// Very fake but demonstrates the transcript + "tool" result
		if strings.Contains(strings.ToLower(userMsg), "todo") || strings.Contains(strings.ToLower(userMsg), "plan") {
			return assistantChunkMsg(ui.Success.Render("Agent: ") + "I'll create a todo list for this task.\n\n" + ui.Dim.Render("[todo_write]") + "\n- [ ] Explore the request\n- [x] Show the TUI works\n- [ ] Wire real LLM + search_replace")
		}

		if strings.Contains(strings.ToLower(userMsg), "edit") || strings.Contains(strings.ToLower(userMsg), "fix") {
			return assistantChunkMsg(ui.Success.Render("Agent: ") + "I would call search_replace here (unique old_string required).\n\n" +
				ui.DiffHeader.Render("```diff\n") + ui.DiffAdd.Render("+ added line\n") + ui.DiffDel.Render("- removed line\n") + "```")
		}

		return assistantChunkMsg(ui.Success.Render("Agent: ") + "Thanks! (This is still a scaffold — real tool-calling loop and Grok/xAI calls are next.)\nTry messages containing 'todo', 'plan', or 'edit' for demo output.")
	}
}

func (m Model) View() string {
	if !m.ready {
		return "Loading Zencode TUI..."
	}

	status := ui.StatusBar.Render(fmt.Sprintf(" zencode • mode:%s • %s • (scaffold TUI — real agent soon) ", m.mode, m.cfg.Model))
	promptView := m.textarea.View()

	// Simple layout
	return lipgloss.JoinVertical(
		lipgloss.Left,
		status,
		m.viewport.View(),
		"",
		promptView,
		ui.Dim.Render("  Enter send • Esc blur prompt • Ctrl+O yolo • Shift+Tab cycle mode • ? help • Ctrl+C quit"),
	)
}

// Run launches the TUI program.
func Run(cfg *config.Config) error {
	p := tea.NewProgram(New(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
