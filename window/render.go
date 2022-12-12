// window/reader.go
// Handles the OpenGL rendering pipeline.

package window

import (
	"github.com/go-gl/gl/v2.1/gl"
)

// Render pipeline.
type RenderPipeline struct {
	objects    []*Object
	newObjects chan *Object
}

// New render pipeline.
func NewRenderPipeline() *RenderPipeline {
	return &RenderPipeline{
		objects:    []*Object{},
		newObjects: make(chan *Object, 100),
	}
}

// Add a new object.
func (r *RenderPipeline) RegisterObject(o *Object) {
	r.newObjects <- o
}

// Object.
type Object struct {
	points           []float32
	vao              uint32
	x                float32
	y                float32
	spName           string
	window           *Window
	shaderAttributes map[*uint8]float32
}

// New object.
func NewObject(points []float32, spName string, w *Window) *Object {
	// Return the new object.
	return &Object{
		points:           points,
		spName:           spName,
		window:           w,
		shaderAttributes: map[*uint8]float32{},
	}
}

// Set the position of the object.
func (o *Object) Move(x, y float32) {
	o.x = x
	o.y = y
}

// Get the position of the object.
func (o *Object) Pos() (float32, float32) {
	return o.x, o.y
}

// Set a shader attribute.
func (o *Object) SetShaderAttribute(name string, value float32) {
	nameString := gl.Str(name + "\x00")
	o.shaderAttributes[nameString] = value
}

// Generate a VBO from the object.
func (o *Object) GenerateVAO() {
	// Generate the VBO first.
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(o.points), gl.Ptr(o.points), gl.STATIC_DRAW)

	// Generate the VAO.
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	o.vao = vao
}

// Render the object.
func (o *Object) Render() {
	// Set the shader program to use.
	shader := o.window.shaderPrograms[o.spName]
	shaderProg := shader.program
	gl.UseProgram(shaderProg)

	// Get the location of the vertex shader position uniform param.
	loc := gl.GetUniformLocation(shaderProg, shader.positionString)
	gl.Uniform3f(loc, o.x, o.y, 1)

	// Set the shader attributes.
	for attribute := range o.shaderAttributes {
		loc := gl.GetUniformLocation(shaderProg, attribute)
		gl.Uniform1f(loc, o.shaderAttributes[attribute])
	}

	gl.BindVertexArray(o.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(o.points)/3))
}
