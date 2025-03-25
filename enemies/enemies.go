package enemies

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func (e *Enemies) createFirstEnemy() (newEnemy Enemy) {
	// enemy Sprite
	var err error
	newEnemy.EnemyType = 1
	newEnemy.EnemyImage, _, err = ebitenutil.NewImageFromFile(EnemyOneURL)
	if err != nil {
		log.Fatal(err)
	}
	// Generate "hitbox" from sprite
	newEnemy.EnemyPixels = e.ET.makeEnemyPixels(newEnemy.EnemyImage)

	newEnemy.SpeedX = 2
	newEnemy.SpeedY = 0
	newEnemy.StepCount = 1

	newEnemy.X, newEnemy.Y = CreateEnemyLocation(newEnemy.EnemyImage.Bounds().Dx(), newEnemy.EnemyImage.Bounds().Dy())
	newEnemy.pathing = e.StraightAhead(newEnemy.X, newEnemy.Y, newEnemy.SpeedX)
	return
}

func (e *Enemies) createSecondEnemy() (newEnemy Enemy) {
	// enemy Sprite
	var err error
	newEnemy.EnemyType = 2

	newEnemy.EnemyImage, _, err = ebitenutil.NewImageFromFile(EnemyTwoURL)
	if err != nil {
		log.Fatal(err)
	}
	// Generate "hitbox"
	newEnemy.EnemyPixels = e.ET.makeEnemyPixels(newEnemy.EnemyImage)

	newEnemy.SpeedX = 3
	newEnemy.StepCount = 1

	newEnemy.X, newEnemy.Y = CreateEnemyLocation(newEnemy.EnemyImage.Bounds().Dx(), newEnemy.EnemyImage.Bounds().Dy())
	newEnemy.pathing = e.GenerateZigzagPath(newEnemy.X, newEnemy.Y, newEnemy.SpeedX, 80.0, 0.02)

	return
}

func (e *Enemies) decideEnemyType() (enemyType int) {
	rando := rand.Intn(100)

	if rando > 95 {
		enemyType = 2
		return
	} else if rando > 80 {
		enemyType = 2
		return
	} else {
		enemyType = 1
		return
	}
}

func (e *Enemies) GenerateEnemy() (newEnemy Enemy) {
	// enemyType := e.decideEnemyType()
	enemyType := 2

	switch enemyType {
	case 1:
		newEnemy = e.createFirstEnemy()
	case 2:
		newEnemy = e.createSecondEnemy()
	}

	return
}

func (e *Enemies) EnemySpawn() {
	tps := ebiten.ActualTPS()
	// Ensure TPS is valid (avoid dividing by zero)
	if tps == 0 {
		tps = 60
	}

	// Start a delay before first wave
	if e.enemySpawnTimer < e.initialSpawnDelay {
		e.enemySpawnTimer += 1 / tps
		return
	}

	// Increment timer after the initial delay
	e.enemySpawnTimer += 1 / tps

	// Spawn a new enemy every 1.5 seconds
	spawnInterval := 1.5
	if e.enemySpawnTimer >= spawnInterval {
		newEnemy := e.GenerateEnemy()
		e.ES = append(e.ES, newEnemy)
		e.enemySpawnTimer -= spawnInterval // Subtract instead of resetting to 0
	}
}

func (e *Enemies) EnemiesMovement() {
	// bail in no mobs have spawned
	if len(e.ES) <= 0 {
		return
	}

	for i, et := range e.ES {
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
