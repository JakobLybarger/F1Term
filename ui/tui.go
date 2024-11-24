package ui

import (
	"fmt"
	"time"

	"github.com/JakobLybarger/formula/models"
	"github.com/JakobLybarger/formula/service"
	tea "github.com/charmbracelet/bubbletea"
)

func InitialModel() model {
	return model{}
}

type TickMsg time.Time

type model struct {
	lastUpdate time.Time

	event models.Event
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

func getDriver(drivers []models.Driver, position models.Position) models.Driver {
	for _, driver := range drivers {
		if driver.Number == position.DriverNumber {
			return driver
		}
	}

	return drivers[0]
}

func (m model) View() string {
	s := fmt.Sprintf("\n%s - %s\n\n", m.event.Session.Name, m.event.Meeting.OfficialName)

	// s += fmt.Sprintf("%s %s\n\n", m.event.Session.Name, m.lastUpdate)

	for _, position := range m.event.Positions {
		driver := getDriver(m.event.Drivers, position)

		s += fmt.Sprintf("%d %s %s\n", position.Position, driver.LastName, driver.TeamName)
	}

	return s
}
