// window/shaders.go
// Shaders for OpenGL.

package window

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

// GLSL shader program, with vertex and fragment shaders.
type ShaderProgram struct {
	VertexShaderSource   string
	FragmentShaderSource string
	PositionParamName    string
	program              uint32
	vertexShader         uint32
	fragmentShader       uint32
	positionString       *uint8
}

// New shader program.
func NewShaderProgram(vertex, fragment string, positionParamName string) *ShaderProgram {
	return &ShaderProgram{
		VertexShaderSource:   vertex,
		FragmentShaderSource: fragment,
		PositionParamName:    positionParamName,
	}
}

// Prepare the shader program and compile the shaders.
func (s *ShaderProgram) Prepare() {
	// Compile the shaders.
	var err error
	s.vertexShader, err = CompileShader(s.VertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	s.fragmentShader, err = CompileShader(s.FragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	// Create the program.
	s.program = gl.CreateProgram()
	gl.AttachShader(s.program, s.vertexShader)
	gl.AttachShader(s.program, s.fragmentShader)
	gl.LinkProgram(s.program)

	s.positionString = gl.Str(s.PositionParamName + "\x00")
}

// Compile a shader from source.
func CompileShader(source string, shaderType uint32) (uint32, error) {
	// Create the shader object.
	shader := gl.CreateShader(shaderType)

	// Load the source.
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	// Try to compile the shader.
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	// Return the completed shader
	return shader, nil
}
