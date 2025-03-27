package enemies

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth = 960
	ScreenHight = 540
	EnemyOneURL = "art/enemyship.png"
	EnemyTwoURL = "art/enemyTwoShip.png"
)

// location of enemy
type Enemy struct {
	X, Y, SpeedX, SpeedY, StartY float64

	Pathing Path

	StepCount int

	MaxHealth     int
	CurrentHealth int

	EnemyType      int
	ProjectileType int

	EnemyImage  *ebiten.Image
	EnemyPixels []byte
}

type Enemies struct {
	ES []Enemy

	enemySpawnTimer   float64
	initialSpawnDelay float64
}

// long term make this support all enemys
func CreateEnemyLocation(enemyWidth, enemyHeight int) (x, y float64) {
	x = ScreenWidth - 1
	y = float64(rand.Intn(ScreenHight - enemyHeight))
	return
}

func (e *Enemies) EnemiesMovement() {
	// bail if no mobs have spawned
	if len(e.ES) <= 0 {
		return
	}

	for i, et := range e.ES {
		// bail if past last step
		if e.ES[i].StepCount >= len(et.Pathing.Cords) {
			continue
		}

		e.ES[i] = e.NextStep(et)
	}
}

func (g *Enemies) DrawHealthBar(e Enemy, screen *ebiten.Image) {
	barWidth := e.EnemyImage.Bounds().Dx()
	barHeight := 5.0 // Height of the health bar

	// Calculate the current width based on health
	healthPercentage := float64(e.CurrentHealth) / float64(e.MaxHealth)
	//log.Printf("health %: %v, currentHealth %v, maxHealth %v", healthPercentage, e.CurrentHealth, e.MaxHealth)
	currentBarWidth := float64(e.EnemyImage.Bounds().Dx()) * healthPercentage

	// Create a red background bar (full size)
	healthBarBg := ebiten.NewImage(int(barWidth), int(barHeight))
	healthBarBg.Fill(color.RGBA{255, 0, 0, 255}) // Red for missing health

	// Draw the background first
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(e.X, e.Y)
	screen.DrawImage(healthBarBg, opts)

	var healthBar *ebiten.Image
	if currentBarWidth > 0 {
		// Create a green foreground bar (based on current health)
		healthBar = ebiten.NewImage(int(currentBarWidth), int(barHeight))
		healthBar.Fill(color.RGBA{0, 255, 0, 255}) // Green for remaining health

		// Draw the foreground (current health)
		opts.GeoM.Reset()
		opts.GeoM.Translate(e.X, e.Y)
		screen.DrawImage(healthBar, opts)
	}
}
