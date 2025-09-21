package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh/"

type statusMsg int

// errMsg represents an error message.
type errMsg struct{ err error }

// Error returns the string value of err.
func (err errMsg) Error() string { return err.err.Error() }

// model stores the application state.
type model struct {
	status int
	err    error
}

// Init returns a tea.Cmd that is run at application start up.
func (model model) Init() tea.Cmd {
	// The Bubble Tea runtime will call the function when the time is right.
	return checkServer
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

// Update updates the application state.
func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		// The server returned a status message. Save it to our model. Also tell
		// the Bubble Tea runtime we want to exit because we have nothing else to
		// do. We'll still be able to render a final view with out status message.
		model.status = int(msg)

		return model, tea.Quit

	case errMsg:
		// There was an error. Note it in the model. Tell the runtime we're done
		// and want to quit.
		model.err = msg

		return model, tea.Quit

	case tea.KeyMsg:
		// Ctrl+c exits. Even with short running programs it's good to have a quit
		// key, just in case your logic is off. Users will be very annoyed if they
		// can't exit.
		if msg.Type == tea.KeyCtrlC {
			return model, tea.Quit
		}
	}

	// If we happen to get any other messages, don't do anything.
	return model, nil
}

// View displays the application state.
func (model model) View() string {
	// If there's an error, print it out and don't do anything else.
	if model.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", model.err)
	}

	// Tell the user we're doing something.
	s := fmt.Sprintf("Checking %s ...", url)

	// When the server responds with a status, add it to the current line.
	if model.status > 0 {
		s += fmt.Sprintf("%d %s!", model.status, http.StatusText(model.status))
	}

	// Send off whatever we came up with above for rendering.
	return "\n" + s + "\n\n"
}

// checkServer sends a GET request to a url and returns the http.StatusCode as a statusMsg.
func checkServer() tea.Msg {
	httpClient := &http.Client{Timeout: time.Duration(10) * time.Second}
	response, err := httpClient.Get(url)
	if err != nil {
		return errMsg{err}
	}

	return statusMsg(response.StatusCode)
}
