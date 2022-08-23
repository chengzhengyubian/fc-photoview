import{s as l,a as h,j as i,H as k,g as b,h as v,r as M,d as y,u as P,L as g,I as w}from"./index.ed54fe83.js";import{P as R,m as _,a as I,u as E}from"./MapboxMap.27d07c02.js";var C="/assets/image-popup.ac743778.svg";const L=l.div`
  width: 56px;
  height: 68px;
  position: relative;
  margin-top: -54px;
  cursor: pointer;
`,O=l.img`
  position: absolute;
  width: 48px;
  height: 48px;
  top: 4px;
  left: 4px;
  border-radius: 2px;
  object-fit: cover;
`,S=l.img`
  width: 100%;
  height: 100%;
`,A=l.div`
  position: absolute;
  top: -10px;
  right: -10px;
  width: 24px;
  height: 24px;
  background-color: #00b3dc;
  border-radius: 50%;
  color: white;
  text-align: center;
  padding-top: 2px;
`,T=({marker:e,dispatchMarkerMedia:t})=>{const r=JSON.parse(e.thumbnail);return h(L,{onClick:()=>{t({type:"replacePresentMarker",marker:{cluster:!!e.cluster,id:e.cluster?e.cluster_id:e.media_id}})},children:[i(S,{src:C}),i(O,{src:r.url}),e.cluster&&i(A,{children:e.point_count_abbreviated})]})},x={};let m={};const j=e=>{const t=q(e);e.map.on("move",t),e.map.on("moveend",t),e.map.on("sourcedata",t),t()},q=({map:e,mapboxLibrary:t,dispatchMarkerMedia:r})=>()=>{const s={},a=e.querySourceFeatures("media");for(const n of a){const p=n.geometry.coordinates,o=n.properties;if(o==null){console.warn("WARN: geojson feature had no properties",n);continue}const d=o.cluster?`cluster_${o.cluster_id}`:`media_${o.media_id}`;let c=x[d];if(!c){const f=Q(o,{dispatchMarkerMedia:r});c=x[d]=new t.Marker({element:f}).setLngLat(p)}s[d]=c,m[d]||c.addTo(e)}for(const n in m)s[n]||m[n].remove();m=s};function Q(e,{dispatchMarkerMedia:t}){const r=document.createElement("div");return k.render(i(T,{marker:e,dispatchMarkerMedia:t}),r),r}const D=b`
  query placePageQueryMedia($mediaIDs: [ID!]!) {
    mediaList(ids: $mediaIDs) {
      id
      title
      blurhash
      thumbnail {
        url
        width
        height
      }
      highRes {
        url
        width
        height
      }
      videoWeb {
        url
        width
        height
      }
      type
    }
  }
`,G=(e,t)=>new Promise((r,s)=>{var u;const{cluster:a,id:n}=t;if(a)e.getSource("media").getClusterLeaves(n,1e3,0,(o,d)=>{if(o){s(o);return}const c=d.map(f=>f.properties);r(c)});else{const o=(u=e.querySourceFeatures("media").find(d=>{var c;return((c=d.properties)==null?void 0:c.media_id)==n}))==null?void 0:u.properties;if(o===void 0){s("ERROR: media is undefined");return}r([o])}}),H=({map:e,markerMediaState:t,dispatchMarkerMedia:r})=>{const[s,{data:a}]=v(D);return M.exports.useEffect(()=>{const n=t.presentMarker;if(n==null||e==null){r({type:"closePresentMode"});return}G(e,n).then(u=>{s({variables:{mediaIDs:u.map(p=>p.media_id)}})})},[t.presentMarker]),M.exports.useEffect(()=>{const n=(a==null?void 0:a.mediaList)||[];r({type:"replaceMedia",media:n}),n.length>0&&r({type:"openPresentMode",activeIndex:0})},[a]),t.presenting?i(R,{activeMedia:t.media[t.activeIndex],dispatchMedia:r,disableSaveCloseInHistory:!0}):null};function J(e,t){switch(t.type){case"replacePresentMarker":return e.presentMarker&&t.marker&&e.presentMarker.cluster===t.marker.cluster&&e.presentMarker.id===t.marker.id?{...e,presenting:!0}:{...e,presentMarker:t.marker};default:return _(e,t)}}const N=l.div`
  width: 100%;
  height: calc(100vh - 120px);
`,W=b`
  query mediaGeoJson {
    myMediaGeoJson
  }
`,B=()=>{const{t:e}=y(),{data:t}=P(W,{fetchPolicy:"cache-first"}),[r,s]=M.exports.useReducer(J,{presenting:!1,activeIndex:-1,media:[]}),{mapContainer:a,mapboxMap:n,mapboxToken:u}=I({configureMapbox:$({mapboxData:t,dispatchMarkerMedia:s})});return E({dispatchMedia:s,openPresentMode:p=>{s({type:"openPresentMode",activeIndex:p.state.activeIndex})}}),t&&u==null?h(g,{title:e("places_page.title","Places"),children:[i("h1",{children:"Mapbox token is not set"}),h("p",{children:["To use map related features a mapbox token is needed.",i("br",{})," A mapbox token can be created for free at"," ",i("a",{href:"https://account.mapbox.com/access-tokens/",children:"mapbox.com"}),"."]}),i("p",{children:"Make sure the access token is added as the MAPBOX_TOKEN environment variable."})]}):h(g,{title:"Places",children:[i(w,{}),i(N,{children:a}),i(H,{map:n,markerMediaState:r,dispatchMarkerMedia:s})]})},$=({mapboxData:e,dispatchMarkerMedia:t})=>(r,s)=>{r.addControl(new s.NavigationControl),r.on("load",()=>{if(r==null){console.error("ERROR: map is null");return}r.addSource("media",{type:"geojson",data:e==null?void 0:e.myMediaGeoJson,cluster:!0,clusterRadius:50,clusterProperties:{thumbnail:["coalesce",["get","thumbnail"],!1]}}),r.addLayer({id:"media-points",type:"circle",source:"media",filter:["!",!0]}),j({map:r,mapboxLibrary:s,dispatchMarkerMedia:t})})};export{B as default};
