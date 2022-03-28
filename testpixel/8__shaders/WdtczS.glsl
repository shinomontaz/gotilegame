#version 330 core

float bounceHeight(float time, float period, float maxHeight) {
    float modulus = mod(time / period, 1.0);
    return modulus * (1.0 - modulus) * 4.0 * maxHeight;
}

in vec2 vTexCoords;
out vec4 fragColor;

// Pixel default uniforms
uniform vec4      uTexBounds;
uniform sampler2D uTexture;
uniform float uTime;

void main() {
    float ballRadius = 50.0;

    float bouncePeriod = 2.0;
    float maxBounceHeight = 0.75 * (uTexBounds.y - ballRadius);

    vec2 ballPosition = vec2(0.5 * uTexBounds.x, ballRadius + bounceHeight(uTime, bouncePeriod, maxBounceHeight));
//    vec2 ballPosition = vec2(0.5 * uTexBounds.x, ballRadius);

    vec3 ballColor = vec3(1.0, 0.0, 0.0);

    vec3 color = vec3(0.0, 0.0, 0.0);
    if (length(vTexCoords - ballPosition) <= ballRadius) {
        color = ballColor;
    }

    fragColor = vec4(color, 1.0);
}