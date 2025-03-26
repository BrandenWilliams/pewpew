package enemies

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth = 960
	ScreenHight = 540
	EnemyOneURL = "art/enemyship.png"
	EnemyTwoURL = "art/enemyTwoShip.png"
)

// Bullet represents the location of a single shot
type Bullet struct {
	X, Y float64
}

// location of enemy
type Enemy struct {
	X, Y, SpeedX, SpeedY, StartY float64

	pathing Path

	StepCount   int
	ZigDuration int

	EnemyType int

	EnemyImage  *ebiten.Image
	EnemyPixels []byte
}

type Enemies struct {
	ES []Enemy

	EnemyBullets []Bullet

	enemySpawnTimer   float64
	initialSpawnDelay float64

	ET EnemyType
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
		if e.ES[i].StepCount >= len(et.pathing.Cords) {
			continue
		}

		e.ES[i] = e.NextStep(et)
	}
}

func (e *Enemies) ManageEnemyBullets() {
	// Make each enemy fire a bullet every 2 seconds (adjust as needed)
	for _, enemySlice := range e.ES {

		if rand.Float64() < 0.02 { // ~2% chance per frame

			e.EnemyBullets = append(e.EnemyBullets, Bullet{X: enemySlice.X, Y: enemySlice.Y + 16})
		}

	}

	// move bullets
	for j := range e.EnemyBullets {
		e.EnemyBullets[j].X -= 3
	}
}
