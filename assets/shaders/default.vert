#version 430 core

layout (location = 0) in vec2 inPos;
layout (location = 1) in vec2 inTexCoord;

uniform mat4 uProjection;
uniform mat4 uView;
uniform mat4 uModel;

out vec2 aTexCoord;

void main() {
    gl_Position = uProjection * uView * uModel * vec4(inPos, 0.0, 1.0);
    aTexCoord = inTexCoord;
}
