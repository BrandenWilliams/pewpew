package enemyprojectiles

import (
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

func (ep *EnemyProjectiles) spawnEnemyBullet(x, y float64) (newBullet Bullet, fired bool) {
	if rand.Float64() < 0.02 { // ~2% chance per frame
		newBullet.X = x
		newBullet.Y = y + 16
		fired = true
		return
	}

	return
}

func (ep *EnemyProjectiles) NewProjectile(x, y float64, projectileType int) (newProjectile EnemyProjectile, fired bool) {
	fired = false
	switch projectileType {
	case 1:
		newProjectile.EnemyBullet, fired = ep.spawnEnemyBullet(x, y)
		return
	case 2:
		// to do
		return
	}

	return
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

	ep.EnemyBullets = newEnemyBullets
}

// Remove enemy bullets old enemy bullets in and new enemy bullets.
func (ep *EnemyProjectiles) DespawnEnemyProjectiles() {
	ep.despawnEnemyBullets()
}
