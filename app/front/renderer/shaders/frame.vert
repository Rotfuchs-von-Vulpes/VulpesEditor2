#version 330

layout (location = 0) in vec3 vertPosition;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

void main() {
    vec4 pos = projection * view * model * vec4(vertPosition, 1); 
    gl_Position = pos;
}
