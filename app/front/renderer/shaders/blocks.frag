#version 330

// uniform sampler2D tex;

in vec2 fragTexCoord;
in vec3 position;
flat in uint textureID;

out vec4 outputColor;

vec3 calculateLocalNormal(vec3 localPos) {
    vec3 localTangentX = dFdx(localPos);
    vec3 localTangentY = dFdy(localPos);
    
    return normalize(cross(localTangentX, localTangentY));
}

void main() {
    vec3 localNormal = calculateLocalNormal(position);
    
    outputColor = vec4(localNormal * 0.5 + vec3(0.5), 1.0);
}
