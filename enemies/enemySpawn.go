package enemies

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func (e *Enemies) CreateEnemy(enemyType EnemyType) (eStruct Enemy) {

	eStruct.EnemyImage = enemyType.Image

	eStruct.EnemyType = enemyType.EType
	eStruct.ProjectileType = enemyType.ProjectileType
	eStruct.EnemyPixels = enemyType.EnemyPixels

	eStruct.SpeedX = enemyType.SpeedX
	eStruct.SpeedY = enemyType.SpeedY

	eStruct.StepCount = 1

	eStruct.X, eStruct.Y = CreateEnemyLocation(enemyType.Image.Bounds().Dx(), enemyType.Image.Bounds().Dy())
	eStruct.Pathing = e.GeneratePath(enemyType.PathingType, eStruct.X, eStruct.Y, eStruct.SpeedX, eStruct.SpeedY)
	return
}

func (e *Enemies) GenerateEnemy() (newEnemy Enemy) {
	enemyTypeInt := e.DecideEnemyType(1)
	enemyTypeStruct := e.GetEnemyType(enemyTypeInt)
	newEnemy = e.CreateEnemy(enemyTypeStruct)

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

func (e *Enemies) DecideEnemyType(spawnSetup int) (enemyType int) {
	switch spawnSetup {
	case 1:
		enemyType = levelOneSpawn()
	}

	return
}

func levelOneSpawn() (enemyType int) {
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
