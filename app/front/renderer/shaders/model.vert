#version 330

layout (location = 0) in vec3 vert;
layout (location = 1) in vec2 vertTexCoord;
layout (location = 2) in uint boneId;

const int MAX_BONES = 64;
uniform mat4 gBones[MAX_BONES];
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    mat4 boneTransform = gBones[boneId] * model;
    gl_Position = projection * view * boneTransform * vec4(vert, 1);
}
