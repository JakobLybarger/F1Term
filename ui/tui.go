package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/JakobLybarger/formula/models"
	"github.com/JakobLybarger/formula/service"
	"github.com/JakobLybarger/formula/utils"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
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

		return m, tickEvery(10)

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
		{Title: "Pos", Width: 3},
		{Title: "Driver", Width: 10},
		{Title: "Team", Width: 20},
	}

	isRace := strings.ToLower(event.Meeting.OfficialName) == "race"

	if isRace {
		columns = append(columns, table.Column{
			Title: "Gap to Leader",
			Width: 20,
		})
	}

	rows := make([]table.Row, 20)
	for i, position := range event.Positions {
		driver := utils.GetDriver(event.Drivers, position.DriverNumber)

		if isRace {
			intervals := utils.GetInterval(event.Intervals, position.DriverNumber)
			rows[i] = table.Row{
				strconv.Itoa(position.Position),
				driver.NameAcronym,
				driver.TeamName,
				utils.DisplayAsString(intervals.GapToLeader),
			}
		} else {
			rows[i] = table.Row{
				strconv.Itoa(position.Position),
				driver.NameAcronym,
				driver.TeamName,
			}
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(40),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = lipgloss.NewStyle()

	s.Cell = s.Cell.
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)

	t.SetStyles(s)

	return t
}

func (m model) View() string {
	if m.loading {
		return fmt.Sprintf("\n%s Loading...", m.spinner.View())
	}

	s := fmt.Sprintf("\n%s - %s\n", m.event.Meeting.OfficialName, m.event.Session.Name)
	// s += fmt.Sprintf("%s\n", baseStyle.Render(m.table.View()))
	s += fmt.Sprintf("%s\n", baseStyle.Render(m.table.View()))

	return s
}
