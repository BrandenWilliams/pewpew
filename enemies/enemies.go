package enemies

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 960
	ScreenHight  = 540
	EnemyShipURL = "art/enemyship.png"
)

// Bullet represents the location of a single shot
type Bullet struct {
	X, Y float64
}

// location of enemy
type Enemy struct {
	X, Y float64

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

func (e *Enemies) createEnemy() (newEnemy Enemy) {
	// enemy Sprite
	var err error
	newEnemy.EnemyImage, _, err = ebitenutil.NewImageFromFile(EnemyShipURL)
	if err != nil {
		log.Fatal(err)
	}

	newEnemy.EnemyPixels = e.ET.makeEnemyPixels(newEnemy.EnemyImage)

	newEnemy.X, newEnemy.Y = CreateEnemyLocation(newEnemy.EnemyImage.Bounds().Dx(), newEnemy.EnemyImage.Bounds().Dy())

	return
}

func (e *Enemies) GenerateEnemy(enemyType int) (newEnemy Enemy) {
	switch enemyType {
	case 1:
		newEnemy = e.createEnemy()
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
		newEnemy := e.GenerateEnemy(1)
		e.ES = append(e.ES, newEnemy)
		e.enemySpawnTimer -= spawnInterval // Subtract instead of resetting to 0
	}
}

func (e *Enemies) EnemyMovement() {
	// Move enemies left
	for i := range e.ES {
		e.ES[i].X -= 2 // Adjust speed as needed
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
