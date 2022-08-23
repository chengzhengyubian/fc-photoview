import{g as N,d as W,m as q,r as w,e as P,u as C,j as t,a as _}from"./index.ed54fe83.js";import{C as T,M as $,l as E,I as M}from"./MessageBox.19e1a4a3.js";import{u as z}from"./index.esm.4b8c3923.js";import{u as j,T as u,S as k}from"./Input.1cb7b307.js";const U=N`
  mutation InitialSetup(
    $username: String!
    $password: String!
    $rootPath: String!
  ) {
    initialSetupWizard(
      username: $username
      password: $password
      rootPath: $rootPath
    ) {
      success
      status
      token
    }
  }
`,B=()=>{var d,m,c,h,g,f;const{t:s}=W(),p=q(),{register:r,handleSubmit:x,formState:{errors:i}}=z();w.exports.useEffect(()=>{P()&&p("/")},[]);const{data:o}=C(M),l=((d=o==null?void 0:o.siteInfo)==null?void 0:d.initialSetup)===!1;w.exports.useEffect(()=>{l&&p("/")},[l]);const[b,{loading:I,data:e}]=j(U,{onCompleted:a=>{if(!a.initialSetupWizard)return;const{success:y,token:S}=a.initialSetupWizard;y&&S&&E(S)}}),v=x(a=>{b({variables:{username:a.username,password:a.password,rootPath:a.rootPath}})});if(P()||l)return null;let n=null;return e&&!((m=e==null?void 0:e.initialSetupWizard)!=null&&m.success)&&(n=(c=e==null?void 0:e.initialSetupWizard)==null?void 0:c.status),t("div",{children:_(T,{children:[t("h1",{className:"text-center text-xl",children:s("login_page.initial_setup.title","Initial Setup")}),_("form",{onSubmit:v,className:"max-w-[500px] mx-auto",children:[t(u,{wrapperClassName:"my-4",fullWidth:!0,...r("username",{required:!0}),label:s("login_page.field.username","Username"),error:((h=i.username)==null?void 0:h.type)=="required"?"Please enter a username":void 0}),t(u,{wrapperClassName:"my-4",fullWidth:!0,...r("password",{required:!0}),label:s("login_page.field.password","Password"),error:((g=i.password)==null?void 0:g.type)=="required"?"Please enter a password":void 0}),t(u,{wrapperClassName:"my-4",fullWidth:!0,...r("rootPath",{required:!0}),label:s("login_page.initial_setup.field.photo_path.label","Photo path"),placeholder:s("login_page.initial_setup.field.photo_path.placeholder","/path/to/photos"),error:((f=i.password)==null?void 0:f.type)=="required"?"Please enter a photo path":void 0}),t($,{type:"negative",message:n,show:!!n}),t(k,{className:"mt-2",disabled:I,children:s("login_page.initial_setup.field.submit","Setup Photoview")})]})]})})};export{B as default};
