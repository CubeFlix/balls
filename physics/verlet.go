// physics/verlet.go
// A Verlet physics engine.

package physics

import (
	"balls/vec"
	"math"
	"time"
)

// A verlet object.
type VerletObject struct {
	CurrentPos vec.Vec2
	OldPos     vec.Vec2
	Accel      vec.Vec2
	Radius     float64
}

// Create a new Verlet object
func NewVerletObject(pos vec.Vec2, radius float64) *VerletObject {
	return &VerletObject{
		CurrentPos: pos,
		OldPos:     pos,
		Radius:     radius,
	}
}

// Update the position.
func (v *VerletObject) Update(dt float64) {
	// Calculate displacement.
	displacement := v.CurrentPos.Sub(v.OldPos)
	v.OldPos = v.CurrentPos

	// Calculate Verlet.
	v.CurrentPos = vec.AddVec2(v.CurrentPos, displacement, v.Accel.MulScalar(dt))
	v.Accel = vec.Vec2{X: 0, Y: 0}
}

// Update the acceleration.
func (v *VerletObject) UpdateAccel(accel vec.Vec2) {
	v.Accel = v.Accel.Add(accel)
}

func (v *VerletObject) SetVelocity(vel vec.Vec2, dt float64) {
	v.CurrentPos = v.CurrentPos.Sub(vel.MulScalar(dt))
}

// A Verlet solver.
type VerletSolver struct {
	Gravity  vec.Vec2
	Objects  []*VerletObject
	SubSteps int

	ConstraintCenter vec.Vec2
	ConstraintRadius float64

	lastCalculatedTime time.Time
}

// Create a new Verlet solver.
func NewVerletSolver(gravity vec.Vec2, substeps int) *VerletSolver {
	return &VerletSolver{
		Gravity:  gravity,
		SubSteps: substeps,
	}
}

// Set the constraint.
func (v *VerletSolver) SetConstraint(center vec.Vec2, radius float64) {
	v.ConstraintCenter = center
	v.ConstraintRadius = radius
}

// Add an object.
func (v *VerletSolver) AddObject(o *VerletObject) {
	v.Objects = append(v.Objects, o)
}

// Set the last calculated time.
func (v *VerletSolver) SetLastCalculatedTime(now time.Time) {
	v.lastCalculatedTime = now
}

// Update the solver.
func (v *VerletSolver) Update() {
	// Perform each sub step calculation.
	for i := 0; i < v.SubSteps; i++ {
		// Calculate the delta time.
		now := time.Now()
		dt := now.Sub(v.lastCalculatedTime)

		// Apply gravity.
		for i := range v.Objects {
			v.Objects[i].UpdateAccel(v.Gravity)
		}

		// Check for collisions.
		v.CheckCollisions(dt.Seconds())

		// Apply the constraint.
		v.ApplyConstraint()

		// Apply updates.
		for i := range v.Objects {
			v.Objects[i].Update(dt.Seconds())
		}

		v.SetLastCalculatedTime(now)
	}
}

// Apply the constraint to each object.
func (v *VerletSolver) ApplyConstraint() {
	for i := range v.Objects {
		// Apply the constraint to each object.
		offset := v.ConstraintCenter.Sub(v.Objects[i].CurrentPos)
		dist := math.Sqrt(offset.X*offset.X + offset.Y*offset.Y)
		if dist > (v.ConstraintRadius - v.Objects[i].Radius) {
			n := offset.MulScalar(1 / dist)
			v.Objects[i].CurrentPos = v.ConstraintCenter.Sub(n.MulScalar(v.ConstraintRadius - v.Objects[i].Radius))
		}
	}
}

// Check for collisions.
func (v *VerletSolver) CheckCollisions(dt float64) {
	responseCoef := 0.1
	nObjects := len(v.Objects)

	// Apply collisions for each object.
	for i := 0; i < nObjects; i++ {
		object1 := v.Objects[i]

		// Iterate over the second object.
		for k := i + 1; k < nObjects; k++ {
			object2 := v.Objects[k]
			offset := object1.CurrentPos.Sub(object2.CurrentPos)
			dist2 := offset.X*offset.X + offset.Y*offset.Y
			minDist := object1.Radius + object2.Radius
			if dist2 < minDist*minDist {
				// Check for overlapping.
				dist := math.Sqrt(dist2)
				n := offset.MulScalar(1 / dist)
				massRatio1 := object1.Radius / (object1.Radius + object2.Radius)
				massRatio2 := object2.Radius / (object1.Radius + object2.Radius)
				delta := 0.5 * responseCoef * (dist - minDist)
				// Update the positions.
				object1.CurrentPos = object1.CurrentPos.Sub(n.MulScalar(massRatio2 * delta))
				object2.CurrentPos = object2.CurrentPos.Add(n.MulScalar(massRatio1 * delta))
			}
		}
	}
}
