package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/JakobLybarger/formula/models"
	"github.com/JakobLybarger/formula/service"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#C81D25"))

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#C81D25"))
	return model{spinner: s, loading: true}
}

type TickMsg time.Time

type model struct {
	lastUpdate time.Time

	event models.Event

	table table.Model

	spinner spinner.Model

	loading bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tickEvery(1))
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
		if m.loading {
			m.loading = false
		}

		m.lastUpdate = time.Time(msg)
		m.event = service.GetLiveData()
		m.table = loadTable(m.event)

		return m, tickEvery(5)

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func loadTable(event models.Event) table.Model {
	columns := []table.Column{
		{Title: "Position", Width: 10},
		{Title: "Driver", Width: 10},
		{Title: "Team", Width: 20},
		{Title: "Gap to Leader", Width: 20},
	}

	rows := make([]table.Row, 20)
	for i, position := range event.Positions {
		driver := getDriver(event.Drivers, position.DriverNumber)
		intervals := getInterval(event.Intervals, position.DriverNumber)
		rows[i] = table.Row{
			strconv.Itoa(position.Position),
			driver.LastName,
			driver.TeamName,
			displayAsString(intervals.GapToLeader),
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows), table.WithHeight(21),
	)

	return t
}

func displayAsString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v

	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)

	default:
		return "null"
	}
}

func getDriver(drivers []models.Driver, driverNumber int) models.Driver {
	for _, driver := range drivers {
		if driver.Number == driverNumber {
			return driver
		}
	}

	return drivers[0]
}

func getInterval(intervals []models.Interval, driverNumber int) models.Interval {
	for _, interval := range intervals {
		if interval.DriverNumber == driverNumber {
			return interval
		}
	}

	return intervals[0]
}

func (m model) View() string {
	if m.loading {
		return fmt.Sprintf("%s Loading...", m.spinner.View())
	}

	s := fmt.Sprintf("\n%s - %s\n\n", m.event.Session.Name, m.event.Meeting.OfficialName)
	s += fmt.Sprintf("%s\n", baseStyle.Render(m.table.View()))

	return s
}
