package enemies

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EnemyType struct {
	EType, ProjectileType, PathingType int

	SpeedX, SpeedY float64

	HealthBase int

	Image       *ebiten.Image
	EnemyPixels []byte
}

type EP struct {
	EnemyPixels []byte
}

type AllEP struct {
	AllEnemyPixels []EP
}

func makeEP(ep []byte) (eps EP) {
	eps.EnemyPixels = ep

	return
}

func (e *Enemies) GetAllEnemyPixels() (allEP AllEP, err error) {
	var enemyList []EnemyType
	enemyList, err = e.GetAllEnemyTypes()

	for _, ep := range enemyList {
		allEP.AllEnemyPixels = append(allEP.AllEnemyPixels, makeEP(ep.EnemyPixels))
	}

	return
}

func (e *Enemies) GetEnemyType(enemyTypeInt int) (ets EnemyType) {
	var (
		enemyTypes []EnemyType
		err        error
	)

	if enemyTypes, err = e.GetAllEnemyTypes(); err != nil {
		log.Fatal(err)
	}

	for _, et := range enemyTypes {
		if et.EType == enemyTypeInt {
			ets = et
			return
		}
	}

	return
}

func (e *Enemies) GetAllEnemyTypes() (enemyTypeList []EnemyType, err error) {
	enemyTypeList = append(enemyTypeList, makeEnemy(1, 1, 1, 2, 0, 2, EnemyOneURL))
	enemyTypeList = append(enemyTypeList, makeEnemy(2, 1, 2, 3, 0, 1, EnemyTwoURL))

	return
}

// EType, Projectile Type, Pathing Type, Speed X, Speed Y, Image URL
func makeEnemy(eType, projectileType, pathingType int, speedX, speedY float64, healthBase int, ImageURL string) (newEnemy EnemyType) {
	newEnemy.EType = eType
	newEnemy.ProjectileType = projectileType
	newEnemy.PathingType = pathingType
	newEnemy.SpeedX = speedX
	newEnemy.SpeedY = speedY
	newEnemy.HealthBase = healthBase
	newEnemy.Image = fetchEbitImage(ImageURL)
	newEnemy.EnemyPixels = makeEnemyPixels(newEnemy.Image)

	return
}

func fetchEbitImage(ImageURL string) (img *ebiten.Image) {
	var err error
	if img, _, err = ebitenutil.NewImageFromFile(ImageURL); err != nil {
		log.Fatal(err)
	}

	return
}

func makeEnemyPixels(img *ebiten.Image) (ep []byte) {
	enemyPoint := img.Bounds().Size()
	ep = make([]byte, 4*enemyPoint.X*enemyPoint.Y)
	img.ReadPixels(ep)
	return
}
