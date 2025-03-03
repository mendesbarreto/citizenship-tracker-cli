package app

import (
	"citizenship-tracker-cli/pkg/api"
	"citizenship-tracker-cli/pkg/model"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss/tree"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	purple        = lipgloss.Color("99")
	gray          = lipgloss.Color("245")
	green         = lipgloss.Color("10")
	red           = lipgloss.Color("170")

	headerStyle   = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	cellStyle     = lipgloss.NewStyle().Padding(0, 1).Width(30)
	grayRowStyle  = cellStyle.Foreground(gray)
	redRowStyle   = cellStyle.Foreground(red).Bold(true).Align(lipgloss.Center)
	greenRowStyle = cellStyle.Foreground(green).Bold(true).Align(lipgloss.Center)
	credentials   = GetUserCredentials()
)

type (
	AppState          int
	statusResponseMsg struct{ resp *model.StatusResponse }
)

const ( // iota is reset to 0
	login        AppState = iota // c0 == 0
	loading      AppState = iota // c1 == 1
	statusResult AppState = iota // c2 == 2
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

type TeaModel struct {
	FocusIndex int
	TextInputs []textinput.Model

	CursorMode      cursor.Mode
	OnSubmit        func(string, string, string) tea.Cmd
	AppState        AppState
	Spinner         spinner.Model
	IsRequesting    bool
	SaveCredentials bool
	StatusResponse  *model.StatusResponse
	err             error
}

func GetUserCredentials() model.Credential {
	uci := api.GetUserUci()
	applicationNumber := api.GetUserApplicationNumber()
	password := api.GetUserPassword()

	return model.Credential{
		Uci:               uci,
		ApplicationNumber: applicationNumber,
		Password:          password,
	}
}

func InitialTeaModel() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := TeaModel{
		FocusIndex: 0,
		TextInputs: make([]textinput.Model, 3),
		AppState:   login,
		Spinner:    s,
	}

	var t textinput.Model

	for i := range m.TextInputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case 0:
			t.Placeholder = "UCI"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 10
			if credentials.Uci != "" {
				t.SetValue(credentials.Uci)
			}
		case 1:
			t.Placeholder = "password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.CharLimit = 40
			if credentials.Password != "" {
				t.SetValue(credentials.Password)
			}
		case 2:
			t.Placeholder = "Application Number"
			t.CharLimit = 10
			if credentials.ApplicationNumber != "" {
				t.SetValue(credentials.ApplicationNumber)
			}
		}

		m.TextInputs[i] = t
	}

	return m
}

func (m TeaModel) submit(uci string, password string, applicationNumber string) tea.Cmd {
	m.IsRequesting = true

	_ = api.AddToKeychain(uci, password, applicationNumber)

	return func() tea.Msg {
		authResponse, err := api.Auth(uci, password)
		if err != nil {
			m.err = err
			return nil
		}

		authToken := fmt.Sprintf("Bearer %s", authResponse.AuthenticationResult.IdToken)
		resp, err := api.GetStatus(authToken, applicationNumber)
		if err != nil {
			m.err = err
			return nil
		}

		checkForUpdates(resp)
		api.SaveStatusResponse(resp)

		m.IsRequesting = false
		m.AppState = statusResult
		return statusResponseMsg{resp}
	}
}

func checkForUpdates(resp *model.StatusResponse) {
	loadStatusResponse, err := api.LoadStatusResponse()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(loadStatusResponse.Activities) == 0 {
		return
	}

	loadedActivities := loadStatusResponse.Activities
	udpateActivities := resp.Activities

	for _, la := range loadedActivities {
		idx := slices.IndexFunc(udpateActivities, func(ca model.Activity) bool { return ca.Activity == la.Activity })

		if idx == -1 {
			continue
		}

		fmt.Printf("old: %s | new: %s\n", la.Status, udpateActivities[idx].Status)
		updatedStatus := udpateActivities[idx].Status
		updatedActivity := udpateActivities[idx].Activity
		if la.Status == updatedStatus {
			continue
		}

		_ = api.SendNotification("Citizenship Tracker", fmt.Sprintf("%s Status Updated", stringsDict[updatedActivity]), stringsDict[updatedStatus])
		time.Sleep(1 * time.Second)
	}
}

