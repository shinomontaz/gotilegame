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
uniform vec2 uLight;

// Utilities
#define drawSDF(dist, col) color = mix(color, col, smoothstep(1.0, 0.0, dist))

float fillMask(float dist)
{
	return clamp(-dist, 0.0, 1.0);
}

float circleDist(vec2 p, float radius)
{
	return length(p) - radius;
}

vec3 drawLight(vec2 p, vec2 pos, vec3 color, float range, float radius)
{
	// distance to light
	float ld = length(p - pos);
	
	// out of range
	if (ld > range) return vec3(0.0);
	
	// shadow and falloff
	float fall = (range - ld)/range;
	fall *= fall;
	float source = fillMask(circleDist(p - pos, radius));
	return (fall + source) * color;
}


// fragCoord -> vTexCoords
// iResolution.xy -> uTexBounds.zw
// mainImage(out vec4 fragColor, in vec2 fragCoord) -> main() + definition of 
// in vec2 vTexCoords;
// out vec4 fragColor;

void main() {
    vec2 uv = vTexCoords.xy / uTexBounds.zw;
    vec4 pixelColor = texture(uTexture, uv);

    vec2 toLight = uv - uLight/uTexBounds.zw;
    vec3 color = pixelColor.rgb / (dot(toLight, toLight) + pixelColor.rgb);

    vec2 circle_pos = uLight + uTexBounds.zw/2;

    color += drawLight(vTexCoords.xy, circle_pos, vec3(1.0, 0.75, 0.5), 100.0, 10.0);

    fragColor = vec4(color, 1.0);
}