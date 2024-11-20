package ui

import (
	"fmt"
	"time"

	"github.com/JakobLybarger/formula/service"
	tea "github.com/charmbracelet/bubbletea"
)

func InitialModel() model {
	return model{}
}

type TickMsg time.Time

type model struct {
	lastUpdate time.Time

	event service.FormulaOneEvent
}

func (m model) Init() tea.Cmd {
	return tickEvery(1)
}

func tickEvery(second time.Duration) tea.Cmd {
	return tea.Every(time.Second*second,
		func(t time.Time) tea.Msg {
			return TickMsg(t)
		})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case TickMsg:
		m.lastUpdate = time.Time(msg)
		m.event = service.GetLiveData()

		return m, tickEvery(5)
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("%s\n\n", m.event.Meeting.OfficialName)

	s += fmt.Sprintf("%s %s\n\n", m.event.Session.Name, m.lastUpdate)

	for _, driver := range m.event.Drivers {
		s += fmt.Sprintf("%s %s %s\n", driver.FirstName, driver.LastName, driver.TeamName)
	}

	return s
}
