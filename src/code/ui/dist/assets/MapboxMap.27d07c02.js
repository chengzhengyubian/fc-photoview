import{r as p,j as e,s as n,o as w,a as s,P as C,n as E,M as h,p as k,W as z,g as y,u as M,q as A,_ as V}from"./index.ed54fe83.js";function Y(t,a){switch(a.type){case"nextImage":return{...t,activeIndex:(t.activeIndex+1)%t.media.length};case"previousImage":return t.activeIndex<=0?{...t,activeIndex:t.media.length-1}:{...t,activeIndex:t.activeIndex-1};case"openPresentMode":return{...t,presenting:!0,activeIndex:a.activeIndex};case"closePresentMode":return{...t,presenting:!1};case"selectImage":return{...t,activeIndex:Math.max(0,Math.min(t.media.length-1,a.index))};case"replaceMedia":return{...t,media:a.media,activeIndex:-1,presenting:!1}}}const $=({dispatchMedia:t,openPresentMode:a})=>{p.exports.useEffect(()=>{const i=l=>{l.state.presenting===!0?a(l):t({type:"closePresentMode"})};return window.addEventListener("popstate",i),history.replaceState({presenting:!1},""),()=>{window.removeEventListener("popstate",i)}},[])},q=({dispatchMedia:t,activeIndex:a})=>{t({type:"openPresentMode",activeIndex:a}),history.pushState({presenting:!0,activeIndex:a},"")},d=({dispatchMedia:t})=>{t({type:"closePresentMode"}),history.back()};function B(t){return e("svg",{viewBox:"0 0 36 36",fillRule:"evenodd",clipRule:"evenodd",strokeLinecap:"round",strokeLinejoin:"round",strokeMiterlimit:1.5,width:"1em",height:"1em",...t,children:e("g",{fill:"none",stroke:"#000",strokeWidth:3,children:e("path",{d:"M2 2l32 32M2 34L34 2"})})})}function H(t){return e("svg",{viewBox:"0 0 28 52",fillRule:"evenodd",clipRule:"evenodd",strokeLinecap:"round",strokeLinejoin:"round",strokeMiterlimit:1.5,width:"1em",height:"1em",...t,children:e("path",{d:"M2 2l24 24L2 50",fill:"none",stroke:"#000",strokeWidth:3})})}function L(t){return e("svg",{viewBox:"0 0 28 52",fillRule:"evenodd",clipRule:"evenodd",strokeLinecap:"round",strokeLinejoin:"round",strokeMiterlimit:1.5,width:"1em",height:"1em",...t,children:e("path",{d:"M26 2L2 26l24 24",fill:"none",stroke:"#000",strokeWidth:3})})}const S=n.div`
  width: 100%;
  height: 100%;
  position: relative;
`,u=n.button`
  width: 64px;
  height: 64px;
  background: none;
  border: none;
  outline: none;
  cursor: pointer;
  position: absolute;

  & svg {
    width: 32px;
    height: 32px;
    overflow: visible !important;
  }

  & svg path {
    stroke: rgba(255, 255, 255, 0.5);
    transition-property: stroke, filter;
    transition-duration: 140ms;
  }

  &:hover svg path {
    stroke: rgba(255, 255, 255, 1);
    filter: drop-shadow(0px 0px 2px rgba(0, 0, 0, 0.6));
  }

  &.hide svg path {
    stroke: rgba(255, 255, 255, 0);
    transition: stroke 300ms;
  }
`,P=n(u)`
  left: 28px;
  top: 28px;
`,b=n(u)`
  height: 80%;
  width: 20%;
  top: 10%;

  ${({align:t})=>t=="left"?"left: 0;":null}
  ${({align:t})=>t=="right"?"right: 0;":null}

  & svg {
    margin: auto;
    width: 48px;
    height: 64px;
  }
`,I=({children:t,dispatchMedia:a,disableSaveCloseInHistory:i})=>{const[l,c]=p.exports.useState(!0),r=p.exports.useRef(null);return p.exports.useEffect(()=>(r.current=w(()=>{c(o=>!o)},2e3,!0),()=>{var o;(o=r.current)==null||o.cancel()}),[]),s(S,{"data-testid":"present-overlay",onMouseMove:()=>{r.current&&r.current()},children:[t,e(b,{"aria-label":"Previous image",className:l?"hide":void 0,align:"left",onClick:()=>a({type:"previousImage"}),children:e(L,{})}),e(b,{"aria-label":"Next image",className:l?"hide":void 0,align:"right",onClick:()=>a({type:"nextImage"}),children:e(H,{})}),e(P,{"aria-label":"Exit presentation mode",className:l?"hide":void 0,onClick:()=>{i===!0?a({type:"closePresentMode"}):d({dispatchMedia:a})},children:e(B,{})})]})},x=n(C)`
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  object-fit: contain;
  object-position: center;
`,R=n(E)`
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
`,T=({media:t,imageLoaded:a,...i})=>{var l,c;switch(t.type){case h.Photo:return s("div",{...i,children:[e(x,{src:(l=t.thumbnail)==null?void 0:l.url,"data-testid":"present-img-thumbnail"},`${t.id}-thumb`),e(x,{style:{display:"none"},src:(c=t.highRes)==null?void 0:c.url,"data-testid":"present-img-highres",onLoad:r=>{const o=r.target;o.style.display="initial",a&&a()}},`${t.id}-highres`)]});case h.Video:return e(R,{media:t,"data-testid":"present-video"})}k(t.type)},_=n.div`
  position: fixed;
  width: 100vw;
  height: 100vh;
  background-color: black;
  color: white;
  top: 0;
  left: 0;
  z-index: 100;
`,j=z`
  * {
    overflow: hidden !important;
  }
`,G=({className:t,imageLoaded:a,activeMedia:i,dispatchMedia:l,disableSaveCloseInHistory:c})=>(p.exports.useEffect(()=>{const r=o=>{o.key=="ArrowRight"&&(o.stopPropagation(),l({type:"nextImage"})),o.key=="ArrowLeft"&&(o.stopPropagation(),l({type:"previousImage"})),o.key=="Escape"&&(o.stopPropagation(),c===!0?l({type:"closePresentMode"}):d({dispatchMedia:l}))};return document.addEventListener("keydown",r),function(){document.removeEventListener("keydown",r)}}),s(_,{className:t,children:[e(j,{}),e(I,{dispatchMedia:l,disableSaveCloseInHistory:!0,children:e(T,{media:i,imageLoaded:a})})]}));const N=y`
  query mapboxToken {
    mapboxToken
    myMediaGeoJson
  }
`,O=n.div`
  width: 100%;
  height: 100%;
`,D=({configureMapbox:t,mapboxOptions:a=void 0})=>{var m;const[i,l]=p.exports.useState(),c=p.exports.useRef(null),r=p.exports.useRef(null),{data:o}=M(N,{fetchPolicy:"cache-first"});return p.exports.useEffect(()=>{async function g(){const v=(await V(()=>import("./mapbox-gl.f8123bca.js").then(function(f){return f.m}),["assets/mapbox-gl.f8123bca.js","assets/index.ed54fe83.js","assets/index.ccbdb765.css"])).default;l(v)}g()},[]),p.exports.useEffect(()=>{var g;i==null||c.current==null||o==null||r.current!=null||(o.mapboxToken&&(i.accessToken=o.mapboxToken),r.current=new i.Map({container:c.current,style:A()?"mapbox://styles/mapbox/dark-v10":"mapbox://styles/mapbox/streets-v11",...a}),t(r.current,i),(g=r.current)==null||g.resize())},[c,i,o]),(m=r.current)==null||m.resize(),{mapContainer:e(O,{ref:c}),mapboxMap:r.current,mapboxLibrary:i,mapboxToken:(o==null?void 0:o.mapboxToken)||null}};export{G as P,D as a,Y as m,q as o,$ as u};
