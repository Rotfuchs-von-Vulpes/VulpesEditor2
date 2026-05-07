#version 330 core

layout (location = 0) in vec2 aPos;
layout (location = 1) in vec2 aTexCoords;

out vec2 TexCoords;

// uniform float zoom;
uniform float aspect;
uniform vec2 move;
uniform mat2 model;

void main() {
   vec2 pos = model * (aPos + move);
   gl_Position = vec4(aspect*pos.x, pos.y, 0.0, 1.0);
   TexCoords = aTexCoords;
}