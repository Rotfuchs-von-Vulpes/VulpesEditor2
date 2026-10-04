#version 330

uniform sampler2DArray tex;

in vec2 fragTexCoord;
in vec3 position;
in vec3 sunDirection;
flat in uint textureID;

out vec4 outputColor;

vec3 calculateLocalNormal(vec3 localPos) {
    vec3 localTangentX = dFdx(localPos);
    vec3 localTangentY = dFdy(localPos);
    
    return normalize(cross(localTangentX, localTangentY));
}

float luminance(vec3 color) {
    return dot(color, vec3(0.2125f, 0.7153f, 0.0721f));
}

const vec3 sunColor = vec3(0.5f, 0.45f, 0.45f);
const vec3 skyColor = vec3(0.8f, 0.8f, 0.9f);
const vec3 _Ambient = vec3(0.02f, 0.04f, 0.08f);

vec3 CalculateLighting(vec3 albedo, vec3 normal, vec3 fragCoords) {
    vec3 ndotl = sunColor * clamp(dot(normal, sunDirection), 0.0f, 1.0f);
    ndotl *= 1.3;
    ndotl *= (luminance(skyColor) + 0.01f);

    vec3 lighting = ndotl + skyColor + _Ambient;

    vec3 diffuse = albedo.rgb;
    diffuse *= lighting;

    return diffuse;
}

void main() {
    vec3 Normal = calculateLocalNormal(position);

    vec4 albedo = texture(tex, vec3(fragTexCoord, textureID));
    vec3 diffuse = CalculateLighting(albedo.rgb, Normal, gl_FragCoord.xyz);
    
    outputColor = vec4(diffuse, 1.0);
}
