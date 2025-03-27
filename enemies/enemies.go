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

// location of enemy
type Enemy struct {
	X, Y, SpeedX, SpeedY, StartY float64

	Pathing Path

	StepCount int

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
