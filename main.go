package main

import (
	"fmt"
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
	screen       [][]string
	lastTick     time.Time
	enemies      []enemy
	n            int
}

type sprite struct {
	width  int
	height int
	char   rune
}

type player struct {
	sprite
	y         float32
	ySpeed    float32
	jumpSpeed float32
	gravity   float32
}

type enemy struct {
	sprite
	xSpeed float32
	x      float32
}

func (s sprite) render(screen [][]string, x int, y int) {
	for i := 0; i < s.width; i++ {
		for j := 0; j < s.height; j++ {
			screen[y-j][x+i] = string(s.char)
		}
	}
}

func (p player) isGrounded(floorHeight float32) bool {
	return p.y == floorHeight
}

func (p *player) jump() {
	p.ySpeed = p.jumpSpeed
}

func (m *model) spawnEnemy() {
	e := enemy{
		sprite: sprite{
			height: 2,
			width:  4,
			char:   'X',
		},
		x: float32(m.screenWidth - 4 - 1),
	}
	m.enemies = append(m.enemies, e)
	m.n = m.n + 1
	w, _ := tea.LogToFile("/tmp/tea.log", "")
	fmt.Fprintf(w, "%v", m.n)
}

func (m *model) resetScreen() {
	for i := 0; i < m.screenHeight; i++ {
		m.screen[i] = make([]string, m.screenWidth)
		for j := 0; j < m.screenWidth; j++ {
			m.screen[i][j] = " "
		}
	}
}

func initialModel() model {
	screenWidth, screenHeight, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to get terminal size: %s", err)
	}
	player := player{
		sprite: sprite{
			width:  4,
			height: 5,
			char:   '*',
		},
		y:         float32(screenHeight - 1),
		ySpeed:    0,
		jumpSpeed: 20,
		gravity:   30,
	}
	m := model{
		player:       player,
		screenHeight: screenHeight,
		screenWidth:  screenWidth,
		screen:       make([][]string, screenWidth),
		lastTick:     time.Now(),
		enemies:      make([]enemy, 0),
	}

	m.resetScreen()

	return m
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
			if m.player.isGrounded(float32(m.screenHeight - 1)) {
				m.player.jump()
			}
		}
	case tickMsg:
		m.mainLoop()
	case spawnMsg:
		m.spawnEnemy()
	}

	return m, nil
}

func (m *model) mainLoop() {
	deltaT := float32(time.Now().Sub(m.lastTick).Seconds())
	m.lastTick = time.Now()
	m.player.y += m.player.ySpeed * deltaT * -1
	m.player.y = min(m.player.y, float32(m.screenHeight-1))
	if !m.player.isGrounded(float32(m.screenHeight - 1)) {
		m.player.ySpeed += m.player.gravity * deltaT * -1
	}
}

func (m model) View() string {
	m.resetScreen()
	m.player.sprite.render(m.screen, 0, int(m.player.y))
	for i := 0; i < len(m.enemies); i++ {
		m.enemies[i].sprite.render(m.screen, int(m.enemies[i].x), m.screenHeight-1)
	}

	res := ""
	for i := 0; i < m.screenHeight; i++ {
		res += strings.Join(m.screen[i], "") + "\n"
	}
	return res
}

type spawnMsg struct{}

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
	go func() {
		for {
			p.Send(spawnMsg{})
			time.Sleep(time.Duration(4) * time.Second)
		}
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
