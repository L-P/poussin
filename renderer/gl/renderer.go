package gl

import (
	"fmt"
	"image"
	"runtime"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/L-P/poussin/emu/ppu"
)

type Renderer struct {
	window *glfw.Window

	program uint32
	vao     uint32
	texture uint32
}

func New() (*Renderer, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("could not init GLFW: %s", err)
	}

	window := initWindow()

	program, err := LoadProgram(shaderDefaultVert, shaderDefaultFrag)
	if err != nil {
		return nil, err
	}

	vao, _ := getPlane(program)
	texture := createTexture(ppu.DotMatrixWidth, ppu.DotMatrixHeight)

	return &Renderer{
		window:  window,
		program: program,
		vao:     vao,
		texture: texture,
	}, nil
}

func (r *Renderer) Close() {
	runtime.LockOSThread()
	glfw.Terminate()
	runtime.UnlockOSThread()
}

func (r *Renderer) Run(nextFrame <-chan *image.RGBA, shouldClose <-chan bool, closed chan<- bool) {
	runtime.LockOSThread()

	gl.UseProgram(r.program)
	projection := mgl32.Ortho2D(0, 1, 0, 1)
	projectionUniform := gl.GetUniformLocation(r.program, gl.Str("uProjection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	for !r.window.ShouldClose() {
		width, height := r.window.GetSize()
		gl.Viewport(0, 0, int32(width), int32(height))

		resetGLState()
		gl.Clear(gl.COLOR_BUFFER_BIT)

		select {
		case fb := <-nextFrame:
			updateTexture(r.texture, fb)
		case <-shouldClose:
			r.window.SetShouldClose(true)
		default:
		}

		drawPlane(r.program, r.vao, r.texture)

		r.window.SwapBuffers()
		glfw.PollEvents()
	}

	runtime.UnlockOSThread()
	close(closed)
}

func drawPlane(program, vao, texture uint32) {
	gl.UseProgram(program)
	gl.BindVertexArray(vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.DrawArrays(gl.TRIANGLES, 0, 3*2)
}

func resetGLState() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)
	gl.ClearColor(0, 0, 0, 1)
}

func getPlane(program uint32) (vao, vbo uint32) {
	//  x  y uvx uvy
	vertices := []float32{
		0, 0, 0, 0,
		1, 0, 1, 0,
		0, 1, 0, 1,

		1, 0, 1, 0,
		1, 1, 1, 1,
		0, 1, 0, 1,
	}

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("inPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("inTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	return vao, vbo
}

func createTexture(width, height int32) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		width,
		height,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(nil),
	)

	return texture
}

func updateTexture(texture uint32, rgba *image.RGBA) {
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexSubImage2D(
		gl.TEXTURE_2D,
		0, 0, 0,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
}

func initWindow() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(640, 576, "Poussin", nil, nil)
	if err != nil {
		panic(fmt.Errorf("could not create window: %s", err))
	}

	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("could not init OpenGL: %s", err))
	}

	// Don't show garbage if the first frame does not come soon enough
	resetGLState()
	gl.Clear(gl.COLOR_BUFFER_BIT)
	window.SwapBuffers()

	glfw.SwapInterval(1)

	return window
}
