package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type model struct {
	player       player
	screenHeight int
	screenWidth  int
}

type player struct {
	spriteWidth  int
	spriteHeight int
	spriteChar   rune
	y            int
}

func (p player) view() string {
	row := strings.Repeat(
		string(p.spriteChar),
		p.spriteWidth,
	)
	return strings.Repeat(row+"\n", p.spriteHeight)
}

func initialModel() model {
	screenWidth, screenHeight, err := term.GetSize(int(os.Stdin.Fd()))
	fmt.Println(screenHeight)
	if err != nil {
		log.Fatalf("Failed to get terminal size: %s", err)
	}
	player := player{
		spriteWidth:  4,
		spriteHeight: 5,
		spriteChar:   '*',
		y:            screenHeight,
	}
	return model{
		player:       player,
		screenHeight: screenHeight,
		screenWidth:  screenWidth,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	y_padding := strings.Repeat("\n", m.player.y-m.player.spriteHeight-1)
	return y_padding + m.player.view()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
