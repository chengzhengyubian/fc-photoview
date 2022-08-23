import{g as h,c as L,i as A,d as y,u as v,r as O,j as l,a as $,L as E}from"./index.ed54fe83.js";import{A as x,u as B,a as F}from"./useOrderingParams.5953a038.js";import{u as D}from"./useURLParameters.1aca2057.js";import{u as I,P as M}from"./MediaGallery.6f345196.js";import"./AlbumBoxes.33e38eeb.js";import"./MapboxMap.27d07c02.js";import"./Table.941991f2.js";import"./Input.1cb7b307.js";import"./Dropdown.73d202b6.js";const G=h`
  ${x}

  query albumQuery(
    $id: ID!
    $onlyFavorites: Boolean
    $mediaOrderBy: String
    $orderDirection: OrderDirection
    $limit: Int
    $offset: Int
  ) {
    album(id: $id) {
      ...AlbumGalleryFields
    }
  }
`;let n=!1,d=!1;function T(){const{id:t}=L();if(A(t))throw new Error("Expected parameter `id` to be defined for AlbumPage");const{t:m}=y(),i=D(),a=B(i),u=i.getParam("favorites")=="1",s=e=>i.setParam("favorites",e?"1":"0"),{loading:o,error:c,data:r,refetch:f,fetchMore:g}=v(G,{variables:{id:t,onlyFavorites:u,mediaOrderBy:a.orderBy,orderDirection:a.orderDirection,offset:0,limit:200}}),{containerElem:b,finished:p}=I({loading:o,fetchMore:g,data:r,getItems:e=>e.album.media}),P=O.exports.useCallback(e=>{n&&!e||d&&e?f({id:t,onlyFavorites:e}).then(()=>{e?d=!1:n=!1,s(e)}):s(e)},[s,f]);return c?l("div",{children:"Error"}):$(E,{title:r?r.album.title:m("title.loading_album","Loading album"),children:[l(F,{ref:b,album:r&&r.album,loading:o,setOnlyFavorites:P,onlyFavorites:u,onFavorite:()=>n=d=!0,showFilter:!0,setOrdering:a.setOrdering,ordering:a}),l(M,{active:!p&&!o,text:m("general.loading.paginate.media","Loading more media")})]})}export{T as default};
