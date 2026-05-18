#version 330 core

uniform sampler2D art;
uniform sampler2D preview;
uniform int outline;

in vec2 TexCoords;

out vec4 frag_color;

void main () {
    if (outline == 1) {
        gl_FragColor = vec4(0.8, 0.284, 0.16, 1.0);
    } else {
        vec4 color1 = texture(preview, TexCoords).rgba;
        if (color1.a < 0.5) {
            vec4 color2 = texture(art, TexCoords).rgba;
            if (color2.a < 0.5) {
                gl_FragColor = vec4(0.5, 0.5, 0.5, 1.0);
            } else {
                gl_FragColor = color2;
            }
        } else {
            gl_FragColor = color1;
        }
    }
}