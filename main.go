package main

import (
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type model struct {
	player       player
	screenHeight int
	screenWidth  int
	lastTick     time.Time
}

type player struct {
	spriteWidth  int
	spriteHeight int
	spriteChar   rune
	y            float32
	ySpeed       float32
	jumpSpeed    float32
	gravity      float32
}

func (p player) view() string {
	row := strings.Repeat(
		string(p.spriteChar),
		p.spriteWidth,
	)
	return strings.Repeat(row+"\n", p.spriteHeight)
}

func (p player) isGrounded(floorHeight float32) bool {
	return p.y == float32(floorHeight)
}

func (p *player) jump() {
	p.ySpeed = p.jumpSpeed
}

func initialModel() model {
	screenWidth, screenHeight, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to get terminal size: %s", err)
	}
	player := player{
		spriteWidth:  4,
		spriteHeight: 5,
		spriteChar:   '*',
		y:            float32(screenHeight),
		ySpeed:       0,
		jumpSpeed:    20,
		gravity:      20,
	}
	return model{
		player:       player,
		screenHeight: screenHeight,
		screenWidth:  screenWidth,
		lastTick:     time.Now(),
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
		case " ":
			if m.player.isGrounded(float32(m.screenHeight)) {
				m.player.jump()
			}
		}
	case tickMsg:
		m.mainLoop()
	}

	return m, nil
}

func (m *model) mainLoop() {
	deltaT := float32(time.Now().Sub(m.lastTick).Seconds())
	m.lastTick = time.Now()
	m.player.y += m.player.ySpeed * deltaT * -1
	m.player.y = min(m.player.y, float32(m.screenHeight))
	if !m.player.isGrounded(float32(m.screenHeight)) {
		m.player.ySpeed += m.player.gravity * deltaT * -1
	}
}

func (m model) View() string {
	y_padding := strings.Repeat("\n", int(m.player.y)-m.player.spriteHeight-1)
	return y_padding + m.player.view()
}

type tickMsg struct{}

func tick(p *tea.Program, fps int) {
	time.Sleep(time.Duration(1/fps) * time.Second)
	p.Send(tickMsg{})
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	go func() {
		for {
			tick(p, 30)
		}
	}()
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
