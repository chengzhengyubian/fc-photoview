import{i as G,r as c,S as F,j as l,a as v,b as C,d as g,E as D,u as L,g as _,G as P,F as E,L as $}from"./index.ed54fe83.js";import{b as N,c as k,d as S,t as R,e as A,u as Y,P as Z}from"./MediaGallery.6f345196.js";import{u as j,P as q}from"./MapboxMap.27d07c02.js";import{F as O,u as Q}from"./useURLParameters.1aca2057.js";import{D as U}from"./Dropdown.73d202b6.js";import"./Input.1cb7b307.js";import"./Table.941991f2.js";function B(t,n){switch(n.type){case"replaceTimelineGroups":{const e=z(n.timeline);return{...t,activeIndex:{album:-1,date:-1,media:-1},timelineGroups:e}}case"nextImage":{const{activeIndex:e,timelineGroups:a}=t;if(e.album==-1&&e.date==-1&&e.media==-1)return t;const r=a[e.date].albums,i=r[e.album].media;return e.media<i.length-1?{...t,activeIndex:{...t.activeIndex,media:e.media+1}}:e.album<r.length-1?{...t,activeIndex:{...t.activeIndex,album:e.album+1,media:0}}:e.date<a.length-1?{...t,activeIndex:{date:e.date+1,album:0,media:0}}:t}case"previousImage":{const{activeIndex:e}=t;if(e.album==-1&&e.date==-1&&e.media==-1)return t;if(e.media>0)return{...t,activeIndex:{...e,media:e.media-1}};if(e.album>0){const r=t.timelineGroups[e.date].albums[e.album-1].media;return{...t,activeIndex:{...e,album:e.album-1,media:r.length-1}}}if(e.date>0){const a=t.timelineGroups[e.date-1].albums,r=a[a.length-1].media;return{...t,activeIndex:{date:e.date-1,album:a.length-1,media:r.length-1}}}return t}case"selectImage":return{...t,activeIndex:n.index};case"openPresentMode":return{...t,presenting:!0,activeIndex:n.activeIndex};case"closePresentMode":return{...t,presenting:!1}}}const W=({mediaState:t,index:n})=>{const{date:e,album:a,media:r}=n;return t.timelineGroups[e].albums[a].media[r]},y=({mediaState:t})=>{if(!Object.values(t.activeIndex).reduce((n,e)=>e===-1||n,!1))return W({mediaState:t,index:t.activeIndex})};function z(t){const n=[];let e=[],a=null;const r=(i,o)=>i.replace(/\d{2}:\d{2}:\d{2}/,"00:00:00")==o.replace(/\d{2}:\d{2}:\d{2}/,"00:00:00");for(const i of t){if(a==null){a={id:i.album.id,title:i.album.title,media:[i]};continue}if(!r(a.media[0].date,i.date)){e.push(a),n.push({date:e[0].media[0].date.replace(/\d{2}:\d{2}:\d{2}/,"00:00:00"),albums:e}),e=[],a={id:i.album.id,title:i.album.title,media:[i]};continue}if(a.id!=i.album.id){e.push(a),a={id:i.album.id,title:i.album.title,media:[i]};continue}a.media.push(i)}return G(a)||(e.push(a),n.push({date:e[0].media[0].date.replace(/\d{2}:\d{2}:\d{2}/,"00:00:00"),albums:e})),n}const H=({dispatchMedia:t,activeIndex:n})=>{t({type:"openPresentMode",activeIndex:n}),history.pushState({presenting:!0,activeIndex:n},"")},V=({dateIndex:t,albumIndex:n,mediaState:e,dispatchMedia:a})=>{const{media:r,title:i,id:o}=e.timelineGroups[t].albums[n],[b]=N(),{updateSidebar:m}=c.exports.useContext(F),d=r.map((s,p)=>{var h;return l(k,{media:s,selectImage:()=>{a({type:"selectImage",index:{album:n,date:t,media:p}}),m(l(S,{media:s}))},clickPresent:()=>{H({dispatchMedia:a,activeIndex:{album:n,date:t,media:p}})},clickFavorite:()=>{R({media:s,markFavorite:b})},active:s.id===((h=y({mediaState:e}))==null?void 0:h.id)},s.id)});return v("div",{className:"mx-2",children:[l(C,{to:`/album/${o}`,className:"hover:underline",children:i}),v("div",{className:"flex flex-wrap items-center relative -mx-1 overflow-hidden",children:[d,l(A,{})]})]})},X={year:"numeric",month:"long",day:"numeric"},J=({groupIndex:t,mediaState:n,dispatchMedia:e})=>{const{i18n:a}=g(),r=n.timelineGroups[t],i=r.albums.map((m,d)=>l(V,{dateIndex:t,albumIndex:d,mediaState:n,dispatchMedia:e},`${r.date}_${m.id}`)),b=new Intl.DateTimeFormat(a.language,X).format(new Date(r.date));return v("div",{className:"mx-3 mb-2",children:[l("div",{className:"text-xl m-0 -mb-2",children:b}),l("div",{className:"flex flex-wrap -mx-2 my-0",children:i})]})},K=t=>c.exports.createElement("svg",{width:"13px",height:"15px",viewBox:"0 0 13 14",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...t},c.exports.createElement("g",{id:"Symbols",stroke:"none",strokeWidth:1,fill:"none",fillRule:"evenodd"},c.exports.createElement("path",{d:"M9.16666667,5.68434189e-14 C9.41979718,5.68434189e-14 9.62899397,0.188102588 9.66210226,0.432152962 L9.66666667,0.5 L9.666,1.333 L11.1666667,1.33333333 C12.1285626,1.33333333 12.9174396,2.07411552 12.9939226,3.01630483 L13,3.16666667 L13,12.5 C13,13.512522 12.1791887,14.3333333 11.1666667,14.3333333 L11.1666667,14.3333333 L1.83333333,14.3333333 C0.820811292,14.3333333 0,13.512522 0,12.5 L0,12.5 L0,3.16666667 C0,2.15414463 0.820811292,1.33333333 1.83333333,1.33333333 L1.83333333,1.33333333 L3.333,1.333 L3.33333333,0.5 C3.33333333,0.223857625 3.55719096,5.68434189e-14 3.83333333,5.68434189e-14 C4.08646384,5.68434189e-14 4.29566064,0.188102588 4.32876892,0.432152962 L4.33333333,0.5 L4.333,1.333 L8.666,1.333 L8.66666667,0.5 C8.66666667,0.223857625 8.89052429,5.68434189e-14 9.16666667,5.68434189e-14 Z M12,6.333 L1,6.333 L1,12.5 C1,12.9248344 1.31790432,13.2754183 1.72880177,13.3268405 L1.83333333,13.3333333 L11.1666667,13.3333333 C11.626904,13.3333333 12,12.9602373 12,12.5 L12,12.5 L12,6.333 Z M3.333,2.333 L1.83333333,2.33333333 C1.37309604,2.33333333 1,2.70642938 1,3.16666667 L1,3.16666667 L1,5.333 L12,5.333 L12,3.16666667 C12,2.74183224 11.6820957,2.39124835 11.2711982,2.33982618 L11.1666667,2.33333333 L9.666,2.333 L9.66666667,3.16666667 C9.66666667,3.44280904 9.44280904,3.66666667 9.16666667,3.66666667 C8.91353616,3.66666667 8.70433936,3.47856408 8.67123108,3.2345137 L8.66666667,3.16666667 L8.666,2.333 L4.333,2.333 L4.33333333,3.16666667 C4.33333333,3.44280904 4.10947571,3.66666667 3.83333333,3.66666667 C3.58020282,3.66666667 3.37100603,3.47856408 3.33789774,3.2345137 L3.33333333,3.16666667 L3.333,2.333 Z",fill:"currentColor",fillRule:"nonzero"}))),ee=D`
  query earliestMedia {
    myMedia(
      order: { order_by: "date_shot", order_direction: ASC }
      paginate: { limit: 1 }
    ) {
      id
      date
    }
  }
`,te=({filterDate:t,setFilterDate:n})=>{const{t:e}=g(),{data:a,loading:r}=L(ee);let i=[{value:"all",label:e("timeline_filter.date.dropdown_all","From today")}];if(a&&a.myMedia.length!=0){const o=a.myMedia[0].date,b=new Date(o),d=new Date().getFullYear(),s=b.getFullYear(),p=[];for(let u=d-1;u>=s;u--)p.push(u);const h=p.map(u=>({value:`${u}`,label:e("timeline_filter.date.dropdown_year","{{year}} and earlier",{year:u})}));i=[...i,...h]}return v("fieldset",{children:[v("legend",{id:"filter_group_date-label",className:"inline-block mb-1",children:[l(K,{className:"inline-block align-baseline mr-1","aria-hidden":"true"}),l("span",{children:e("timeline_filter.date.label","Date")})]}),l("div",{children:l(U,{"aria-labelledby":"filter_group_date-label",setSelected:o=>o=="all"?n(null):n(o),value:t||"all",items:i,disabled:r})})]})},ae=({onlyFavorites:t,setOnlyFavorites:n,filterDate:e,setFilterDate:a})=>v("div",{className:"flex items-end gap-4 flex-wrap mb-4",children:[l(te,{filterDate:e,setFilterDate:a}),l(O,{onlyFavorites:t,setOnlyFavorites:n})]}),ie=_`
  query myTimeline(
    $onlyFavorites: Boolean
    $limit: Int
    $offset: Int
    $fromDate: Time
  ) {
    myTimeline(
      onlyFavorites: $onlyFavorites
      fromDate: $fromDate
      paginate: { limit: $limit, offset: $offset }
    ) {
      id
      title
      type
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
      }
      favorite
      album {
        id
        title
      }
      date
    }
  }
`,ne=()=>{const{t}=g(),{getParam:n,setParam:e}=Q(),a=n("favorites")=="1",r=f=>e("favorites",f?"1":null),i=n("date"),o=f=>e("date",f),b=c.exports.useRef(!1),[m,d]=c.exports.useReducer(B,{presenting:!1,timelineGroups:[],activeIndex:{media:-1,album:-1,date:-1}}),{data:s,error:p,loading:h,refetch:u,fetchMore:I}=L(ie,{variables:{onlyFavorites:a,fromDate:i?`${parseInt(i)+1}-01-01T00:00:00Z`:void 0,offset:0,limit:200}}),{containerElem:w,finished:M}=Y({loading:h,fetchMore:I,data:s,getItems:f=>f.myTimeline});if(c.exports.useEffect(()=>{d({type:"replaceTimelineGroups",timeline:(s==null?void 0:s.myTimeline)||[]})},[s]),c.exports.useEffect(()=>{(async()=>(await P.resetStore(),await u({onlyFavorites:a,fromDate:i?`${parseInt(i)+1}-01-01T00:00:00Z`:void 0,offset:0,limit:200})))()},[i]),j({dispatchMedia:d,openPresentMode:f=>{d({type:"openPresentMode",activeIndex:f.state.activeIndex})}}),c.exports.useEffect(()=>{b.current=!1,u({onlyFavorites:a})},[a]),p)return l("div",{children:p.message});const T=m.timelineGroups.map((f,x)=>l(J,{groupIndex:x,mediaState:m,dispatchMedia:d},x));return v("div",{className:"overflow-x-hidden",children:[l(ae,{onlyFavorites:a,setOnlyFavorites:r,filterDate:i,setFilterDate:o}),l("div",{className:"-mx-3 flex flex-wrap",ref:w,children:T}),l(Z,{active:!M&&!h,text:t("general.loading.paginate.media","Loading more media")}),m.presenting&&l(q,{activeMedia:y({mediaState:m}),dispatchMedia:d})]})},ce=()=>{const{t}=g();return l(E,{children:l($,{title:t("photos_page.title","Timeline"),children:l(ne,{})})})};export{ce as default};
