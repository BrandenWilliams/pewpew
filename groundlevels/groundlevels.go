package groundlevels

const (
	levelOneBase_seg1 = "art/groundlevels/levelOneBase.png"
	levelOneBase_seg2 = "art/groundlevels/levelOneBase_part2.png"
)

type GroundLevels struct {
	GroundImagesURLs []string
	CurrentBGIndex   int
}

func (gl *GroundLevels) SetFirstLevelBackgrounds() {
	gl.GroundImagesURLs = append(gl.GroundImagesURLs, levelOneBase_seg1)
	gl.GroundImagesURLs = append(gl.GroundImagesURLs, levelOneBase_seg2)
	gl.GroundImagesURLs = append(gl.GroundImagesURLs, levelOneBase_seg1)
	gl.GroundImagesURLs = append(gl.GroundImagesURLs, levelOneBase_seg2)
}
