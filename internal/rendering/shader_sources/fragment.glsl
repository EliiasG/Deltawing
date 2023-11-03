<version>

out vec4 FragColor;

in vec4 vertexColor;
flat in uint layer;

void main() {
    FragColor = vertexColor;
    // maybe should be 1 higher but that would be bigger than int
    // 1/(2^24-1)
    gl_FragDepth = float(layer) * 5.96046448e-8;
}