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
	ATTACKING                  = 2
	JUMPING                    = 3
	FACE_LEFT                  = -1
	FACE_RIGHT                 = 1
	ZOMBIE_DEAD                = 2
)

var (
	backgroundImage *ebiten.Image
	zombieImages    [][]*ebiten.Image
	knightImages    [][]*ebiten.Image
	frameCounter    float64
	direction       int
)

type Actor struct {
	direction  int
	frameCount float64
	state      int
}

func (a *Actor) SetState(state int) {
	a.state = state
}
func (a *Actor) SetDirection(direction int) {
	a.direction = direction
}
func (a *Actor) SetFrameCount(fc float64) {
	a.frameCount = fc
}

type Position struct {
	posX float64
	posY float64
}

func (p *Position) Move(x, y float64) {
	p.posX += x
	p.posY += y
}

type Zombie struct {
	*Actor
	*Position
}

func (z *Zombie) Update() {
	switch z.Actor.state {
	case IDLE:
		if z.Actor.frameCount > 14 {
			z.Actor.SetFrameCount(0)
		} else {
			z.Actor.SetFrameCount(z.Actor.frameCount + 0.2)
		}
		break
	case WALKING:
		if z.Actor.frameCount > 9 {
			z.Actor.SetFrameCount(0)
		} else {
			z.Actor.SetFrameCount(z.Actor.frameCount + 0.2)
		}
		break
	case ZOMBIE_DEAD:
		if z.Actor.frameCount > 9 {
			z.Actor.SetFrameCount(0)
		} else {
			z.Actor.SetFrameCount(z.Actor.frameCount + 0.2)
		}
		break
	default:
		panic("UNKNOWN ZOMBIE STATE!!!")
	}
}
func NewZombie(posX, posY float64) *Zombie {
	return &Zombie{
		&Actor{1, 0, 0},
		&Position{
			posX: posX,
			posY: posY,
		},
	}
}

type Player struct {
	*Actor
	*Position
	vy      float64
	jumping bool
	width   int
	height  int
}

func NewPlayer(posX, posY float64) *Player {
	return &Player{
		&Actor{1, 0, 0},
		&Position{posX, posY},
		0,
		false,
		74,
		90,
	}
}

func collision(p *Player, z *Zombie) bool {
	pX, pY := p.posX, p.posY
	zX, zY := z.posX, z.posY
	w, h := float64(p.width), float64(p.height)

	return (pX < zX+w &&
		pX+w > zX &&
		pY < zY+h &&
		pY+h > zY)
}
func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if p.posX > WINDOW_LIMIT_LEFT {
			p.Position.Move(-2, 0)
		}
		p.SetState(WALKING)
		p.SetDirection(FACE_LEFT)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if p.posX < WINDOW_LIMIT_RIGHT {
			p.Move(2, 0)
		}
		p.SetState(WALKING)
		p.SetDirection(FACE_RIGHT)
	} else {
		p.SetState(IDLE)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.SetState(ATTACKING)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if p.jumping == false && p.posY == 150 {
			p.SetState(JUMPING)
			p.vy = -6
			p.jumping = true
		}
	}
	if p.jumping == true && p.vy+p.posY <= 25 {
		p.vy = 2
		p.jumping = false
	}
	if 0 <= p.posY+p.vy && p.posY+p.vy <= 150 {
		p.Move(0, p.vy)
	}
	// there are 10 frames in all of the animations for the knight model so we don't need fancy logic for this one
	if p.frameCount > 9 {
		p.SetFrameCount(0)
	} else {
		if p.state == ATTACKING {
			p.SetFrameCount(p.Actor.frameCount + 0.5)
		} else {
			p.SetFrameCount(p.Actor.frameCount + 0.15)
		}
	}
}

type Game struct {
	p *Player
	z *Zombie
}

func (g *Game) Update(*ebiten.Image) error {
	if collision(g.p, g.z) == true && g.p.state == ATTACKING {
		g.z.SetState(ZOMBIE_DEAD)
	} else {
		fmt.Println("NO COLLISION")
	}
	g.p.Update()
	g.z.Update()
	return nil
}
func drawPlayer(screen *ebiten.Image, player *Player) {
	op := &ebiten.DrawImageOptions{}
	state := player.Actor.state
	frame := int(math.Floor(player.Actor.frameCount))
	w, _ := knightImages[state][frame].Bounds().Dx(), knightImages[state][frame].Bounds().Dy()
	//fmt.Println("Width:", w, "Height:", h)
	x, y := player.posX, player.posY
	op.GeoM.Scale(float64(player.Actor.direction), 1)
	op.GeoM.Translate(x+float64(w), y)
	screen.DrawImage(knightImages[state][frame], op)
}
func drawZombie(screen *ebiten.Image, z *Zombie) {
	op := &ebiten.DrawImageOptions{}
	frame := int(math.Floor(z.frameCount))
	w, _ := zombieImages[z.state][frame].Bounds().Dx(), zombieImages[z.state][frame].Bounds().Dy()
	op.GeoM.Scale(1, 1)
	op.GeoM.Translate(50+float64(w), 150)
	screen.DrawImage(zombieImages[z.state][frame], op)
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, &ebiten.DrawImageOptions{})
	drawPlayer(screen, g.p)
	drawZombie(screen, g.z)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("my-first-game")
	player := NewPlayer(0, 150)
	zombie := NewZombie(25, 150)
	game := &Game{player, zombie}
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
	zombieImages = append(zombieImages, []*ebiten.Image{}, []*ebiten.Image{}, []*ebiten.Image{})
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
	for i := 1; i < 11; i++ {
		zombieImage, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("png/male/zombie-dead%d.png", i), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
			return
		}
		zombieImages[2] = append(zombieImages[2], zombieImage)
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
	knightAttack := []*ebiten.Image{}
	for i := 1; i < 11; i++ {
		knightImage, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("png/knight/knight-attack%d.png", i), ebiten.FilterDefault)
		knightAttack = append(knightAttack, knightImage)
	}
	knightImages = append(knightImages, knightAttack)
	knightJump := []*ebiten.Image{}
	for i := 1; i < 11; i++ {
		knightImage, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("png/knight/knight-jump%d.png", i), ebiten.FilterDefault)
		knightJump = append(knightJump, knightImage)
	}
	knightImages = append(knightImages, knightJump)
	frameCounter = 1
}
