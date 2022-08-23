import{g as i,s as n,J as o,j as r}from"./index.ed54fe83.js";const g=i`
  query CheckInitialSetup {
    siteInfo {
      initialSetup
    }
  }
`;function d(e){o(e),window.location.href="/"}const c=n.div.attrs({className:"mt-20"})``,u=({message:e,show:a,type:s})=>{if(!a)return null;let t="bg-gray-100";return s=="positive"&&(t="bg-green-200 text-green-900"),s=="negative"&&(t="bg-red-200 text-red-900"),r("div",{className:`py-2 px-3 my-4 rounded-md ${t}`,children:e})};export{c as C,g as I,u as M,d as l};
