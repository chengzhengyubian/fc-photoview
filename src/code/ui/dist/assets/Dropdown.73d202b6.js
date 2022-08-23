import{s as p,j as o,f as i}from"./index.ed54fe83.js";const b=p.select`
  appearance: none;

  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='9px' height='5px' viewBox='0 0 9 5'%3E%3Cpolygon fill='%23D8D8D8' points='0 0 8.36137659 0 4.1806883 4.1806883'%3E%3C/polygon%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: center right 10px;
`,u=({items:a,selected:n,setSelected:t,className:d,...r})=>{const s=e=>{t(e.target.value),r.onChange&&r.onChange(e)},g=a.map(({value:e,label:l})=>o("option",{value:e,children:l},e));return o(b,{className:i("bg-gray-50 px-2 py-0.5 pr-6 rounded border border-gray-200 focus:outline-none focus:border-blue-300 text-[#222] hover:bg-gray-100 disabled:hover:bg-gray-50 disabled:text-gray-500 disabled:cursor-default","dark:bg-dark-input-bg dark:border-dark-input-border dark:text-dark-input-text dark:focus:border-blue-300",d),value:n,onChange:s,...r,children:g})};export{u as D};
