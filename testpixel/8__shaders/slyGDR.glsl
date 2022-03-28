#version 330 core

// Pixel default uniforms
uniform vec4      uTexBounds;
uniform sampler2D uTexture;

// Our custom uniforms
uniform float uTime;

// fragCoord -> vTexCoords
// iResolution.xy -> uTexBounds.zw
// mainImage(out vec4 fragColor, in vec2 fragCoord) -> main() + definition of 
in vec2 vTexCoords;
out vec4 fragColor;

void main( )
{
    vec2 uv = vTexCoords.xy;
    vec2 circle_pos = uTexBounds.zw/2.0;
    //vec2 circle_pos = uv/2.0;
    
   
    float radius =  0.2 * uv.y;
    float d = length(uv - circle_pos) - radius;
    
    vec4 col = vec4 (0.2,0.1,d,d);

    fragColor = vec4(col);
}