func (m TeaModel) updateLoginState(msg tea.Msg) (TeaModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.CursorMode++
			if m.CursorMode > cursor.CursorHide {
				m.CursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.TextInputs))
			for i := range m.TextInputs {
				cmds[i] = m.TextInputs[i].Cursor.SetMode(m.CursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.FocusIndex == len(m.TextInputs) {
				if !m.IsRequesting {
					fmt.Println("Submitting...")
					m.AppState = loading

					fetch := m.submit(m.TextInputs[0].Value(), m.TextInputs[1].Value(), m.TextInputs[2].Value())
					return m, tea.Batch(fetch, spinner.Tick)
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}

			if m.FocusIndex > len(m.TextInputs) {
				m.FocusIndex = 0
			} else if m.FocusIndex < 0 {
				m.FocusIndex = len(m.TextInputs)
			}

			cmds := make([]tea.Cmd, len(m.TextInputs))
			for i := 0; i <= len(m.TextInputs)-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.TextInputs[i].Focus()
					m.TextInputs[i].PromptStyle = focusedStyle
					m.TextInputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.TextInputs[i].Blur()
				m.TextInputs[i].PromptStyle = noStyle
				m.TextInputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	return m.udateTextInputs(msg)
}

func (m TeaModel) udateTextInputs(msg tea.Msg) (TeaModel, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.TextInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.TextInputs {
		m.TextInputs[i], cmds[i] = m.TextInputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m TeaModel) updateLoadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusResponseMsg:
		return m.updateStatusResultState(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}

	default:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}
}

func (m TeaModel) updateStatusResultState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusResponseMsg:
		r := statusResponseMsg(msg)
		m.AppState = statusResult
		m.StatusResponse = r.resp
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	}

	return m, nil
}

func (m TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.AppState {
	case login:
		return m.updateLoginState(msg)
	case loading:
		return m.updateLoadingState(msg)
	case statusResult:
		return m.updateStatusResultState(msg)
	}

	return m, nil
}

func (m TeaModel) View() string {
	switch m.AppState {
	case login:
		return m.renderLogin()
	case loading:
		return m.renderLoading()
	case statusResult:
		return m.renderStatusResult()
	}

	return ""
}

var stringsDict = map[string]string{
	"language":               "Language Skills",
	"backgroundVerification": "Background Check",
	"residency":              "Physical Presence",
	"prohibitions":           "Prohibitions",
	"citizenshipTest":        "Citizenship test",
	"citizenshipOath":        "Citizenship Ceremony",
	"inProgress":             "In Progress",
	"completed":              "Completed",
	"notStarted":             "Not Started",
}

func clean(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "\r", ""), "\n", "")
}

func (m TeaModel) renderStatusResult() string {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("70")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("190"))

	t := tree.
		Root("⁜ History")

	sort.Slice(m.StatusResponse.History, func(i, j int) bool {
		return m.StatusResponse.History[i].Time > m.StatusResponse.History[j].Time
	})

	for _, h := range m.StatusResponse.History {
		t.Child(
			clean(h.Title.En),
			tree.New().Child(
				fmt.Sprintf("Description => %s", clean(h.Text.En)),
				fmt.Sprintf("Date => %s", formatUnixTime(h.Time)),
			),
		)
	}

	t.Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyle(itemStyle)

	rows := [][]string{}

	for _, activity := range m.StatusResponse.Activities {
		rows = append(rows, []string{stringsDict[activity.Activity], stringsDict[activity.Status]})
	}

	t2 := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case table.HeaderRow:
				return headerStyle
			default:
				if col == 0 {
					return grayRowStyle
				}

				if rows[row][col] == stringsDict["completed"] {
					return greenRowStyle
				}

				return redRowStyle
			}
		}).
		Headers("STEP", "STATUS").
		Rows(rows...)

	textRender := fmt.Sprintf("%s\n\n%s\n\n\n", t2.String(), t.String())

	return textRender
}

func formatUnixTime(unixTime int64) string {
	t := time.Unix(unixTime/1000, 0)
	return t.Format("2 Jan 2006 3:04pm")
}

func (m TeaModel) renderLoading() string {
	str := fmt.Sprintf("\n\n\n\n\n%s Requesting Information... press q to quit\n\n", m.Spinner.View())
	return str
}

func (m TeaModel) renderLogin() string {
	var b strings.Builder

	for i := range m.TextInputs {
		b.WriteString(m.TextInputs[i].View())
		if i < len(m.TextInputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.FocusIndex == len(m.TextInputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.CursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func (m TeaModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink)
}
