#version 420
#extension GL_ARB_explicit_uniform_location : enable

layout(location = 0) in vec3 position;
layout(location = 1) in vec2 texCoord;
layout(location = 2) in vec3 normal;

layout(location = 3) uniform mat4 projection;
layout(location = 4) uniform mat4 view;
layout(location = 5) uniform mat4 model;

layout(location = 0) out vec2 TexCoord;

void main() {
    gl_Position = projection * view * model * vec4(position, 1.0);
    TexCoord = texCoord;
}
