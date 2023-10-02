#version 300 es
precision mediump float;

out vec4 FragColor;

in vec4 vertexColor;
in flat uint layer;

void main() {
    FragColor = vertexColor;
    // maybe should be 1 higher but that would be bigger than uint
    gl_FragDepth = float(layer) / 4294967295.0;
}