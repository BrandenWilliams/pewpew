package enemies

import "math"

func (g *Enemies) enemyMovement(eIn Enemy) Enemy {

	switch eIn.EnemyType {
	case 1:
		eIn.X = g.straightAhead(eIn.X, eIn.SpeedX)
	case 2:
		eIn.X, eIn.Y, eIn.SpeedY, eIn.StepCount = g.zigzagPattern(eIn.X, eIn.Y, eIn.SpeedX, eIn.SpeedY, eIn.StartY, eIn.StepCount)
	}

	return eIn

}

func (g *Enemies) zigzagPattern(inX, inY, speedX, speedY, startY float64, stepCount int) (newX, newY, newSpeedY float64, newStepCount int) {
	amplitude := 80.0
	frequency := 0.02

	// Move horizontally and oscillate vertically
	newX = inX - speedX
	newY = startY + amplitude*math.Sin(float64(stepCount)*frequency)

	stepCount++
	newStepCount = stepCount

	// Return everything, including new speedY to track the direction change
	return newX, newY, speedY, stepCount
}

func (g *Enemies) straightAhead(inX, speedX float64) (newX float64) {
	newX = inX - speedX
	return
}
