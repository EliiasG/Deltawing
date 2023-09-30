// using 300 es to be compatible with web
// maybe that will be changed in case special features are needed for desktop
#version 300 es
precision mediump float;
layout(location = 0) in vec2 aPos;
// color.w/a/q is layer
layout(location = 1) in ivec4 aColor;
// layouts from channels
<attributes>

out vec4 vertexColor;
out uint layer;

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
        (<xAxis> * aPos.x + <yAxis> * aPos.y + <pos>),
        0, 1.0
    );

    // TODO maybe add modulate like godot
    vertexColor = vec4(aColor.rgb, 1);
    layer = <layer>*256 + aColor.a;
}