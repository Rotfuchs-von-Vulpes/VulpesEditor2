#version 330 core

uniform sampler2D art;
uniform float zoom;
uniform float aspect;

in vec2 TexCoords;

out vec4 frag_color;

void main () {
    vec2 borderThickness = vec2(aspect*zoom, zoom);

    if (TexCoords.x < borderThickness.x || TexCoords.x > (1.0 - borderThickness.x) || TexCoords.y < borderThickness.y || TexCoords.y > (1.0 - borderThickness.y)) {
        gl_FragColor = vec4(0.8, 0.284, 0.16, 1.0); // Render border in red
    } else {
        vec4 color = texture(art, TexCoords).rgba;
        if (color.a < 0.5) {
            gl_FragColor = vec4(0.5, 0.5, 0.5, 1.0);
        } else {
            gl_FragColor = color;
        }
    }
}