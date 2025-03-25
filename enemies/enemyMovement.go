package enemies

import (
	"log"
	"math"
)

type Path struct {
	Cords []CordSet
}

type CordSet struct {
	Step int
	X, Y float64
}

func (e *Enemies) NextStep(eIn Enemy) Enemy {
	newEnemy := eIn

	newEnemy.X = eIn.pathing.Cords[eIn.StepCount].X
	newEnemy.Y = eIn.pathing.Cords[eIn.StepCount].Y
	newEnemy.StepCount = eIn.StepCount + 1

	return newEnemy
}

// Mobs go Stright only
func (e *Enemies) StraightAhead(startX, fixedY, speedX float64) Path {
	var path Path
	stepCount := 0

	for x := startX; x > -50; x -= speedX {
		stepCount++
		path.Cords = append(path.Cords, CordSet{Step: stepCount, X: x, Y: fixedY})
	}

	return path
}

// First mob using **temp note**
// amplitude := 80.0  frequency := 0.02
func (e *Enemies) GenerateZigzagPath(startX, startY, speedX, amplitude, frequency float64) Path {
	var path Path
	stepCount := 0

	// Loop until the enemy moves off-screen
	for x := startX; x > -50; x -= speedX {
		stepCount++
		// create y cord using Sin controlled by amplitude and freq
		y := startY + amplitude*math.Sin(float64(stepCount)*frequency)
		// append new set of cords along with step count
		path.Cords = append(path.Cords, CordSet{Step: stepCount, X: x, Y: y})
	}

	return path
}

// FOR DEBUGGING PATHING KEEP FOR LATER PATHS
func PrintPathing(pathing Path) {
	log.Printf("Pathing START Len count: %v\n", len(pathing.Cords))

	for s, c := range pathing.Cords {
		log.Printf("Cords step: %v, stepCount: %v, X: %v, Y: %v", s, c.Step, c.X, c.Y)
	}

	log.Printf("Pathing STOP \n")

}
