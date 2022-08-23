import{A as e}from"./AlbumBoxes.33e38eeb.js";import{g as s,u as l,j as r,L as m}from"./index.ed54fe83.js";const o=s`
  query getMyAlbums {
    myAlbums(order: { order_by: "title" }, onlyRoot: true, showEmpty: true) {
      id
      title
      thumbnail {
        id
        thumbnail {
          url
        }
      }
    }
  }
`,i=()=>{const{error:u,data:t}=l(o);return r(m,{title:"Albums",children:r(e,{error:u,albums:t==null?void 0:t.myAlbums})})};export{i as default};
