#define PI 3.14159265359

// https://www.shadertoy.com/view/4dSXR1

// fragCoord -> vTexCoords
// iResolution.xy -> uTexBounds.zw

const vec2 constantList = vec2(1.0, 0.0);

float antiAlias(float x) {return (x-(1.0-2.0/iResolution.y))*(iResolution.y/2.0);}

float asrat;
struct Obj{    int type;    vec3 col;    vec2 pos;    float emis;}T_Obj;
struct Hit{    float d,t;    vec2 pos;    Obj obj;}T_Hit;
struct Ray{    vec2 ro,rd;    float tmin,tmax;}T_Ray;
const Obj NoObj=Obj(-1,vec3(0.),vec2(0.),0.);
const Hit NoHit=Hit(400.,400.,vec2(0.),NoObj);

float sdCircle(vec2 p){    return distance(p,vec2(0.));}
float sdBox(vec2 p, vec2 size, float radius)
{
    vec2 d = abs(p) - size-vec2(radius);
    return clamp(min(max(d.x, d.y), 0.0) + length(max(d, 0.0)) - radius,0.,1.);
}

Hit renderSdf(vec2 p,inout Obj o)
{
    Hit hit=NoHit;
    if(o.type==1) hit= Hit(sdCircle(p-o.pos),0.,p,o);
    if(o.type==2) hit= Hit(sdBox(p-o.pos,vec2(0.02),0.),090.1,p,o);
    return hit;
}
Hit sdU(in Hit a,vec2 p,in Obj d)
{
    Hit hit=renderSdf(p,d);
    if (a.d<hit.d) 
        return (a); else return (hit);
}

Hit map(vec2 p, in Obj oo)
{
	//----8<----- SOLIDS ----8<-----//

    const vec3 solidCol=vec3(0.95,0.57,0.);

    Hit hit=NoHit;
    hit=sdU(hit,p,Obj(1,solidCol,vec2(0.5*asrat,0.5),0.));
    hit=sdU(hit,p,Obj(1,solidCol,vec2(0.5*asrat,0.8),0.));
    hit=sdU(hit,p,Obj(2,solidCol,vec2(0.53*asrat,0.2),0.));
    hit=sdU(hit,p,Obj(2,solidCol,vec2(0.53*asrat,0.3),0.));
    hit=sdU(hit,p,Obj(2,solidCol,vec2(0.47*asrat,0.2),0.));
    hit=sdU(hit,p,Obj(2,solidCol,vec2(0.47*asrat,0.3),0.));
    return hit;
}

vec3 render(vec3 incol,inout Hit hit)
{
    float aares=clamp(antiAlias(smoothstep(((1.-hit.d)), .977, 0.)),0.,1.);
    if(aares==0.)
    {
        hit=NoHit;    
    }
    return mix(incol,hit.obj.col,aares);
}
Hit castRay(in Ray ray){
    float precis =.023;
    float t = ray.tmin;
    Hit hit=NoHit;
    vec2 p;
    for( int i=0; i<200; i++ )
    {
        p=ray.ro+ray.rd*t;
        hit = map( p,NoObj );
        if( hit.d<precis || t>ray.tmax ) break;
        t += hit.d;
    }
    hit.t=t;

    if( t>ray.tmax ) return NoHit;
    return hit;

}
float getFloor(vec2 uv){    return mod( floor(32.0*(uv.x/asrat)) + floor(18.0*(uv.y)), 2.0);}

vec3 getLight(vec2 uv, Obj light)
{
    float rot=atan(uv.y-light.pos.y, uv.x-light.pos.x);
    Hit hit=castRay(Ray(light.pos.xy,vec2(sin(rot),cos(rot)).yx,.0001,distance(light.pos,uv)));

    render(vec3(0.),hit);

    float lamnt=10.5*distance(uv,light.pos);
    lamnt=pow((1.2/(lamnt*lamnt)),.5)*0.39*light.emis;

    vec3 lcol=vec3(0.01);        
    if(hit==NoHit)
    {
        lcol+=light.col*lamnt;
    }

    lcol-=lamnt*light.col*0.5*getFloor(uv);


    return lcol;
}

void main()
{
    vec2 uv = vTexCoords.xy / iResolution.xy;
    vec2 um = iMouse.xy / iResolution.xy;
    vec4 col = vec4(0.0);

    asrat=iResolution.x/iResolution.y;

    if(iMouse.z<0.5)
    {
        um=0.5+0.4*vec2(sin(iTime),cos(iTime));
    }

    uv.x*=asrat;um.x*=asrat;

    Hit hit=map(uv,NoObj);
    col.rgb +=render(vec3(0.),hit);

    if(hit==NoHit)
    {
        //----8<----- LIGHTS ----8<-----//

        col.rgb+=getLight(uv, Obj(1,vec3(0.5,1.,.5),um,1.));
        col.rgb+=getLight(uv, Obj(1,vec3(1.,.5,.5),vec2(asrat-um.x,1.-um.y),2.));
    }

    fragColor = col.xyzw * constantList.xxxy + constantList.yyyx;
}