package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JakobLybarger/formula/models"
	"github.com/JakobLybarger/formula/service"
	"github.com/JakobLybarger/formula/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

type ApiResponseMsg models.Event

type model struct {
	lastUpdate time.Time

	event models.Event

	table *table.Table

	spinner spinner.Model

	loading bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tickEvery(1))
}

func tickEvery(second time.Duration) tea.Cmd {
	return tea.Every(time.Second*second,
		func(t time.Time) tea.Msg {
			event := service.GetLiveData()
			return ApiResponseMsg(event)
		})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case ApiResponseMsg:
		if m.loading {
			m.loading = false
		}

		m.lastUpdate = time.Now()
		m.event = models.Event(msg)
		m.table = loadTable(m.event)

		return m, tickEvery(10)
	}

	return m, nil
}

func loadTable(event models.Event) *table.Table {

	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)

	headers := []string{"Pos", "Driver", "Team"}

	isRace := strings.ToLower(event.Session.Name) == "race"
	if isRace {
		headers = append(headers, "Gap to Leader")
	}

	colors := map[string]string{}

	rows := make([][]string, 0)
	for _, position := range event.Positions {
		driver := utils.GetDriver(event.Drivers, position.DriverNumber)
		colors[driver.TeamName] = driver.TeamColor
		row := []string{strconv.Itoa(position.Position), driver.NameAcronym, driver.TeamName}

		if isRace {
			if interval, ok := utils.GetInterval(event.Intervals, position.DriverNumber); ok {
				row = append(row, utils.DisplayAsString(interval.GapToLeader))
			}
		}

		rows = append(rows, row)
	}

	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {

			if row < 0 {
				return headerStyle
			}

			if row >= len(rows) {
				return baseStyle
			}

			if col == 2 {
				team := rows[row][2]
				if color, ok := colors[team]; ok {
					return baseStyle.Foreground(lipgloss.Color("#" + color))
				}
			}

			return baseStyle.Foreground(lipgloss.Color("252"))
		})

	return t
}

func (m model) View() string {
	if m.loading {
		return fmt.Sprintf("\n%s Loading...", m.spinner.View())
	}

	s := fmt.Sprintf("\n%s - %s\n", m.event.Meeting.OfficialName, m.event.Session.Name)
	s += fmt.Sprintf("%s\n", m.table)

	return s
}
