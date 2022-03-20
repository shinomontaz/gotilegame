#version 330 core

// https://www.shadertoy.com/view/7dlXWM

// Whether or not shadows can hide objects
//#define OBSTRUCT

in vec2 vTexCoords;
out vec4 fragColor;

// Pixel default uniforms
uniform vec4      uTexBounds;
uniform sampler2D uTexture;

// Our custom uniforms
uniform float uTime;
uniform vec2 uLight;
uniform mat2 uObject;

// Utilities
#define drawSDF(dist, col) color = mix(color, col, smoothstep(unit, 0.0, dist))

struct ShadowVol2D {
    vec2 ap;
    vec2 ad;
    vec2 bp;
    vec2 bd;
};

// Shadow volumes
ShadowVol2D shadowVolBox(in vec2 l, in vec2 b) {
    vec2 s = vec2(l.x < 0.0 ? -1.0 : 1.0, l.y < 0.0 ? -1.0 : 1.0);
    vec2 c1 = vec2(b.x * sign(b.y - abs(l.y)), b.y) * s;
    vec2 c2 = vec2(b.x, b.y * sign(b.x - abs(l.x))) * s;
    return ShadowVol2D(c1, normalize(c1 - l), c2, normalize(c2 - l));
}

float sdBox(in vec2 p, in vec2 b) {
    vec2 q = abs(p) - b;
    return length(max(q, 0.0)) + min(0.0, max(q.x, q.y));
}

float sdShadowVolume2D(in vec2 p, in vec2 ap, in vec2 ad, in vec2 bp, in vec2 bd) {
    vec2 pa = p - ap, pb = p - bp, ba = bp - ap;

    vec2 b = pa - ba * clamp(dot(pa, ba) / dot(ba, ba), 0.0, 1.0);
    vec2 e1 = pa - ad * max(0.0, dot(pa, ad) / dot(ad, ad));
    vec2 e2 = pb - bd * max(0.0, dot(pb, bd) / dot(bd, bd));

    vec2 bap = vec2(-ba.y, ba.x), h = 0.5 * (ad + bd);
    float s = sign(max(dot(pa, vec2(-ad.y, ad.x)) * dot(pb, vec2(-bd.y, bd.x)), dot(pa, bap) * sign(dot(bap, -h))));
    return sqrt(min(dot(b, b), min(dot(e1, e1), dot(e2, e2)))) * s;
}

// fragCoord -> vTexCoords
// iResolution.xy -> uTexBounds.zw
// mainImage(out vec4 fragColor, in vec2 fragCoord) -> main() + definition of 
// in vec2 vTexCoords;
// out vec4 fragColor;

void main() {
    vec2 t = vTexCoords / uTexBounds.zw;
    vec3 col = texture(uTexture, t).rgb;

    vec2 center = 0.5 * uTexBounds.zw;
    vec2 uv = (vTexCoords) / uTexBounds.zw;
//    float unit = 8.0 / uTexBounds.y;
    float unit = 1.0;

    // Inverse square (kinda)
    vec2 toLight = uv - uLight.xy;
    vec3 color = vec3(1.0 / (1.0 + dot(toLight, toLight)));

    // Shapes and shadow volumes

    vec2 bp = vec2(0.0, -1.0);
    vec2 bb = vec2(0.5, 0.2);

//    vec2 bp = vec2(uObject[0][0], uObject[0][1]);
//    vec2 bb = vec2(uObject[1][1], uObject[1][1]);

    ShadowVol2D boxShadow = shadowVolBox(uLight.xy - bp, bb); // Object space
    boxShadow.ap += bp, boxShadow.bp += bp; // Back to world space
    float boxShadowVol = sdShadowVolume2D(uv, boxShadow.ap, boxShadow.ad, boxShadow.bp, boxShadow.bd); // Shadow volume distance
    float box = sdBox(uv - bp, bb); // Box distance

    if ( sdBox(uLight.xy - bp, bb) > 0.0 ) {
        drawSDF(boxShadowVol, vec3(0.0)); // Draw shadow volumes
    }
    else { // Light is inside an object
        color = vec3(0.0);
    }

    drawSDF(box, vec3(0.8, 0.0, 0.0));
    fragColor = vec4(col *color, 1.0);
}