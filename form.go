package main

import (
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style_form = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(1).Margin(1, 1).BorderForeground(lipgloss.Color("62"))

type EntryForm struct {
	content textarea.Model
	date    time.Time
	index   int
}

func NewForm(content string) *EntryForm {
	form := EntryForm{
		content: textarea.New(),
	}
	if content == "" {
		form.content.Placeholder = "Enter the content here"
	} else {
		form.content.SetValue(content)
	}
	form.content.Focus()
	return &form
}
func (f EntryForm) editEntry(date time.Time, content string) Entry {
	return Entry{date: date, content: content}
}
func (f EntryForm) Init() tea.Cmd {
	return textarea.Blink
}

func (f EntryForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if f.content.Focused() {
				f.content.Blur()
			}

		case tea.KeyCtrlC:
			return f, tea.Quit
		case tea.KeyBackspace:
			if !f.content.Focused() {

				return journal.Update(f)
			}
		default:
			if !f.content.Focused() {
				cmd = f.content.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	f.content, cmd = f.content.Update(msg)
	cmds = append(cmds, cmd)
	return f, tea.Batch(cmds...)

}

func (f EntryForm) View() string {
	return style_form.Render(lipgloss.NewStyle().Background(lipgloss.Color("24")).Padding(0, 1).Render(f.date.Format("Mon Jan _2")) + "\n\n" + f.content.View())
}
