#version 330

layout (location = 0) in vec3 vertPosition;
layout (location = 1) in vec2 vertTexCoord;
layout (location = 2) in uint TextureID;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform mat3 invMV;

out vec2 fragTexCoord;
out vec3 position;
out uint textureID;
out vec3 sunDirection;

void main() {
    sunDirection = normalize(invMV * vec3(-0.3, 0.75, 1.0));
    fragTexCoord = vertTexCoord;
    textureID = TextureID;
    vec4 pos = projection * view * model * vec4(vertPosition, 1); 
    gl_Position = pos;
    position = vertPosition;
}
