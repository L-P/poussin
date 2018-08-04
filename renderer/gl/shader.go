package gl

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
)

/// Load compiles both frag/vert GL shaders and returns a GL program
func LoadProgram(vert, frag string) (uint32, error) {
	vertShader, err := compileShader(vert, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("could not compile vertex shader: %s", err)
	}
	defer gl.DeleteShader(vertShader)

	fragShader, err := compileShader(frag, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("could not compile fragment shader: %s", err)
	}
	defer gl.DeleteShader(fragShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	if status == gl.FALSE {
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(log))

		return 0, fmt.Errorf("cannot link program: %s", log)
	}

	return program, nil
}

func compileShader(src string, type_ uint32) (uint32, error) {
	shader := gl.CreateShader(type_)

	cSrc, free := gl.Strs(src)
	defer free()
	gl.ShaderSource(shader, 1, cSrc, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(shader, length, nil, gl.Str(log))

		return 0, fmt.Errorf("cannot compile shader: %s", log)
	}

	return shader, nil
}
