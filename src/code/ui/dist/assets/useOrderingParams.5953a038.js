import{g as y,R as M,r as f,j as t,a as v}from"./index.ed54fe83.js";import{M as x,A as h,a as G}from"./MediaGallery.6f345196.js";import{A as u}from"./AlbumBoxes.33e38eeb.js";import{A as R}from"./useURLParameters.1aca2057.js";import{m as _,u as g}from"./MapboxMap.27d07c02.js";const k=y`
  ${x}

  fragment AlbumGalleryFields on Album {
    id
    title
    subAlbums(order: { order_by: "title", order_direction: $orderDirection }) {
      id
      title
      thumbnail {
        id
        thumbnail {
          url
        }
      }
    }
    media(
      paginate: { limit: $limit, offset: $offset }
      order: { order_by: $mediaOrderBy, order_direction: $orderDirection }
      onlyFavorites: $onlyFavorites
    ) {
      ...MediaGalleryFields
    }
  }
`,B=M.forwardRef(({album:e,loading:i=!1,customAlbumLink:a,showFilter:p=!1,setOnlyFavorites:n,setOrdering:l,ordering:o,onlyFavorites:s=!1},r)=>{const[A,d]=f.exports.useReducer(_,{presenting:!1,activeIndex:-1,media:(e==null?void 0:e.media)||[]});f.exports.useEffect(()=>{d({type:"replaceMedia",media:(e==null?void 0:e.media)||[]})},[e==null?void 0:e.media]),g({dispatchMedia:d,openPresentMode:m=>{d({type:"openPresentMode",activeIndex:m.state.activeIndex})}});let c=null;return e?e.subAlbums.length>0&&(c=t(u,{albums:e.subAlbums,getCustomLink:a})):c=t(u,{}),v("div",{ref:r,children:[p&&t(R,{onlyFavorites:s,setOnlyFavorites:n,setOrdering:l,ordering:o}),t(h,{album:e,disableLink:!0}),c,t(G,{loading:i,mediaState:A,dispatchMedia:d})]})});function S({getParam:e,setParams:i}){const a=e("orderBy","date_shot"),n=e("orderDirection","ASC")||"hello",l=f.exports.useCallback(({orderBy:o,orderDirection:s})=>{const r=[];o!==void 0&&r.push({key:"orderBy",value:o}),s!==void 0&&r.push({key:"orderDirection",value:s}),i(r)},[i]);return{orderBy:a,orderDirection:n,setOrdering:l}}export{k as A,B as a,S as u};
