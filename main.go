package main

import (
	"balls/window"
)

const (
	vertexShaderSource = `
		#version 410
		in vec3 vp;
		uniform vec3 position;
		void main() {
			gl_Position = vec4(vp + position, 1.0f);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"
)

type MyDrawContext struct {
	
}

func main() {
	square := []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}

	d := 
	
	s := window.NewShaderProgram(vertexShaderSource, fragmentShaderSource, "position")
	w := window.New(500, 500, "hi")
	w.RegisterShaderProgram("basic", s)
	w.Init()
	// add an object
	o := window.NewObject(square, "basic", w)
	w.RegisterObject(o)

	w.Start()
}
