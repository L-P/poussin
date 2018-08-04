package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/gobuffalo/packr"
	"home.leo-peltier.fr/poussin/render"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not init GLFW: %s", err))
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "Poussin", nil, nil)
	if err != nil {
		panic(fmt.Errorf("could not create window: %s", err))
	}

	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("could not init OpenGL: %s", err))
	}

	glfw.SwapInterval(1)

	box := packr.NewBox("./assets/shaders")
	program, err := render.LoadProgram(
		box.String("default.vert"),
		box.String("default.frag"),
	)
	if err != nil {
		panic(fmt.Errorf("could not load shaders: %s", err))
	}

	for !window.ShouldClose() {
		gl.Disable(gl.DEPTH_TEST)
		gl.Disable(gl.BLEND)
		gl.Disable(gl.CULL_FACE)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(program)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
