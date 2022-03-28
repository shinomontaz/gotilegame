#version 330 core
// fragCoord -> vTexCoords
// iResolution.xy -> uTexBounds.zw
// mainImage(out vec4 fragColor, in vec2 fragCoord) -> main() + definition of 
// in vec2 vTexCoords;
// out vec4 fragColor;

const float objectSize=0.04;
const float lampSize = 0.02;

vec2[] points = vec2[] (
        vec2(0.1, 0.7),
        vec2(0.2, 0.7),
        vec2(0.3, 0.7),
        vec2(0.4, 0.7),
        vec2(0.5, 0.7),
        vec2(0.6, 0.7),
        vec2(0.7, 0.7),
        vec2(0.8, 0.7),
        
        vec2(0.1, 0.3),
        vec2(0.2, 0.3),
        vec2(0.3, 0.3),
        vec2(0.4, 0.3),
        vec2(0.5, 0.3),
        vec2(0.6, 0.3),
        vec2(0.7, 0.3),
        vec2(0.8, 0.3)
);

float isObject(vec2 uv, float bias) {
    float dist = 1.0;
    for (int i = 0; i < 16; ++i ) {
        vec2 bpoint = vec2(points[i].x * bias, points[i].y);
        float d = (distance(bpoint, uv));
        if  (d < dist) dist =d;
    }
    return dist;
}

float hitsObject(vec2 a, vec2 n, vec2 p) {
    vec2 projection = (dot((a-p),n) * n);
    vec2 projectionNormal = (a-p)- (dot((a-p),n) * n);
    float l = length(projectionNormal);
    return l;
}

float isLamp(vec2 uv,vec2 posLamp) {
   	    return distance(uv, posLamp);
}

float isShadow(vec2 uv, vec2 posLamp, float bias) {
    int count = 0;
    float dist = 1.0;
    for (int i = 0; i< 16; i++ ) {
        vec2 p = points[i];
        p.x *= bias;
        
        vec2 normUvLamp = normalize(posLamp-uv);
        vec2 normPLamp = normalize(posLamp-p);
        
        
        if (length(normUvLamp+normPLamp) > 0.5) {
        
        	if ((distance(uv, posLamp) > distance(p, posLamp)))  {
                float d = hitsObject(uv, normalize(uv-posLamp), p);
     			if (d< dist )  dist = d ;
        	}
        }
        
    }
    return dist;
}

void mainImage( out vec4 fragColor, in vec2 fragCoord )
{
    vec2 uv = fragCoord.xy / iResolution.y;
    float bias = iResolution.x / iResolution.y;
    vec2 posLamp = vec2((0.5 * sin(iTime*0.2) + 0.5) * bias , 0.5);
    vec4 cShadow = vec4(0.0,0.0,0.0, 1.0);  
    vec4 cObject = vec4(0.0,1.0,0.0, 1.0);
    vec4 cFloor = vec4(0.0,0.0,1.0, 1.0);
    vec4 cLamp = vec4(1.0,1.0,1.0, 1.0);

    float tolerance = 0.01;
    
    float distObject = isObject(uv, bias);
    float distShadow = isShadow(uv, posLamp, bias);
    float distLamp  = isLamp(uv, posLamp);
    
    float t = 0.005;
    float shadow = 1.0;
    if ((distShadow > objectSize - t) && distShadow < objectSize ) {
    	 shadow = smoothstep(0.0, 1.0, (distShadow-(objectSize-t))/t);
    } else if (distShadow < objectSize) {
        shadow = 0.0;
    }
    
    cFloor = vec4(0.0,0.0,shadow,1.0);
    
    if ((distObject > objectSize -t) && distObject < objectSize) {
        float step = smoothstep(1.0, 0.0, (distObject-(objectSize-t))/t);
        vec4 m = mix(cFloor, cObject, step);
        fragColor = m;
    } else if ((distObject < objectSize -t)) {
        fragColor = cObject;
    } else if (distLamp < lampSize) {
        float tlamp = 0.015;
        if (distLamp < lampSize - tlamp) {
        	fragColor = vec4(1.0,1.0,1.0,1.0);
        } else {
            
            float step = smoothstep(1.0, 0.0, (distLamp-(lampSize-tlamp))/tlamp);
            vec4 m = mix(cFloor, cLamp, step);
            fragColor = m;
        }
    } else {
        fragColor = cFloor;
    }
}