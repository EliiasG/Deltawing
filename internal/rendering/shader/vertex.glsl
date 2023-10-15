// using 300 es to be compatible with web
// maybe that will be changed in case special features are needed for desktop
<version>
layout(location = 0) in vec2 aPos;
// color.w/a/q is layer
layout(location = 1) in ivec4 aColor;
// layouts from channels
<attributes>

out vec4 vertexColor;
flat out uint layer;

uniform ivec2 screenSize;
// uniforms from channels
<uniforms>

// functions from procedure
<functions>

void main() {
    // variables from channels
    <variables>

    // function calls from channels
    <calls>

    gl_Position = vec4(
        // position
        (<xAxis> * aPos.x + <yAxis> * -aPos.y + <pos>) / vec2(screenSize) * vec2(2, -2) + vec2(-1, 1),
        0, 1.0
    );

    vertexColor = vec4(vec3(<color>)/256.0, 1);
    layer = <layer>*uint(256) + uint(aColor.a+1);
}