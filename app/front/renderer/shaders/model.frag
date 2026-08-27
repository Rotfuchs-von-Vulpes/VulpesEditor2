#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = vec4(texture(tex, fragTexCoord).rgb, 1.0);
}
