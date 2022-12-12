package main

import (
	"balls/physics"
	"balls/vec"
	"balls/window"
	"math/rand"
	"time"
)

const (
	vertexShaderSource = `
		#version 410
		layout (location=0) in vec3 vp;
		out vec2 texCoordV;
		uniform vec3 position;
		uniform float radius;
		void main() {
			texCoordV = (vp.xy+1)/2;
			gl_Position = vec4(vp * radius * 0.1 + position, 1.0f);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 frag_colour;
		in vec2 texCoordV;

		void main() {
			vec2 center = vec2(0.5, 0.5);
			float dist = length(texCoordV-center);
			if (dist < 0.5) {
				frag_colour = vec4(texCoordV, 0.0, 1.0);
			} else {
				discard;
			}
		}
	` + "\x00"
)

type MyDrawContext struct {
	Window         *window.Window
	Solver         *physics.VerletSolver
	Objects        []Circle
	lastObjectTime time.Time
}

type Circle struct {
	GLObj *window.Object
	VObj  *physics.VerletObject
}

func (c *MyDrawContext) PreDraw(elapsedTime time.Duration) {
	now := time.Now()
	if now.Sub(c.lastObjectTime) >= time.Second {
		c.lastObjectTime = now
		// We should add a new object.
		square := []float32{
			-1, 1, 0,
			-1, -1, 0,
			1, -1, 0,

			-1, 1, 0,
			1, 1, 0,
			1, -1, 0,
		}
		radius := rand.Float32() * 2
		glObj := window.NewObject(square, "circle", c.Window)
		glObj.SetShaderAttribute("radius", radius)
		c.Window.RegisterObject(glObj)

		vObj := physics.NewVerletObject(vec.NewVec2(-1, 0), float64(radius))
		c.Solver.AddObject(vObj)
		vObj.SetVelocity(vec.NewVec2(0, -1), elapsedTime.Seconds())
		circle := Circle{
			GLObj: glObj,
			VObj:  vObj,
		}
		c.Objects = append(c.Objects, circle)
	}
	c.Solver.Update()
	for i := range c.Objects {
		c.Objects[i].GLObj.Move(float32(c.Objects[i].VObj.CurrentPos.X)/10, float32(c.Objects[i].VObj.CurrentPos.Y)/10)
	}
}

func main() {
	dc := &MyDrawContext{}
	dc.lastObjectTime = time.Now()

	s := window.NewShaderProgram(vertexShaderSource, fragmentShaderSource, "position")
	w := window.New(500, 500, "hi", dc)
	w.RegisterShaderProgram("circle", s)
	w.Init()
	dc.Window = w
	dc.Objects = []Circle{}
	dc.Solver = physics.NewVerletSolver(vec.NewVec2(0, -1), 8)
	dc.Solver.SetConstraint(vec.NewVec2(0, 0), 10)
	// add an object
	w.Start()
}
