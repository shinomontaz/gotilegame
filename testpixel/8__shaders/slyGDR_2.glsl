#version 330 core

// Pixel default uniforms
uniform vec4      uTexBounds;
uniform sampler2D uTexture;

// Our custom uniforms
uniform float uTime;
uniform vec2 uLight;
uniform mat2 uObject;

in vec2 vTexCoords;
out vec4 fragColor;

// Utilities
#define drawSDF(dist, col) color = mix(color, col, smoothstep(1.0, 0.0, dist))

float sdDisc(in vec2 p, in float r) {
    return length(p) - r;
}

float sdBox(in vec2 p, in vec2 b) {
    vec2 q = abs(p) - b;
    return length(max(q, 0.0)) + min(0.0, max(q.x, q.y));
}

float sdBox2(in vec2 uv, in vec2 tl, in vec2 br) {
    vec2 d = max(tl-uv, uv-br);
    return length(max(vec2(0.0), d)) + min(0.0, max(d.x, d.y));
}

void main( )
{
    vec2 uv = vTexCoords.xy;
    vec2 uv2 = vTexCoords / uTexBounds.zw;

    vec2 circle_pos = uLight + uTexBounds.zw/2;

   
    vec4 pixelColor = texture(uTexture, uv2);
    float radius = 50;

    vec2 tl = vec2(uObject[0][0], uObject[0][1]) + uTexBounds.zw/2;
    vec2 br = vec2(uObject[1][1], uObject[1][1]) + uTexBounds.zw/2;
    float box = sdBox2(uv, tl, br);
    float circle = sdDisc(uv - circle_pos, radius);

    vec3 color = pixelColor.rgb;
    drawSDF(circle, vec3(1.0, 0.8, 0.0));
    drawSDF(box, vec3(1.0, 0.0, 0.0));

//    float d = distance(uv, circle_pos) - radius;
//    drawSDF(d, vec3(1.0, 0.8, 0.0));

//    if (d < 0) {
//        pixelColor.r = 0;
//        pixelColor.g = 0;
//        pixelColor.b = 1;
//    }

    fragColor = vec4(color, 1.0);
}