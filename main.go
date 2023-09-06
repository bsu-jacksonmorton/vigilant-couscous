package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"
	"github.com/hajimehoshi/ebiten"
)

const (
	WINDOW_LIMIT_LEFT  float64 = -32
	WINDOW_LIMIT_RIGHT float64 = 206
	// Actor States
	IDLE           = 0
	WALKING        = 1
	ATTACKING      = 2
	PLAYER_JUMPING = 3
	FACE_LEFT      = -1
	FACE_RIGHT     = 1
	ZOMBIE_DEAD    = 2
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
	switch z.state {
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
		if z.frameCount <= 9 {
			z.Actor.SetFrameCount(z.Actor.frameCount + 0.2)
		}
		break
	default:
		panic("UNKNOWN ZOMBIE STATE!!!")
	}
}
func (z *Zombie) Kill() {
	z.state = ZOMBIE_DEAD
	z.SetFrameCount(0)
	z.Move(0, 25)
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
	vy            float64
	jumping       bool
	width         int
	height        int
	zombiesKilled int
}

func (p *Player) KillZombie() {
	p.zombiesKilled += 1
	fmt.Println("Zombies Killed: ", p.zombiesKilled)
}
func NewPlayer(posX, posY float64) *Player {
	return &Player{
		&Actor{1, 0, 0},
		&Position{posX, posY},
		0,
		false,
		74,
		90,
		0,
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
			p.SetState(PLAYER_JUMPING)
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
	p       *Player
	z       *Zombie
	zombies []*Zombie
}

func (g *Game) Update(*ebiten.Image) error {
	for _, z := range g.zombies {
		if collision(g.p, z) == true && g.p.state == ATTACKING && z.state != ZOMBIE_DEAD {
			z.Kill()
			g.p.KillZombie()
		}
		z.Update()
	}
	g.p.Update()
	return nil
}
func drawPlayer(screen *ebiten.Image, player *Player) {
	op := &ebiten.DrawImageOptions{}
	state := player.state
	frame := int(math.Floor(player.frameCount))
	w, _ := knightImages[state][frame].Bounds().Dx(), knightImages[state][frame].Bounds().Dy()
	//fmt.Println("Width:", w, "Height:", h)
	x, y := player.posX, player.posY
	op.GeoM.Scale(float64(player.direction), 1)
	op.GeoM.Translate(x+float64(w), y)
	screen.DrawImage(knightImages[state][frame], op)
}
func drawZombie(screen *ebiten.Image, z *Zombie) {
	op := &ebiten.DrawImageOptions{}
	frame := int(math.Floor(z.frameCount))
	w, _ := zombieImages[z.state][frame].Bounds().Dx(), zombieImages[z.state][frame].Bounds().Dy()
	op.GeoM.Scale(1, 1)
	op.GeoM.Translate(50+float64(w), z.posY)
	screen.DrawImage(zombieImages[z.state][frame], op)
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, &ebiten.DrawImageOptions{})
	drawPlayer(screen, g.p)
	for _, z := range g.zombies {
		drawZombie(screen, z)
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("my-first-game")
	player := NewPlayer(0, 150)
	zombie := NewZombie(25, 150)
	game := &Game{player, zombie, []*Zombie{}}
	go game.SpawnZombies()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func (g *Game) SpawnZombies() {
	for range time.Tick(time.Second * 5) {
		g.zombies = append(g.zombies, NewZombie(50, 150))
	}
}
func getImageFromPath(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}
func init() {
	// Load background
	var err error
	bgImg, err := getImageFromPath("png/test.png")
	if err != nil {
		log.Fatal(err)
    return
	}
	backgroundImage, err = ebiten.NewImageFromImage(bgImg, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
    return
	}
	// LOAD ZOMBIE FRAMES
	// zombie idle
	var zombieImage *ebiten.Image
	zombieImages = [][]*ebiten.Image{}
	zombieImages = append(zombieImages, []*ebiten.Image{}, []*ebiten.Image{}, []*ebiten.Image{})
	for i := 1; i < 16; i++ {
		tmp, err := getImageFromPath(fmt.Sprintf("png/male/idle%d.png", i))
		if err != nil {
			log.Fatal(err)
			return
		}
		zombieImage, err = ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
			return
		}
		zombieImages[0] = append(zombieImages[0], zombieImage)
	}
	// zombie walk
	for i := 1; i < 11; i++ {
		tmp, err := getImageFromPath(fmt.Sprintf("png/male/walk%d.png", i))
		if err != nil {
			log.Fatal(err)
			return
		}
		zombieImage, err = ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
			return
		}
		zombieImages[1] = append(zombieImages[1], zombieImage)
	}
	for i := 1; i < 11; i++ {
		tmp, err := getImageFromPath(fmt.Sprintf("png/male/zombie-dead%d.png", i))
		if err != nil {
			log.Fatal(err)
			return
		}
		zombieImage, err = ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
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
		tmp, err := getImageFromPath(fmt.Sprintf("png/knight/knight-real-idle%d.png", i))
		if err != nil {
			log.Fatal(err)
			return
		}
		knightImage, err = ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
		knightIdle = append(knightIdle, knightImage)
	}
	knightImages = append(knightImages, knightIdle)
	knightRun := []*ebiten.Image{}
	for i := 1; i < 11; i++ {
		tmp, err := getImageFromPath(fmt.Sprintf("png/knight/knight-run%d.png", i))
		if err != nil {
			log.Fatal(err)
			return
		}
		knightImage, err = ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
		knightRun = append(knightRun, knightImage)
	}
	knightImages = append(knightImages, knightRun)
	knightAttack := []*ebiten.Image{}
	for i := 1; i < 11; i++ {
		tmp, err := getImageFromPath(fmt.Sprintf("png/knight/knight-attack%d.png", i))
		if err != nil {
			log.Fatal(err)
			return
		}
		knightImage, err = ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
		knightAttack = append(knightAttack, knightImage)
	}
	knightImages = append(knightImages, knightAttack)
	knightJump := []*ebiten.Image{}
	for i := 1; i < 11; i++ {
		tmp, err := getImageFromPath(fmt.Sprintf("png/knight/knight-jump%d.png", i))
		if err != nil {
			log.Fatal(err)
			return
		}
		knightImage, err = ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
		knightJump = append(knightJump, knightImage)
	}
	knightImages = append(knightImages, knightJump)
	frameCounter = 1
}
