package enemies

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EnemyType struct {
	EType       int
	Image       *ebiten.Image
	EnemyPixels []byte
}

type EnemyTypes struct {
	EnemyTypeStuct EnemyType
}

func (et *EnemyType) makeEnemyPixels(img *ebiten.Image) (ep []byte) {
	enemyPoint := img.Bounds().Size()
	ep = make([]byte, 4*enemyPoint.X*enemyPoint.Y)
	img.ReadPixels(ep)
	return
}

func (et *EnemyType) TypeOne() (ret EnemyType) {
	// enemy sprite
	var (
		err error
	)

	ret.Image, _, err = ebitenutil.NewImageFromFile(EnemyOneURL)
	if err != nil {
		log.Fatal(err)
	}

	ret.EnemyPixels = et.makeEnemyPixels(ret.Image)

	return
}
