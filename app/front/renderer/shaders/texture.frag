#version 330 core

uniform sampler2D art;
uniform float zoom;

in vec2 TexCoords;

out vec4 frag_color;

void main () {
    float borderThickness = zoom;

    if (TexCoords.x < borderThickness || TexCoords.x > (1.0 - borderThickness) || TexCoords.y < borderThickness || TexCoords.y > (1.0 - borderThickness)) {
        gl_FragColor = vec4(0.8, 0.284, 0.16, 1.0); // Render border in red
    } else {
        vec4 color = texture(art, TexCoords).rgba;
        if (color.a < 0.5) discard; 
        gl_FragColor = color;
    }
}