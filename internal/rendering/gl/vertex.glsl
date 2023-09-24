// using 300 es to be compatible with web
// maybe that will be changed in case special features are needed for desktop
#version 300 es
precision mediump float;
layout(location = 0) in vec3 aPos;
// color.w/a is layer
layout(location = 1) in vec4 color;
// layouts from channels
<attributes>
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
        (<xAxis> * aPos.x + <yAxis> * aPos.y + pos),
        // layer
        // TODO
        (<layer>*256 + color.w*255),
        1.0
    )
}