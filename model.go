package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style_model = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).Padding(1).Margin(1).BorderForeground(lipgloss.Color("62"))

// main model

type Model struct {
	loaded bool
	list   list.Model
	err    error
}

func New() *Model {
	return &Model{}
}

func fd_curr_week() time.Time {
	date := time.Now()
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	return date
}

func (m *Model) initLists(width int, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height-12)
	m.list = defaultList
	//initializing our list
	entries, err := getWeekData(db, fd_curr_week(), &counter)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 7; i++ {
		m.list.InsertItem(i, entries[i])
	}
}

func (m *Model) generateLists() {
	entries, err := getWeekData(db, fd_curr_week(), &counter)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 7; i++ {
		m.list.SetItem(i, entries[i])
	}

}
func (m *Model) getWeekNum(counter int) int {
	_, weeknum := time.Now().ISOWeek()
	if counter+weeknum > 52 {
		return counter + weeknum - 52
	} else if counter+weeknum < 0 {
		return counter + weeknum + 52
	}
	return weeknum + counter
}
func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Set(index int, e Entry) tea.Cmd {
	UpdateValue(db, e)
	return m.list.SetItem(index, e)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.loaded = true
			m.initLists(msg.Width, msg.Height)
		}
	case EntryForm:
		newEntry := msg.editEntry(msg.date, msg.content.Value())
		return m, m.Set(msg.index, newEntry)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			item := m.list.SelectedItem().(Entry)
			index := m.list.Index()
			f := NewForm(item.content)
			f.index = index
			f.date = item.date
			return f.Update(nil)
		case tea.KeyLeft:
			counter -= 1
			m.generateLists()
		case tea.KeyRight:
			counter += 1
			m.generateLists()
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	week_num := m.getWeekNum(counter)
	m.list.Title = fmt.Sprintf("Week %d", week_num)
	if m.loaded {
		weeklyView := m.list.View()
		return style_model.Render(weeklyView)
	} else {
		return "loading.."
	}
}
