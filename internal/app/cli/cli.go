package cli

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultInterval = 5
	defaultTimeout  = 1.5
)

var (
	interval time.Duration
	timeout  time.Duration
)

func Run() (time.Duration, time.Duration, error) {
	p := tea.NewProgram(initialModel())

	m, err := p.Run()
	if err != nil {
		return 0, 0, err
	}

	model, ok := m.(model)
	if !ok {
		return 0, 0, fmt.Errorf("model is not of type model")
	}

	if model.err != nil {
		return 0, 0, model.err
	}

	return interval, timeout, nil
}

const defaultWidth = 50

const (
	intervalKey = iota
	timeoutKey
)

const (
	redColor      = lipgloss.Color("#FF0000")
	darkGrayColor = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle()
	errMsgStyle   = lipgloss.NewStyle().Foreground(redColor)
	continueStyle = lipgloss.NewStyle().Foreground(darkGrayColor)
)

type model struct {
	inputs  []textinput.Model
	focused int
	err     error
}

type validationError error

var (
	errInvalidInterval validationError = errors.New("Interval must be a number")
	errInvalidTimeout  validationError = errors.New("Timeout must be a number")
)

func intervalValidator(s string) error {
	if s == "" {
		return nil
	}
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errInvalidInterval
	}
	return nil
}

func timeoutValidator(s string) error {
	if s == "" {
		return nil
	}
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errInvalidTimeout
	}
	return nil
}

func initialModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 2)

	inputs[intervalKey] = textinput.New()
	inputs[intervalKey].Placeholder = fmt.Sprintf("%d", defaultInterval)
	inputs[intervalKey].Focus()
	inputs[intervalKey].CharLimit = 4
	inputs[intervalKey].Width = defaultWidth
	inputs[intervalKey].Prompt = ""
	inputs[intervalKey].Validate = intervalValidator

	inputs[timeoutKey] = textinput.New()
	inputs[timeoutKey].Placeholder = fmt.Sprintf("%.1f", defaultTimeout)
	inputs[timeoutKey].CharLimit = 4
	inputs[timeoutKey].Width = defaultWidth
	inputs[timeoutKey].Prompt = ""
	inputs[timeoutKey].Validate = timeoutValidator

	return model{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if ok := m.validateInputs(); !ok {
				return m, nil
			}
			if m.focused == len(m.inputs)-1 {
				m.setValues()
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyEsc:
			m.setValues()
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case error:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	errMsg := ""
	if m.err != nil {
		errMsg = errMsgStyle.
			Render(m.err.Error())
	}

	return fmt.Sprintf(
		`
%s
%s
%s

%s
%s

%s
%s
`,
		inputStyle.Width(defaultWidth).
			Bold(true).
			Render("Press Esc to use default values and Ctrl+C to exit"),
		inputStyle.Width(defaultWidth).
			Render("How often should I press a key? (in minutes)"),
		m.inputs[intervalKey].View(),
		inputStyle.Width(defaultWidth).
			Render("How much time should I run? (in hours)"),
		m.inputs[timeoutKey].View(),
		errMsg,
		continueStyle.Render("Continue ->"),
	) + "\n"
}

// setValues sets the interval and timeout values with the input values
func (m *model) setValues() {
	strInterval := m.inputs[intervalKey].Value()
	if strInterval == "" {
		strInterval = fmt.Sprintf("%d", defaultInterval)
	}

	strTimeout := m.inputs[timeoutKey].Value()
	if strTimeout == "" {
		strTimeout = fmt.Sprintf("%f", defaultTimeout)
	}

	float64Interval, _ := strconv.ParseFloat(strInterval, 64)
	interval = time.Duration(float64Interval * float64(time.Minute))

	float64Timeout, _ := strconv.ParseFloat(strTimeout, 64)
	timeout = time.Duration(float64Timeout * float64(time.Hour))
}

// validateInputs returns true if all inputs are valid
func (m *model) validateInputs() bool {
	for i := range m.inputs {
		if err := m.inputs[i].Validate(m.inputs[i].Value()); err != nil {
			m.err = fmt.Errorf("%v", err)
			return false
		}
	}

	switch (m.err).(type) {
	case validationError:
		m.err = nil
	}

	return true
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
