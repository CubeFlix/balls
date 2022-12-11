// window/window.go
// Package window provides tools for using OpenGL windows.

package window

import (
	"balls/context"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Simulation window struct.
type Window struct {
	Width  int
	Height int
	Title  string

	DrawContext context.DrawContext

	renderPipeline *RenderPipeline

	isReady        bool
	window         *glfw.Window
	shaderPrograms map[string]*ShaderProgram
}

// Create a new window.
func New(width, height int, title string, drawContext context.DrawContext) *Window {
	return &Window{
		Width:          width,
		Height:         height,
		Title:          title,
		DrawContext:    drawContext,
		renderPipeline: NewRenderPipeline(),
		shaderPrograms: map[string]*ShaderProgram{},
	}
}

// Initialize the window.
func (w *Window) Init() {
	// We need to be on the main thread.
	runtime.LockOSThread()

	// Init GLFW.
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	// Set the values on the window.
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create the window.
	var err error
	w.window, err = glfw.CreateWindow(w.Width, w.Height, w.Title, nil, nil)
	if err != nil {
		panic(err)
	}
	w.window.MakeContextCurrent()

	// Init OpenGL.
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("starting OpenGL version", version)

	// Prepare the shader programs.
	for name := range w.shaderPrograms {
		w.shaderPrograms[name].Prepare()
	}

	log.Println("finished initialization")

	w.isReady = true
}

// Start the window.
func (w *Window) Start() {
	if !w.isReady {
		panic("window not initialized")
	}

	// Main loop.
	for !w.window.ShouldClose() {
		w.DrawContext.PreDraw()
		w.draw()
		w.DrawContext.PostDraw()
	}

	// Terminate the window.
	glfw.Terminate()
}

// Draw function.
func (w *Window) draw() {
	// Load new objects.
	shouldContinueChecking := true
	for shouldContinueChecking {
		select {
		case obj := <-w.renderPipeline.newObjects:
			// New object, generate the VAO and add it to the render pipeline's objects.
			obj.GenerateVAO()
			w.renderPipeline.objects = append(w.renderPipeline.objects, obj)
		default:
			// No new objects this cycle, continue.
			shouldContinueChecking = false
		}
	}

	// Prepare the window.
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Render each object in the pipeline.
	for i := range w.renderPipeline.objects {
		w.renderPipeline.objects[i].Render()
	}

	// Handle events and draw.
	glfw.PollEvents()
	w.window.SwapBuffers()
}

// Register a new object.
func (w *Window) RegisterObject(o *Object) {
	w.renderPipeline.RegisterObject(o)
}

// Register a new shader program. Should be done BEFORE initialization.
func (w *Window) RegisterShaderProgram(name string, prog *ShaderProgram) {
	w.shaderPrograms[name] = prog
}
