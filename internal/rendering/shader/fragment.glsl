#version 300 es
precision highp float;

out vec4 FragColor;

in vec4 vertexColor;
flat in uint layer;

void main() {
    FragColor = vertexColor;
    // maybe should be 1 higher but that would be bigger than int
    gl_FragDepth = float(layer) / 16777216.0;// * 5.96046448e-8;
}