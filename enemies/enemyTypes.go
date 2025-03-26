package enemies

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EnemyType struct {
	EType, ProjectileType, PathingType int

	SpeedX, SpeedY float64

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
	enemyTypeList = append(enemyTypeList, makeEnemy(1, 1, 1, 2, 0, EnemyOneURL))
	enemyTypeList = append(enemyTypeList, makeEnemy(2, 1, 2, 3, 0, EnemyTwoURL))

	return
}

// EType, Projectile Type, Pathing Type, Speed X, Speed Y, Image URL
func makeEnemy(eType, projectileType, pathingType int, speedX, speedY float64, ImageURL string) (newEnemy EnemyType) {
	newEnemy.EType = eType
	newEnemy.ProjectileType = projectileType
	newEnemy.PathingType = pathingType
	newEnemy.SpeedX = speedX
	newEnemy.SpeedY = speedY
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

/*
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
*/
