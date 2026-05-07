#version 330 core

uniform sampler2D art;
uniform vec2 thickness;
uniform vec2 size;

in vec2 TexCoords;

out vec4 frag_color;

void main () {
    if (TexCoords.x < thickness.x || TexCoords.x > (1.0 - thickness.x) || TexCoords.y < thickness.y || TexCoords.y > (1.0 - thickness.y)) {
        gl_FragColor = vec4(0.8, 0.284, 0.16, 1.0);
    } else {
        vec4 color = texture(art, TexCoords).rgba;
        if (color.a < 0.5) {
            gl_FragColor = vec4(0.5, 0.5, 0.5, 1.0);
        } else {
            gl_FragColor = color;
        }
    }
}