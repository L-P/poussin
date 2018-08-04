#version 430 core

out vec4 outFragColor;

in vec2 aTexCoord;

uniform sampler2D uTexture;

void main() {
    outFragColor = texture(uTexture, aTexCoord);
}
