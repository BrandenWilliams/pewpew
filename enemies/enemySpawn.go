package enemies

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (e *Enemies) GenerateEnemy() (newEnemy Enemy) {
	enemyType := e.decideEnemyType()

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

func (e *Enemies) EnemyDespawn() {
	var newEnemies []Enemy

	for _, e := range e.ES {
		if e.X > -150 {
			newEnemies = append(newEnemies, e)
		}
	}

	e.ES = newEnemies
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
