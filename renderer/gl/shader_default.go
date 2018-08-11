package gl

const shaderDefaultVert = `
#version 430 core

layout (location = 0) in vec2 inPos;
layout (location = 1) in vec2 inTexCoord;

uniform mat4 uProjection;

out vec2 aTexCoord;

void main() {
    gl_Position = uProjection * vec4(inPos, 0.0, 1.0);
    aTexCoord = inTexCoord;
	aTexCoord.y = 1 - aTexCoord.y;
}
`

const shaderDefaultFrag = `
#version 430 core

out vec4 outFragColor;

in vec2 aTexCoord;

uniform sampler2D uTexture;

void main() {
    outFragColor = texture(uTexture, aTexCoord);
}
`
