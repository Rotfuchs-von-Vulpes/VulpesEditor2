#version 330 core

uniform sampler2D art;
uniform sampler2D preview;
uniform int outline;
uniform vec2 size;

in vec2 TexCoords;

out vec4 frag_color;

void main () {
    if (outline == 1) {
        gl_FragColor = vec4(0.8, 0.284, 0.16, 1.0);
    } else {
        vec4 background;
        vec4 color;
        vec4 textureColor = texture(art, TexCoords).rgba;
        vec4 previewColor = texture(preview, TexCoords).rgba;
        if (previewColor.a > 0) {
            color = previewColor;
        } else {
            color = textureColor;
        }
        float w = size.x * TexCoords.x * 0.0625;
        float h = size.y * TexCoords.y * 0.0625;
        if (fract(w) > 0.5 ^^ fract(h) > 0.5) {
            background = vec4(0.25, 0.25, 0.25, 1.0);
        } else {
            background = vec4(0.375, 0.375, 0.375, 1.0);
        }
        gl_FragColor = vec4(mix(background.rgb, color.rgb, color.a), 1.0);
    }
}