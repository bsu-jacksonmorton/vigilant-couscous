package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	// "github.com/hajimehoshi/ebiten/inpututil"
)

const (
	WINDOW_LIMIT_LEFT  float64 = -32
	WINDOW_LIMIT_RIGHT float64 = 206
	IDLE                       = 0
	WALKING                    = 1
	FACE_LEFT                  = -1
	FACE_RIGHT                 = 1
)

var (
	backgroundImage *ebiten.Image
	zombieImages    [][]*ebiten.Image
	knightImages    [][]*ebiten.Image
	frameCounter    float64
	direction       int
)

type Zombie struct {
	state     int
	frame     float32
	direction float64
	posX      float64
	posY      float64
}

func (p *Zombie) SetState(state int) {
	p.state = state
}
func (p *Zombie) IncrementFrameCounter() {
	p.frame = p.frame + 0.2
}
func (p *Zombie) SetDirection(direction float64) {
	p.direction = direction
}
func (p *Zombie) SetPosX(posX float64) {
	p.posX = posX
}
func (p *Zombie) SetPosY(posY float64) {
	p.posY = posY
}
func (p *Zombie) SetPos(x float64, y float64) {
	p.posX = x
	p.posY = y
}
func (p *Zombie) GetPos() (float64, float64) {
	return p.posX, p.posY
}
func (p *Zombie) GetState() int {
	return p.state
}
func (p *Zombie) GetFrame() float32 {
	return p.frame
}
func (p *Zombie) GetDirection() float64 {
	return p.direction
}

type Player struct {
	state     int
	frame     float32
	direction float64
	posX      float64
	posY      float64
}

func (p *Player) SetState(state int) {
	p.state = state
}
func (p *Player) IncrementFrameCounter() {
	if p.frame > 9 {
		p.frame = 0
	} else {
		p.frame = p.frame + 0.2
	}
}
func (p *Player) SetDirection(direction float64) {
	p.direction = direction
}
func (p *Player) SetPosX(posX float64) {
	p.posX = posX
}
func (p *Player) SetPosY(posY float64) {
	p.posY = posY
}
func (p *Player) SetPos(x float64, y float64) {
	p.posX = x
	p.posY = y
}
func (p *Player) GetPos() (float64, float64) {
	return p.posX, p.posY
}
func (p *Player) GetState() int {
	return p.state
}
func (p *Player) GetFrame() float32 {
	return p.frame
}
func (p *Player) GetDirection() float64 {
	return p.direction
}

type Game struct {
	p *Player
}

func (g *Game) Update(*ebiten.Image) error {
	playerX, _ := g.p.GetPos()
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if playerX > WINDOW_LIMIT_LEFT {
			g.p.SetPosX(playerX - 2)
		}
		g.p.SetState(WALKING)
		g.p.SetDirection(FACE_LEFT)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if playerX < 206 {
			g.p.SetPosX(playerX + 2)
		}
		g.p.SetState(WALKING)
		g.p.SetDirection(FACE_RIGHT)
	} else {
		g.p.SetState(IDLE)
	}
	frameCounter += 0.2
	playerState := g.p.GetState()
	if playerState == WALKING && frameCounter > 9 {
		frameCounter = 0
	} else if playerState == IDLE && frameCounter > 9 {
		frameCounter = 0
	}
	return nil
}
func drawPlayer(screen *ebiten.Image, player *Player, frame int) {
	op := &ebiten.DrawImageOptions{}
	state := player.GetState()
	w, _ := knightImages[state][frame].Bounds().Dx(), knightImages[state][frame].Bounds().Dy()
	x, y := player.GetPos()
	op.GeoM.Scale(player.GetDirection(), 1)
	op.GeoM.Translate(x+float64(w), y)
	screen.DrawImage(knightImages[state][frame], op)
}
func drawZombie(screen *ebiten.Image, frame int) {
	op := &ebiten.DrawImageOptions{}
	w, _ := zombieImages[0][frame].Bounds().Dx(), zombieImages[0][frame].Bounds().Dy()
	op.GeoM.Scale(1, 1)
	op.GeoM.Translate(5+float64(w), 150)
	screen.DrawImage(zombieImages[0][frame], op)
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, &ebiten.DrawImageOptions{})
	drawPlayer(screen, g.p, int(math.Floor(frameCounter)))
	drawZombie(screen, int(math.Floor(frameCounter)))
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("my-first-game")
	player := &Player{0, 0, 1, 1, 150}
	game := &Game{player}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func init() {
	// Load background
	var err error
	backgroundImage, _, err = ebitenutil.NewImageFromFile("png/test.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	// LOAD ZOMBIE FRAMES
	// zombie idle
	var zombieImage *ebiten.Image
	zombieImages = [][]*ebiten.Image{}
	zombieImages = append(zombieImages, []*ebiten.Image{}, []*ebiten.Image{})
	for i := 1; i < 16; i++ {
		zombieImage, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("png/male/idle%d.png", i), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
			return
		}
		zombieImages[0] = append(zombieImages[0], zombieImage)
	}
	// zombie walk
	for i := 1; i < 11; i++ {
		zombieImage, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("png/male/walk%d.png", i), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
			return
		}
		zombieImages[1] = append(zombieImages[1], zombieImage)
	}
	// LOAD KNIGHT FRAMES
	var knightImage *ebiten.Image
	knightImages = [][]*ebiten.Image{}
	// Idle
	knightIdle := []*ebiten.Image{}
	for i := 1; i < 11; i++ {
		knightImage, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("png/knight/knight-real-idle%d.png", i), ebiten.FilterDefault)
		knightIdle = append(knightIdle, knightImage)
	}
	knightImages = append(knightImages, knightIdle)
	knightRun := []*ebiten.Image{}
	for i := 1; i < 11; i++ {
		knightImage, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("png/knight/knight-run%d.png", i), ebiten.FilterDefault)
		knightRun = append(knightRun, knightImage)
	}
	knightImages = append(knightImages, knightRun)
	frameCounter = 1
}
