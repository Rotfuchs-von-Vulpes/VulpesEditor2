#version 330

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = vec4(fragTexCoord.x, 0.0, fragTexCoord.y, 1.0);
}
