package enemyprojectiles

import (
	"log"
	"math/rand/v2"
)

// Bullet represents the location of a single shot
type Bullet struct {
	X, Y float64
}

type EnemyProjectiles struct {
	EnemyBullets []Bullet
}

type EnemyProjectile struct {
	ProjectileType int

	EnemyBullet Bullet
}

func (ep *EnemyProjectiles) SpawnEnemyBullet(x, y float64) (newBullet Bullet) {
	if rand.Float64() < 0.02 { // ~2% chance per frame
		newBullet.X = x
		newBullet.Y = y + 16
		ep.EnemyBullets = append(ep.EnemyBullets, newBullet)

		return
	}

	return
}

func (ep *EnemyProjectiles) NewProjectile(x, y float64, projectileType int) {
	switch projectileType {
	case 1:
		ep.SpawnEnemyBullet(x, y)
		return
	case 2:
		log.Printf("BAD BAD BAD projectileType: %v\n", projectileType)
		// to do
		return
	default:
		return
	}
}

func (ep *EnemyProjectiles) ManageEnemyProjectiles() {
	// move bullets
	for j := range ep.EnemyBullets {
		ep.EnemyBullets[j].X -= 3
	}
}

func (ep *EnemyProjectiles) despawnEnemyBullets() {
	var newEnemyBullets []Bullet
	for _, eb := range ep.EnemyBullets {
		if eb.X > 0 {
			newEnemyBullets = append(newEnemyBullets, eb)
		}
	}
	log.Printf("bullet size: %v\n", len(newEnemyBullets))
	ep.EnemyBullets = newEnemyBullets
}

// Remove enemy bullets old enemy bullets in and new enemy bullets.
func (ep *EnemyProjectiles) DespawnEnemyProjectiles() {
	ep.despawnEnemyBullets()
}
