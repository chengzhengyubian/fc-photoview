import{g as S,s as m,d as h,u as P,j as e,a as u,L as y,P as E,n as O,r as b,S as R,M as _,p as D,K as $,Q as M,c as v,i as w,T as I,U as k}from"./index.ed54fe83.js";import{u as V,a as U}from"./useOrderingParams.5953a038.js";import{u as q}from"./useURLParameters.1aca2057.js";import{u as N,P as Q,d as W}from"./MediaGallery.6f345196.js";import{u as C}from"./index.esm.4b8c3923.js";import{T as B}from"./Input.1cb7b307.js";import"./AlbumBoxes.33e38eeb.js";import"./MapboxMap.27d07c02.js";import"./Table.941991f2.js";import"./Dropdown.73d202b6.js";const j=S`
  query shareAlbumQuery(
    $id: ID!
    $token: String!
    $password: String
    $mediaOrderBy: String
    $mediaOrderDirection: OrderDirection
    $limit: Int
    $offset: Int
  ) {
    album(id: $id, tokenCredentials: { token: $token, password: $password }) {
      id
      title
      subAlbums(order: { order_by: "title" }) {
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
        order: {
          order_by: $mediaOrderBy
          order_direction: $mediaOrderDirection
        }
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
        downloads {
          title
          mediaUrl {
            url
            width
            height
            fileSize
          }
        }
        highRes {
          url
          width
          height
        }
        videoWeb {
          url
        }
        exif {
          id
          description
          camera
          maker
          lens
          dateShot
          exposure
          aperture
          iso
          focalLength
          flash
          exposureProgram
          coordinates {
            latitude
            longitude
          }
        }
      }
    }
  }
`,z=m.div`
  height: 100%;
`,x=({albumID:r,token:t,password:a})=>{const{t:i}=h(),d=q(),s=V(d),{data:o,error:n,loading:c,fetchMore:g}=P(j,{variables:{id:r,token:t,password:a,limit:200,offset:0,mediaOrderBy:s.orderBy,mediaOrderDirection:s.orderDirection}}),{containerElem:l,finished:L}=N({loading:c,fetchMore:g,data:o,getItems:f=>f.album.media});if(n)return e("div",{children:n.message});const p=o==null?void 0:o.album;return e(z,{"data-testid":"AlbumSharePage",children:u(y,{title:p?p.title:i("general.loading.album","Loading album"),children:[e(U,{ref:l,album:p,customAlbumLink:f=>`/share/${t}/${f}`,showFilter:!0,setOrdering:s.setOrdering,ordering:s}),e(Q,{active:!L&&!c,text:i("general.loading.paginate.media","Loading more media")})]})})},F=m(E)`
  /* width: 100%; */
  max-height: calc(80vh);
  object-fit: contain;
`,K=m(O)`
  /* width: 100%; */
  max-height: calc(80vh);
`,Y=({media:r})=>{var a;const{updateSidebar:t}=b.exports.useContext(R);switch(b.exports.useEffect(()=>{t(e(W,{media:r,hidePreview:!0}))},[r]),r.type){case _.Photo:return e(F,{src:(a=r.highRes)==null?void 0:a.url});case _.Video:return e(K,{media:r})}D(r.type)},G=({media:r})=>{const{t}=h();return e(y,{title:t("share_page.media.title","Shared media"),children:u("div",{"data-testid":"MediaSharePage",children:[e("h1",{className:"font-semibold text-xl mb-4",children:r.title}),e(Y,{media:r})]})})},H=({refetchWithPassword:r,loading:t=!1})=>{const{t:a}=h(),{register:i,watch:d,formState:{errors:s},handleSubmit:o}=C(),[n,c]=b.exports.useState(!1),g=()=>{r(d("password")),c(!0)};let l;return n&&!t?l=a("share_page.wrong_password","Wrong password, please try again."):s.password&&(l=a("share_page.protected_share.password_required_error","Password is required")),u(A,{children:[e("h1",{className:"text-xl",children:a("share_page.protected_share.title","Protected share")}),e("p",{className:"mb-4",children:a("share_page.protected_share.description","This share is protected with a password.")}),e(B,{...i("password",{required:!0}),label:a("login_page.field.password","Password"),type:"password",loading:t,disabled:t,action:o(g),error:l,fullWidth:!0,sizeVariant:"big"})]})},J=S`
  query SharePageToken($token: String!, $password: String) {
    shareToken(credentials: { token: $token, password: $password }) {
      token
      album {
        id
      }
      media {
        id
        title
        type
        thumbnail {
          url
          width
          height
        }
        downloads {
          title
          mediaUrl {
            url
            width
            height
            fileSize
          }
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
        exif {
          id
          description
          camera
          maker
          lens
          dateShot
          exposure
          aperture
          iso
          focalLength
          flash
          exposureProgram
          coordinates {
            longitude
            latitude
          }
        }
      }
    }
  }
`,X=S`
  query ShareTokenValidatePassword($token: String!, $password: String) {
    shareTokenValidatePassword(
      credentials: { token: $token, password: $password }
    )
  }
`,T=()=>{const{token:r}=v();if(w(r))throw new Error("Expected `token` param to be defined");return r},Z=()=>{const{t:r}=h(),t=T(),a=$(t),{loading:i,error:d,data:s}=P(J,{variables:{token:t,password:a}});return w(d)?i?e("div",{children:r("general.loading.default","Loading...")}):s!=null&&s.shareToken.album?u(I,{children:[e(k,{path:":subAlbum",element:e(()=>{const{subAlbum:n}=v();if(w(n))throw new Error("Expected `subAlbum` param to be defined");return e(x,{albumID:n,token:t,password:a})},{})}),e(k,{index:!0,element:e(x,{albumID:s.shareToken.album.id,token:t,password:a})})]}):s!=null&&s.shareToken.media?e(G,{media:s.shareToken.media}):e("h1",{children:r("share_page.share_not_found","Share not found")}):e("div",{children:d.message})},A=m.div`
  max-width: 400px;
  margin: 100px auto 0;
`,he=()=>{const{t:r}=h(),t=T(),{loading:a,error:i,data:d,refetch:s}=P(X,{notifyOnNetworkStatusChange:!0,variables:{token:t,password:$(t)}});return i?i.message=="GraphQL error: share not found"?u(A,{children:[e("h1",{children:r("share_page.share_not_found","Share not found")}),e("p",{children:r("share_page.share_not_found_description","Maybe the share has expired or has been deleted.")})]}):e("div",{children:i.message}):d&&d.shareTokenValidatePassword==!1?e(H,{refetchWithPassword:o=>{M(t,o),s({token:t,password:o})},loading:a}):a?e("div",{children:r("general.loading.default","Loading...")}):e(Z,{})};export{A as MessageContainer,J as SHARE_TOKEN_QUERY,he as TokenRoute,X as VALIDATE_TOKEN_PASSWORD_QUERY};
