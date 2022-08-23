import{r as d,R as A,s as y,g as M,d as S,u as de,a as f,j as o,h as ue,e as V,S as xe,A as Mt,b as ne,t as we,P as Ce,M as be,k as me,N as ke,i as N,l as Pt,F as G,m as ye,L as at,c as Ft,f as Nt,n as Dt}from"./index.ed54fe83.js";import{m as Rt,a as Tt,o as At,P as Ot}from"./MapboxMap.27d07c02.js";import{u as T,T as se,b as Bt,B as te}from"./Input.1cb7b307.js";import{t as Z,a as ot,H,y as B,L as nt,S as lt,b as it,c as st,u as q,p as j,_ as W,I as z,o as $,d as E,r as Ae,x as le,n as _e,s as Oe,v as Ut,T as Gt,E as zt,e as ct,h as Me,f as Pe,g as re,i as ae,j as oe,k as Ht,C as Be,l as ge,m as dt,q as ie,w as ve,z as ut,A as Fe,M as Ue}from"./Table.941991f2.js";var Wt=function(){var e=document.getSelection();if(!e.rangeCount)return function(){};for(var t=document.activeElement,a=[],n=0;n<e.rangeCount;n++)a.push(e.getRangeAt(n));switch(t.tagName.toUpperCase()){case"INPUT":case"TEXTAREA":t.blur();break;default:t=null;break}return e.removeAllRanges(),function(){e.type==="Caret"&&e.removeAllRanges(),e.rangeCount||a.forEach(function(r){e.addRange(r)}),t&&t.focus()}},Qt=Wt,je={"text/plain":"Text","text/html":"Url",default:"Text"},Xt="Copy to clipboard: #{key}, Enter";function qt(e){var t=(/mac os x/i.test(navigator.userAgent)?"\u2318":"Ctrl")+"+C";return e.replace(/#{\s*key\s*}/g,t)}function Yt(e,t){var a,n,r,l,s,i,c=!1;t||(t={}),a=t.debug||!1;try{r=Qt(),l=document.createRange(),s=document.getSelection(),i=document.createElement("span"),i.textContent=e,i.style.all="unset",i.style.position="fixed",i.style.top=0,i.style.clip="rect(0, 0, 0, 0)",i.style.whiteSpace="pre",i.style.webkitUserSelect="text",i.style.MozUserSelect="text",i.style.msUserSelect="text",i.style.userSelect="text",i.addEventListener("copy",function(p){if(p.stopPropagation(),t.format)if(p.preventDefault(),typeof p.clipboardData=="undefined"){a&&console.warn("unable to use e.clipboardData"),a&&console.warn("trying IE specific stuff"),window.clipboardData.clearData();var h=je[t.format]||je.default;window.clipboardData.setData(h,e)}else p.clipboardData.clearData(),p.clipboardData.setData(t.format,e);t.onCopy&&(p.preventDefault(),t.onCopy(p.clipboardData))}),document.body.appendChild(i),l.selectNodeContents(i),s.addRange(l);var u=document.execCommand("copy");if(!u)throw new Error("copy command was unsuccessful");c=!0}catch(p){a&&console.error("unable to copy using execCommand: ",p),a&&console.warn("trying IE specific stuff");try{window.clipboardData.setData(t.format||"text",e),t.onCopy&&t.onCopy(window.clipboardData),c=!0}catch(h){a&&console.error("unable to copy using clipboardData: ",h),a&&console.error("falling back to prompt"),n=qt("message"in t?t.message:Xt),window.prompt(n,e)}}finally{s&&(typeof s.removeRange=="function"?s.removeRange(l):s.removeAllRanges()),i&&document.body.removeChild(i),r()}return c}var Zt=Yt;function ce(){let e=[],t=[],a={enqueue(n){t.push(n)},addEventListener(n,r,l,s){return n.addEventListener(r,l,s),a.add(()=>n.removeEventListener(r,l,s))},requestAnimationFrame(...n){let r=requestAnimationFrame(...n);return a.add(()=>cancelAnimationFrame(r))},nextFrame(...n){return a.requestAnimationFrame(()=>a.requestAnimationFrame(...n))},setTimeout(...n){let r=setTimeout(...n);return a.add(()=>clearTimeout(r))},add(n){return e.push(n),()=>{let r=e.indexOf(n);if(r>=0){let[l]=e.splice(r,1);l()}}},dispose(){for(let n of e.splice(0))n()},async workQueue(){for(let n of t.splice(0))await n()}};return a}function pt(){let[e]=d.exports.useState(ce);return d.exports.useEffect(()=>()=>e.dispose(),[e]),e}function Ke(e){var t;if(e.type)return e.type;let a=(t=e.as)!=null?t:"button";if(typeof a=="string"&&a.toLowerCase()==="button")return"button"}function mt(e,t){let[a,n]=d.exports.useState(()=>Ke(e));return Z(()=>{n(Ke(e))},[e.type,e.as]),Z(()=>{a||!t.current||t.current instanceof HTMLButtonElement&&!t.current.hasAttribute("type")&&n("button")},[a,t]),a}function Vt({container:e,accept:t,walk:a,enabled:n=!0}){let r=d.exports.useRef(t),l=d.exports.useRef(a);d.exports.useEffect(()=>{r.current=t,l.current=a},[t,a]),Z(()=>{if(!e||!n)return;let s=ot(e);if(!s)return;let i=r.current,c=l.current,u=Object.assign(h=>i(h),{acceptNode:i}),p=s.createTreeWalker(e,NodeFilter.SHOW_ELEMENT,u,!1);for(;p.nextNode();)c(p.currentNode)},[e,n,r,l])}function jt(e){throw new Error("Unexpected object: "+e)}var O=(e=>(e[e.First=0]="First",e[e.Previous=1]="Previous",e[e.Next=2]="Next",e[e.Last=3]="Last",e[e.Specific=4]="Specific",e[e.Nothing=5]="Nothing",e))(O||{});function Kt(e,t){let a=t.resolveItems();if(a.length<=0)return null;let n=t.resolveActiveIndex(),r=n!=null?n:-1,l=(()=>{switch(e.focus){case 0:return a.findIndex(s=>!t.resolveDisabled(s));case 1:{let s=a.slice().reverse().findIndex((i,c,u)=>r!==-1&&u.length-c-1>=r?!1:!t.resolveDisabled(i));return s===-1?s:a.length-1-s}case 2:return a.findIndex((s,i)=>i<=r?!1:!t.resolveDisabled(s));case 3:{let s=a.slice().reverse().findIndex(i=>!t.resolveDisabled(i));return s===-1?s:a.length-1-s}case 4:return a.findIndex(s=>t.resolveId(s)===e.id);case 5:return null;default:jt(e)}})();return l===-1?n:l}var Jt=(e=>(e[e.Open=0]="Open",e[e.Closed=1]="Closed",e))(Jt||{}),er=(e=>(e[e.Pointer=0]="Pointer",e[e.Other=1]="Other",e))(er||{}),tr=(e=>(e[e.OpenMenu=0]="OpenMenu",e[e.CloseMenu=1]="CloseMenu",e[e.GoToItem=2]="GoToItem",e[e.Search=3]="Search",e[e.ClearSearch=4]="ClearSearch",e[e.RegisterItem=5]="RegisterItem",e[e.UnregisterItem=6]="UnregisterItem",e))(tr||{});function $e(e,t=a=>a){let a=e.activeItemIndex!==null?e.items[e.activeItemIndex]:null,n=Ut(t(e.items.slice()),l=>l.dataRef.current.domRef.current),r=a?n.indexOf(a):null;return r===-1&&(r=null),{items:n,activeItemIndex:r}}let rr={[1](e){return e.menuState===1?e:{...e,activeItemIndex:null,menuState:1}},[0](e){return e.menuState===0?e:{...e,menuState:0}},[2]:(e,t)=>{var a;let n=$e(e),r=Kt(t,{resolveItems:()=>n.items,resolveActiveIndex:()=>n.activeItemIndex,resolveId:l=>l.id,resolveDisabled:l=>l.dataRef.current.disabled});return{...e,...n,searchQuery:"",activeItemIndex:r,activationTrigger:(a=t.trigger)!=null?a:1}},[3]:(e,t)=>{let a=e.searchQuery!==""?0:1,n=e.searchQuery+t.value.toLowerCase(),r=(e.activeItemIndex!==null?e.items.slice(e.activeItemIndex+a).concat(e.items.slice(0,e.activeItemIndex+a)):e.items).find(s=>{var i;return((i=s.dataRef.current.textValue)==null?void 0:i.startsWith(n))&&!s.dataRef.current.disabled}),l=r?e.items.indexOf(r):-1;return l===-1||l===e.activeItemIndex?{...e,searchQuery:n}:{...e,searchQuery:n,activeItemIndex:l,activationTrigger:1}},[4](e){return e.searchQuery===""?e:{...e,searchQuery:"",searchActiveItemIndex:null}},[5]:(e,t)=>{let a=$e(e,n=>[...n,{id:t.id,dataRef:t.dataRef}]);return{...e,...a}},[6]:(e,t)=>{let a=$e(e,n=>{let r=n.findIndex(l=>l.id===t.id);return r!==-1&&n.splice(r,1),n});return{...e,...a,activationTrigger:1}}},Ge=d.exports.createContext(null);Ge.displayName="MenuContext";function Ie(e){let t=d.exports.useContext(Ge);if(t===null){let a=new Error(`<${e} /> is missing a parent <Menu /> component.`);throw Error.captureStackTrace&&Error.captureStackTrace(a,Ie),a}return t}function ar(e,t){return q(t.type,rr,e,t)}let or=d.exports.Fragment,nr=H(function(e,t){let a=d.exports.useReducer(ar,{menuState:1,buttonRef:d.exports.createRef(),itemsRef:d.exports.createRef(),items:[],searchQuery:"",activeItemIndex:null,activationTrigger:1}),[{menuState:n,itemsRef:r,buttonRef:l},s]=a,i=B(t);nt([l,r],(h,w)=>{var x;s({type:1}),lt(w,it.Loose)||(h.preventDefault(),(x=l.current)==null||x.focus())},n===0);let c=d.exports.useMemo(()=>({open:n===0}),[n]),u=e,p={ref:i};return A.createElement(Ge.Provider,{value:a},A.createElement(st,{value:q(n,{[0]:j.Open,[1]:j.Closed})},W({ourProps:p,theirProps:u,slot:c,defaultTag:or,name:"Menu"})))}),lr="button",ir=H(function(e,t){var a;let[n,r]=Ie("Menu.Button"),l=B(n.buttonRef,t),s=`headlessui-menu-button-${z()}`,i=pt(),c=$(g=>{switch(g.key){case E.Space:case E.Enter:case E.ArrowDown:g.preventDefault(),g.stopPropagation(),r({type:0}),i.nextFrame(()=>r({type:2,focus:O.First}));break;case E.ArrowUp:g.preventDefault(),g.stopPropagation(),r({type:0}),i.nextFrame(()=>r({type:2,focus:O.Last}));break}}),u=$(g=>{switch(g.key){case E.Space:g.preventDefault();break}}),p=$(g=>{if(Ae(g.currentTarget))return g.preventDefault();e.disabled||(n.menuState===0?(r({type:1}),i.nextFrame(()=>{var b;return(b=n.buttonRef.current)==null?void 0:b.focus({preventScroll:!0})})):(g.preventDefault(),r({type:0})))}),h=d.exports.useMemo(()=>({open:n.menuState===0}),[n]),w=e,x={ref:l,id:s,type:mt(e,n.buttonRef),"aria-haspopup":!0,"aria-controls":(a=n.itemsRef.current)==null?void 0:a.id,"aria-expanded":e.disabled?void 0:n.menuState===0,onKeyDown:c,onKeyUp:u,onClick:p};return W({ourProps:x,theirProps:w,slot:h,defaultTag:lr,name:"Menu.Button"})}),sr="div",cr=le.RenderStrategy|le.Static,dr=H(function(e,t){var a,n;let[r,l]=Ie("Menu.Items"),s=B(r.itemsRef,t),i=_e(r.itemsRef),c=`headlessui-menu-items-${z()}`,u=pt(),p=Oe(),h=(()=>p!==null?p===j.Open:r.menuState===0)();d.exports.useEffect(()=>{let m=r.itemsRef.current;!m||r.menuState===0&&m!==(i==null?void 0:i.activeElement)&&m.focus({preventScroll:!0})},[r.menuState,r.itemsRef,i]),Vt({container:r.itemsRef.current,enabled:r.menuState===0,accept(m){return m.getAttribute("role")==="menuitem"?NodeFilter.FILTER_REJECT:m.hasAttribute("role")?NodeFilter.FILTER_SKIP:NodeFilter.FILTER_ACCEPT},walk(m){m.setAttribute("role","none")}});let w=$(m=>{var I,P;switch(u.dispose(),m.key){case E.Space:if(r.searchQuery!=="")return m.preventDefault(),m.stopPropagation(),l({type:3,value:m.key});case E.Enter:if(m.preventDefault(),m.stopPropagation(),l({type:1}),r.activeItemIndex!==null){let{dataRef:k}=r.items[r.activeItemIndex];(P=(I=k.current)==null?void 0:I.domRef.current)==null||P.click()}ce().nextFrame(()=>{var k;return(k=r.buttonRef.current)==null?void 0:k.focus({preventScroll:!0})});break;case E.ArrowDown:return m.preventDefault(),m.stopPropagation(),l({type:2,focus:O.Next});case E.ArrowUp:return m.preventDefault(),m.stopPropagation(),l({type:2,focus:O.Previous});case E.Home:case E.PageUp:return m.preventDefault(),m.stopPropagation(),l({type:2,focus:O.First});case E.End:case E.PageDown:return m.preventDefault(),m.stopPropagation(),l({type:2,focus:O.Last});case E.Escape:m.preventDefault(),m.stopPropagation(),l({type:1}),ce().nextFrame(()=>{var k;return(k=r.buttonRef.current)==null?void 0:k.focus({preventScroll:!0})});break;case E.Tab:m.preventDefault(),m.stopPropagation();break;default:m.key.length===1&&(l({type:3,value:m.key}),u.setTimeout(()=>l({type:4}),350));break}}),x=$(m=>{switch(m.key){case E.Space:m.preventDefault();break}}),g=d.exports.useMemo(()=>({open:r.menuState===0}),[r]),b=e,C={"aria-activedescendant":r.activeItemIndex===null||(a=r.items[r.activeItemIndex])==null?void 0:a.id,"aria-labelledby":(n=r.buttonRef.current)==null?void 0:n.id,id:c,onKeyDown:w,onKeyUp:x,role:"menu",tabIndex:0,ref:s};return W({ourProps:C,theirProps:b,slot:g,defaultTag:sr,features:cr,visible:h,name:"Menu.Items"})}),ur=d.exports.Fragment,pr=H(function(e,t){let{disabled:a=!1,...n}=e,[r,l]=Ie("Menu.Item"),s=`headlessui-menu-item-${z()}`,i=r.activeItemIndex!==null?r.items[r.activeItemIndex].id===s:!1,c=d.exports.useRef(null),u=B(t,c);Z(()=>{if(r.menuState!==0||!i||r.activationTrigger===0)return;let C=ce();return C.requestAnimationFrame(()=>{var m,I;(I=(m=c.current)==null?void 0:m.scrollIntoView)==null||I.call(m,{block:"nearest"})}),C.dispose},[c,i,r.menuState,r.activationTrigger,r.activeItemIndex]);let p=d.exports.useRef({disabled:a,domRef:c});Z(()=>{p.current.disabled=a},[p,a]),Z(()=>{var C,m;p.current.textValue=(m=(C=c.current)==null?void 0:C.textContent)==null?void 0:m.toLowerCase()},[p,c]),Z(()=>(l({type:5,id:s,dataRef:p}),()=>l({type:6,id:s})),[p,s]);let h=$(C=>{if(a)return C.preventDefault();l({type:1}),ce().nextFrame(()=>{var m;return(m=r.buttonRef.current)==null?void 0:m.focus({preventScroll:!0})})}),w=$(()=>{if(a)return l({type:2,focus:O.Nothing});l({type:2,focus:O.Specific,id:s})}),x=$(()=>{a||i||l({type:2,focus:O.Specific,id:s,trigger:0})}),g=$(()=>{a||!i||l({type:2,focus:O.Nothing})}),b=d.exports.useMemo(()=>({active:i,disabled:a}),[i,a]);return W({ourProps:{id:s,ref:u,role:"menuitem",tabIndex:a===!0?void 0:-1,"aria-disabled":a===!0?!0:void 0,disabled:void 0,onClick:h,onFocus:w,onPointerMove:x,onMouseMove:x,onPointerLeave:g,onMouseLeave:g},theirProps:n,slot:b,defaultTag:ur,name:"Menu.Item"})}),fe=Object.assign(nr,{Button:ir,Items:dr,Item:pr});var mr=(e=>(e[e.Open=0]="Open",e[e.Closed=1]="Closed",e))(mr||{}),hr=(e=>(e[e.TogglePopover=0]="TogglePopover",e[e.ClosePopover=1]="ClosePopover",e[e.SetButton=2]="SetButton",e[e.SetButtonId=3]="SetButtonId",e[e.SetPanel=4]="SetPanel",e[e.SetPanelId=5]="SetPanelId",e))(hr||{});let fr={[0]:e=>({...e,popoverState:q(e.popoverState,{[0]:1,[1]:0})}),[1](e){return e.popoverState===1?e:{...e,popoverState:1}},[2](e,t){return e.button===t.button?e:{...e,button:t.button}},[3](e,t){return e.buttonId===t.buttonId?e:{...e,buttonId:t.buttonId}},[4](e,t){return e.panel===t.panel?e:{...e,panel:t.panel}},[5](e,t){return e.panelId===t.panelId?e:{...e,panelId:t.panelId}}},ze=d.exports.createContext(null);ze.displayName="PopoverContext";function Se(e){let t=d.exports.useContext(ze);if(t===null){let a=new Error(`<${e} /> is missing a parent <Popover /> component.`);throw Error.captureStackTrace&&Error.captureStackTrace(a,Se),a}return t}let He=d.exports.createContext(null);He.displayName="PopoverAPIContext";function We(e){let t=d.exports.useContext(He);if(t===null){let a=new Error(`<${e} /> is missing a parent <Popover /> component.`);throw Error.captureStackTrace&&Error.captureStackTrace(a,We),a}return t}let Qe=d.exports.createContext(null);Qe.displayName="PopoverGroupContext";function ht(){return d.exports.useContext(Qe)}let Xe=d.exports.createContext(null);Xe.displayName="PopoverPanelContext";function br(){return d.exports.useContext(Xe)}function gr(e,t){return q(t.type,fr,e,t)}let vr="div",xr=H(function(e,t){var a;let n=`headlessui-popover-button-${z()}`,r=`headlessui-popover-panel-${z()}`,l=d.exports.useRef(null),s=B(t,Gt(v=>{l.current=v})),i=d.exports.useReducer(gr,{popoverState:1,button:null,buttonId:n,panel:null,panelId:r,beforePanelSentinel:d.exports.createRef(),afterPanelSentinel:d.exports.createRef()}),[{popoverState:c,button:u,panel:p,beforePanelSentinel:h,afterPanelSentinel:w},x]=i,g=_e((a=l.current)!=null?a:u);d.exports.useEffect(()=>x({type:3,buttonId:n}),[n,x]),d.exports.useEffect(()=>x({type:5,panelId:r}),[r,x]);let b=d.exports.useMemo(()=>{if(!u||!p)return!1;for(let v of document.querySelectorAll("body > *"))if(Number(v==null?void 0:v.contains(u))^Number(v==null?void 0:v.contains(p)))return!0;return!1},[u,p]),C=d.exports.useMemo(()=>({buttonId:n,panelId:r,close:()=>x({type:1})}),[n,r,x]),m=ht(),I=m==null?void 0:m.registerPopover,P=$(()=>{var v;return(v=m==null?void 0:m.isFocusWithinPopoverGroup())!=null?v:(g==null?void 0:g.activeElement)&&((u==null?void 0:u.contains(g.activeElement))||(p==null?void 0:p.contains(g.activeElement)))});d.exports.useEffect(()=>I==null?void 0:I(C),[I,C]),zt(g==null?void 0:g.defaultView,"focus",v=>{var L,F,Q,Y;c===0&&(P()||!u||!p||(F=(L=h.current)==null?void 0:L.contains)!=null&&F.call(L,v.target)||(Y=(Q=w.current)==null?void 0:Q.contains)!=null&&Y.call(Q,v.target)||x({type:1}))},!0),nt([u,p],(v,L)=>{x({type:1}),lt(L,it.Loose)||(v.preventDefault(),u==null||u.focus())},c===0);let k=$(v=>{x({type:1});let L=(()=>v?v instanceof HTMLElement?v:v.current instanceof HTMLElement?v.current:u:u)();L==null||L.focus()}),_=d.exports.useMemo(()=>({close:k,isPortalled:b}),[k,b]),R=d.exports.useMemo(()=>({open:c===0,close:k}),[c,k]),D=e,U={ref:s};return A.createElement(ze.Provider,{value:i},A.createElement(He.Provider,{value:_},A.createElement(st,{value:q(c,{[0]:j.Open,[1]:j.Closed})},W({ourProps:U,theirProps:D,slot:R,defaultTag:vr,name:"Popover"}))))}),wr="button",Cr=H(function(e,t){let[a,n]=Se("Popover.Button"),{isPortalled:r}=We("Popover.Button"),l=d.exports.useRef(null),s=`headlessui-focus-sentinel-${z()}`,i=ht(),c=i==null?void 0:i.closeOthers,u=br(),p=u===null?!1:u===a.panelId,h=B(l,t,p?null:v=>n({type:2,button:v})),w=B(l,t),x=_e(l),g=$(v=>{var L,F,Q;if(p){if(a.popoverState===1)return;switch(v.key){case E.Space:case E.Enter:v.preventDefault(),(F=(L=v.target).click)==null||F.call(L),n({type:1}),(Q=a.button)==null||Q.focus();break}}else switch(v.key){case E.Space:case E.Enter:v.preventDefault(),v.stopPropagation(),a.popoverState===1&&(c==null||c(a.buttonId)),n({type:0});break;case E.Escape:if(a.popoverState!==0)return c==null?void 0:c(a.buttonId);if(!l.current||(x==null?void 0:x.activeElement)&&!l.current.contains(x.activeElement))return;v.preventDefault(),v.stopPropagation(),n({type:1});break}}),b=$(v=>{p||v.key===E.Space&&v.preventDefault()}),C=$(v=>{var L,F;Ae(v.currentTarget)||e.disabled||(p?(n({type:1}),(L=a.button)==null||L.focus()):(v.preventDefault(),v.stopPropagation(),a.popoverState===1&&(c==null||c(a.buttonId)),n({type:0}),(F=a.button)==null||F.focus()))}),m=$(v=>{v.preventDefault(),v.stopPropagation()}),I=a.popoverState===0,P=d.exports.useMemo(()=>({open:I}),[I]),k=mt(e,l),_=e,R=p?{ref:w,type:k,onKeyDown:g,onClick:C}:{ref:h,id:a.buttonId,type:k,"aria-expanded":e.disabled?void 0:a.popoverState===0,"aria-controls":a.panel?a.panelId:void 0,onKeyDown:g,onKeyUp:b,onClick:C,onMouseDown:m},D=ct(),U=$(()=>{let v=a.panel;if(!v)return;function L(){q(D.current,{[oe.Forwards]:()=>re(v,ae.First),[oe.Backwards]:()=>re(v,ae.Last)})}L()});return A.createElement(A.Fragment,null,W({ourProps:R,theirProps:_,slot:P,defaultTag:wr,name:"Popover.Button"}),I&&!p&&r&&A.createElement(Me,{id:s,features:Pe.Focusable,as:"button",type:"button",onFocus:U}))}),yr="div",_r=le.RenderStrategy|le.Static,Ir=H(function(e,t){let[{popoverState:a},n]=Se("Popover.Overlay"),r=B(t),l=`headlessui-popover-overlay-${z()}`,s=Oe(),i=(()=>s!==null?s===j.Open:a===0)(),c=$(p=>{if(Ae(p.currentTarget))return p.preventDefault();n({type:1})}),u=d.exports.useMemo(()=>({open:a===0}),[a]);return W({ourProps:{ref:r,id:l,"aria-hidden":!0,onClick:c},theirProps:e,slot:u,defaultTag:yr,features:_r,visible:i,name:"Popover.Overlay"})}),Sr="div",kr=le.RenderStrategy|le.Static,$r=H(function(e,t){let{focus:a=!1,...n}=e,[r,l]=Se("Popover.Panel"),{close:s,isPortalled:i}=We("Popover.Panel"),c=`headlessui-focus-sentinel-before-${z()}`,u=`headlessui-focus-sentinel-after-${z()}`,p=d.exports.useRef(null),h=B(p,t,_=>{l({type:4,panel:_})}),w=_e(p),x=Oe(),g=(()=>x!==null?x===j.Open:r.popoverState===0)(),b=$(_=>{var R;switch(_.key){case E.Escape:if(r.popoverState!==0||!p.current||(w==null?void 0:w.activeElement)&&!p.current.contains(w.activeElement))return;_.preventDefault(),_.stopPropagation(),l({type:1}),(R=r.button)==null||R.focus();break}});d.exports.useEffect(()=>{var _;e.static||r.popoverState===1&&((_=e.unmount)!=null?_:!0)&&l({type:4,panel:null})},[r.popoverState,e.unmount,e.static,l]),d.exports.useEffect(()=>{if(!a||r.popoverState!==0||!p.current)return;let _=w==null?void 0:w.activeElement;p.current.contains(_)||re(p.current,ae.First)},[a,p,r.popoverState]);let C=d.exports.useMemo(()=>({open:r.popoverState===0,close:s}),[r,s]),m={ref:h,id:r.panelId,onKeyDown:b,onBlur:a&&r.popoverState===0?_=>{var R,D,U,v,L;let F=_.relatedTarget;!F||!p.current||(R=p.current)!=null&&R.contains(F)||(l({type:1}),(((U=(D=r.beforePanelSentinel.current)==null?void 0:D.contains)==null?void 0:U.call(D,F))||((L=(v=r.afterPanelSentinel.current)==null?void 0:v.contains)==null?void 0:L.call(v,F)))&&F.focus({preventScroll:!0}))}:void 0,tabIndex:-1},I=ct(),P=$(()=>{let _=p.current;if(!_)return;function R(){q(I.current,{[oe.Forwards]:()=>{re(_,ae.First)},[oe.Backwards]:()=>{var D;(D=r.button)==null||D.focus({preventScroll:!0})}})}R()}),k=$(()=>{let _=p.current;if(!_)return;function R(){q(I.current,{[oe.Forwards]:()=>{var D,U,v;if(!r.button)return;let L=Ht(),F=L.indexOf(r.button),Q=L.slice(0,F+1),Y=[...L.slice(F+1),...Q];for(let pe of Y.slice())if(((U=(D=pe==null?void 0:pe.id)==null?void 0:D.startsWith)==null?void 0:U.call(D,"headlessui-focus-sentinel-"))||((v=r.panel)==null?void 0:v.contains(pe))){let Ve=Y.indexOf(pe);Ve!==-1&&Y.splice(Ve,1)}re(Y,ae.First,!1)},[oe.Backwards]:()=>re(_,ae.Last)})}R()});return A.createElement(Xe.Provider,{value:r.panelId},g&&i&&A.createElement(Me,{id:c,ref:r.beforePanelSentinel,features:Pe.Focusable,as:"button",type:"button",onFocus:P}),W({ourProps:m,theirProps:n,slot:C,defaultTag:Sr,features:kr,visible:g,name:"Popover.Panel"}),g&&i&&A.createElement(Me,{id:u,ref:r.afterPanelSentinel,features:Pe.Focusable,as:"button",type:"button",onFocus:k}))}),Lr="div",Er=H(function(e,t){let a=d.exports.useRef(null),n=B(a,t),[r,l]=d.exports.useState([]),s=$(g=>{l(b=>{let C=b.indexOf(g);if(C!==-1){let m=b.slice();return m.splice(C,1),m}return b})}),i=$(g=>(l(b=>[...b,g]),()=>s(g))),c=$(()=>{var g;let b=ot(a);if(!b)return!1;let C=b.activeElement;return(g=a.current)!=null&&g.contains(C)?!0:r.some(m=>{var I,P;return((I=b.getElementById(m.buttonId))==null?void 0:I.contains(C))||((P=b.getElementById(m.panelId))==null?void 0:P.contains(C))})}),u=$(g=>{for(let b of r)b.buttonId!==g&&b.close()}),p=d.exports.useMemo(()=>({registerPopover:i,unregisterPopover:s,isFocusWithinPopoverGroup:c,closeOthers:u}),[i,s,c,u]),h=d.exports.useMemo(()=>({}),[]),w=e,x={ref:n};return A.createElement(Qe.Provider,{value:p},W({ourProps:x,theirProps:w,slot:h,defaultTag:Lr,name:"Popover.Group"}))}),Le=Object.assign(xr,{Button:Cr,Overlay:Ir,Panel:$r,Group:Er});const J=y.div.attrs({className:"mb-8"})``,ee=y.h2.attrs({className:"text-xl mb-2 mx-4"})``,Mr=e=>d.exports.createElement("svg",{width:"12px",height:"14px",viewBox:"0 0 12 12",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("path",{stroke:"none",strokeWidth:1,fill:"currentColor",d:"M5.01241676,3.99743718 L5.2123926,4.0051107 C6.08060614,4.06729045 6.87900923,4.50349496 7.40038252,5.20051071 C7.56578589,5.42163587 7.52061445,5.73497915 7.29948929,5.90038252 C7.07836413,6.06578589 6.76502085,6.02061445 6.59961748,5.79948929 C6.25203528,5.33481212 5.71976656,5.04400912 5.14095753,5.00255596 C4.62002941,4.96524811 4.10763194,5.1332244 3.71124008,5.46642286 L3.58355339,5.58355339 L2.0896441,7.0773559 C1.33169968,7.86211367 1.34253936,9.10952545 2.11400695,9.88099305 C2.84487099,10.6118571 4.00289823,10.6600523 4.78475751,10.0238225 L4.91144661,9.91144661 L5.76644661,9.05644661 C5.96170876,8.86118446 6.27829124,8.86118446 6.47355339,9.05644661 C6.64711974,9.23001296 6.66640489,9.49943736 6.53140884,9.6943055 L6.47355339,9.76355339 L5.6123559,10.6246441 C4.43521924,11.7615607 2.56410157,11.7453012 1.40690017,10.5880998 C0.295986831,9.47718649 0.236565536,7.70832791 1.24484904,6.52165836 L1.37644661,6.37644661 L2.87638252,4.87651071 C3.44442608,4.30826116 4.21444071,3.99361146 5.01241676,3.99743718 Z M10.5930998,1.41190017 C11.7040132,2.52281351 11.7634345,4.29167209 10.755151,5.47834164 L10.6235534,5.62355339 L9.12361748,7.12348929 C8.50823696,7.73909297 7.65582094,8.05706905 6.7876074,7.9948893 C5.91939386,7.93270955 5.12099077,7.49650504 4.59961748,6.79948929 C4.43421411,6.57836413 4.47938555,6.26502085 4.70051071,6.09961748 C4.92163587,5.93421411 5.23497915,5.97938555 5.40038252,6.20051071 C5.74796472,6.66518788 6.28023344,6.95599088 6.85904247,6.99744404 C7.37997059,7.03475189 7.89236806,6.8667756 8.28875992,6.53357714 L8.41644661,6.41644661 L9.9103559,4.9226441 C10.6683003,4.13788633 10.6574606,2.89047455 9.88599305,2.11900695 C9.15512901,1.38814291 7.99710177,1.3399477 7.21436478,1.97705267 L7.08752112,2.08958265 L6.22752112,2.94458265 C6.03169053,3.13927469 5.71510939,3.13835172 5.52041735,2.94252112 C5.34735776,2.76844949 5.32885818,2.49897001 5.46442178,2.30449627 L5.52247888,2.23541735 L6.3876441,1.3753559 C7.56478076,0.238439269 9.43589843,0.254698775 10.5930998,1.41190017 Z",fillRule:"nonzero"})),Pr=e=>d.exports.createElement("svg",{width:"14px",height:"14px",viewBox:"0 0 14 14",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("path",{stroke:"none",d:"M12.0909091,4.45454545 C13.1452709,4.45454545 14,5.30927457 14,6.36363636 L14,6.36363636 L14,12.0909091 C14,13.1452709 13.1452709,14 12.0909091,14 L12.0909091,14 L6.36363636,14 C5.30927457,14 4.45454545,13.1452709 4.45454545,12.0909091 L4.45454545,12.0909091 L4.45454545,6.36363636 C4.45454545,5.30927457 5.30927457,4.45454545 6.36363636,4.45454545 L6.36363636,4.45454545 Z M12.0909091,5.72727273 L6.36363636,5.72727273 C6.01218243,5.72727273 5.72727273,6.01218243 5.72727273,6.36363636 L5.72727273,6.36363636 L5.72727273,12.0909091 C5.72727273,12.442363 6.01218243,12.7272727 6.36363636,12.7272727 L6.36363636,12.7272727 L12.0909091,12.7272727 C12.442363,12.7272727 12.7272727,12.442363 12.7272727,12.0909091 L12.7272727,12.0909091 L12.7272727,6.36363636 C12.7272727,6.01218243 12.442363,5.72727273 12.0909091,5.72727273 L12.0909091,5.72727273 Z M7.63636364,0 C8.65306965,0 9.48414783,0.79476725 9.54221374,1.79691732 L9.54545455,1.90909091 L9.54545455,2.54545455 C9.54545455,2.89690848 9.26054484,3.18181818 8.90909091,3.18181818 C8.58274083,3.18181818 8.31376817,2.93615624 8.27700856,2.61966799 L8.27272727,2.54545455 L8.27272727,1.90909091 C8.27272727,1.58274083 8.02706533,1.31376817 7.71057708,1.27700856 L7.63636364,1.27272727 L1.90909091,1.27272727 C1.58274083,1.27272727 1.31376817,1.51838921 1.27700856,1.83487747 L1.27272727,1.90909091 L1.27272727,7.63636364 C1.27272727,7.96271372 1.51838921,8.23168637 1.83487747,8.26844599 L1.90909091,8.27272727 L2.54545455,8.27272727 C2.89690848,8.27272727 3.18181818,8.55763698 3.18181818,8.90909091 C3.18181818,9.23544099 2.93615624,9.50441365 2.61966799,9.54117326 L2.54545455,9.54545455 L1.90909091,9.54545455 C0.892384892,9.54545455 0.0613067131,8.7506873 0.0032408049,7.74853723 L0,7.63636364 L0,1.90909091 C0,0.892384892 0.79476725,0.0613067131 1.79691732,0.0032408049 L1.90909091,0 L7.63636364,0 Z",fill:"currentColor",fillRule:"nonzero"})),Fr=e=>d.exports.createElement("svg",{width:"13px",height:"14px",viewBox:"0 0 13 14",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("path",{stroke:"none",d:"M7.63636364,0 C8.65306965,0 9.48414783,0.79476725 9.54221374,1.79691732 L9.54545455,1.90909091 L9.545,2.545 L12.0909091,2.54545455 C12.442363,2.54545455 12.7272727,2.83036425 12.7272727,3.18181818 C12.7272727,3.50816826 12.4816108,3.77714092 12.1651225,3.81390053 L12.0909091,3.81818182 L11.454,3.818 L11.4545455,12.0909091 C11.4545455,13.1076151 10.6597782,13.9386933 9.65762814,13.9967592 L9.54545455,14 L3.18181818,14 C2.16511217,14 1.33403399,13.2052327 1.27596808,12.2030827 L1.27272727,12.0909091 L1.272,3.818 L0.636363636,3.81818182 C0.284909705,3.81818182 2.27373675e-13,3.53327211 2.27373675e-13,3.18181818 C2.27373675e-13,2.8554681 0.245661939,2.58649544 0.562150193,2.54973583 L0.636363636,2.54545455 L3.181,2.545 L3.18181818,1.90909091 C3.18181818,0.892384892 3.97658543,0.0613067131 4.9787355,0.0032408049 L5.09090909,0 L7.63636364,0 Z M10.181,3.818 L2.545,3.818 L2.54545455,12.0909091 C2.54545455,12.4172592 2.79111648,12.6862318 3.10760474,12.7229914 L3.18181818,12.7272727 L9.54545455,12.7272727 C9.87180462,12.7272727 10.1407773,12.4816108 10.1775369,12.1651225 L10.1818182,12.0909091 L10.181,3.818 Z M5.09090909,5.72727273 C5.41725917,5.72727273 5.68623183,5.97293467 5.72299144,6.28942292 L5.72727273,6.36363636 L5.72727273,10.1818182 C5.72727273,10.5332721 5.44236302,10.8181818 5.09090909,10.8181818 C4.76455901,10.8181818 4.49558635,10.5725199 4.45882674,10.2560316 L4.45454545,10.1818182 L4.45454545,6.36363636 C4.45454545,6.01218243 4.73945516,5.72727273 5.09090909,5.72727273 Z M7.63636364,5.72727273 C7.96271372,5.72727273 8.23168637,5.97293467 8.26844599,6.28942292 L8.27272727,6.36363636 L8.27272727,10.1818182 C8.27272727,10.5332721 7.98781757,10.8181818 7.63636364,10.8181818 C7.31001356,10.8181818 7.0410409,10.5725199 7.00428128,10.2560316 L7,10.1818182 L7,6.36363636 C7,6.01218243 7.2849097,5.72727273 7.63636364,5.72727273 Z M7.63636364,1.27272727 L5.09090909,1.27272727 C4.76455901,1.27272727 4.49558635,1.51838921 4.45882674,1.83487747 L4.45454545,1.90909091 L4.454,2.545 L8.272,2.545 L8.27272727,1.90909091 C8.27272727,1.58274083 8.02706533,1.31376817 7.71057708,1.27700856 L7.63636364,1.27272727 Z",fill:"currentColor",fillRule:"nonzero"})),Nr=e=>d.exports.createElement("svg",{width:"10px",height:"14px",viewBox:"0 0 10 2",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("path",{stroke:"currentColor",strokeWidth:1,fill:"none",d:"M0.646446609,0.646446609 C1.13807119,0.5 1.26307119,0.555964406 1.35355339,0.646446609 C1.44403559,0.736928813 1.5,0.861928813 1.5,1 C1.5,1.13807119 1.44403559,1.26307119 1.35355339,1.35355339 C1.26307119,1.44403559 1.13807119,1.5 1,1.5 C0.861928813,1.5 0.736928813,1.44403559 0.646446609,1.35355339 C0.555964406,1.26307119 0.5,1.13807119 0.5,1 C0.5,0.861928813 0.555964406,0.736928813 0.646446609,0.646446609 Z M4.64644661,0.646446609 C5.13807119,0.5 5.26307119,0.555964406 5.35355339,0.646446609 C5.44403559,0.736928813 5.5,0.861928813 5.5,1 C5.5,1.13807119 5.44403559,1.26307119 5.35355339,1.35355339 C5.26307119,1.44403559 5.13807119,1.5 5,1.5 C4.86192881,1.5 4.73692881,1.44403559 4.64644661,1.35355339 C4.55596441,1.26307119 4.5,1.13807119 4.5,1 C4.5,0.861928813 4.55596441,0.736928813 4.64644661,0.646446609 Z M8.64644661,0.646446609 C9.13807119,0.5 9.26307119,0.555964406 9.35355339,0.646446609 C9.44403559,0.736928813 9.5,0.861928813 9.5,1 C9.5,1.13807119 9.44403559,1.26307119 9.35355339,1.35355339 C9.26307119,1.44403559 9.13807119,1.5 9,1.5 C8.86192881,1.5 8.73692881,1.44403559 8.64644661,1.35355339 C8.55596441,1.26307119 8.5,1.13807119 8.5,1 C8.5,0.861928813 8.55596441,0.736928813 8.64644661,0.646446609 Z"})),Dr=e=>d.exports.createElement("svg",{width:"12px",height:"12px",viewBox:"0 0 12 12",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("path",{d:"M5.66666667,1 L5.66666667,10.3333333 M1,5.66666667 L10.3333333,5.66666667",fill:"none",stroke:"currentColor",strokeWidth:1.5})),Ne=M`
  query sidebarGetPhotoShares($id: ID!) {
    media(id: $id) {
      id
      shares {
        id
        token
        hasPassword
      }
    }
  }
`,De=M`
  query sidebarGetAlbumShares($id: ID!) {
    album(id: $id) {
      id
      shares {
        id
        token
        hasPassword
      }
    }
  }
`,Rr=M`
  mutation sidebarPhotoAddShare($id: ID!, $password: String, $expire: Time) {
    shareMedia(mediaId: $id, password: $password, expire: $expire) {
      token
    }
  }
`,Tr=M`
  mutation sidebarAlbumAddShare($id: ID!, $password: String, $expire: Time) {
    shareAlbum(albumId: $id, password: $password, expire: $expire) {
      token
    }
  }
`,Ar=M`
  mutation sidebarProtectShare($token: String!, $password: String) {
    protectShareToken(token: $token, password: $password) {
      token
      hasPassword
    }
  }
`,Or=M`
  mutation sidebareDeleteShare($token: String!) {
    deleteShareToken(token: $token) {
      token
    }
  }
`,ft=y.div.attrs({className:"absolute -top-3 bg-white dark:bg-dark-bg rounded shadow-md border border-gray-200 dark:border-dark-border z-10"})`
  width: ${({width:e})=>e}px;

  ${({flipped:e})=>e?`
      left: 32px;
        `:`
      right: 24px;
    `}

  &::after {
    content: '';
    position: absolute;
    top: 18px;
    width: 8px;
    height: 14px;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 14'%3E%3Cpolyline stroke-width='1' stroke='%23E2E2E2' fill='%23FFFFFF' points='1 0 7 7 1 14'%3E%3C/polyline%3E%3C/svg%3E");

    ${({flipped:e})=>e?`
      left: -7px;
      transform: rotate(180deg);
        `:`
      right: -7px;
    `}
  }
`,Br=({share:e,query:t,id:a})=>{const[n,r]=d.exports.useState(!1),l=n||e.hasPassword,[s,i]=d.exports.useState(e.hasPassword?"**********":""),[c,u]=d.exports.useState(e.hasPassword),[p,{loading:h}]=T(Ar,{refetchQueries:[{query:t,variables:{id:a}}],onCompleted:b=>{w(b.protectShareToken.hasPassword)},variables:{token:e.token}}),w=b=>{b&&i("**********"),c&&!b&&i(""),u(b)};return f("div",{className:"px-4 py-2",children:[o(Be,{label:"Password protected",checked:l,onChange:()=>{const b=!l;r(b),b||(p({variables:{token:e.token,password:null}}),i(""))}}),o(se,{disabled:!l,type:c?"password":"text",value:s,className:"mt-2 w-full",onKeyDown:b=>{b.shiftKey||b.altKey||b.ctrlKey||b.metaKey||b.key=="Enter"||b.key=="Tab"||b.key=="Escape"||w(!1)},onChange:b=>{i(b.target.value)},action:()=>{!c&&s!=""&&p({variables:{token:e.token,password:s}})},loading:h})]})},Ur=({id:e,share:t,query:a})=>{const{t:n}=S();return f(Le,{className:"relative",children:[o(Le.Button,{className:"align-middle p-1 ml-2",title:n("sidebar.sharing.more","More"),children:o(Nr,{})}),o(Le.Panel,{children:f(ft,{width:260,children:[o(Br,{id:e,share:t,query:a}),f("div",{className:"px-4 py-2 border-t border-gray-200 dark:border-dark-border mt-2 mb-2",children:[o(Be,{label:"Expiration date"}),o(se,{className:"mt-2 w-full"})]})]})})]})},Gr=({id:e})=>{const{t}=S(),{loading:a,error:n,data:r}=de(De,{variables:{id:e}}),[l,{loading:s}]=T(Tr,{refetchQueries:[{query:De,variables:{id:e}}]}),i=a||s;return n?f("div",{children:["Error: ",n.message]}):i?o("div",{children:t("general.loading.shares","Loading shares...")}):o(bt,{id:e,isPhoto:!1,loading:i,shares:r==null?void 0:r.album.shares,shareItem:l})},zr=({id:e})=>{const{t}=S(),[a,{loading:n,error:r,data:l}]=ue(Ne),[s,{loading:i}]=T(Rr,{refetchQueries:[{query:Ne,variables:{id:e}}]});d.exports.useEffect(()=>{V()&&a({variables:{id:e}})},[]);const c=n||i;return r?f("div",{children:["Error: ",r.message]}):c?o("div",{children:t("general.loading.shares","Loading shares...")}):o(bt,{id:e,isPhoto:!0,loading:c,shares:l==null?void 0:l.media.shares,shareItem:s})},bt=({loading:e,shares:t,isPhoto:a,id:n,shareItem:r})=>{const{t:l}=S(),s=a?Ne:De,[i]=T(Or,{refetchQueries:[{query:s,variables:{id:n}}]});if(t===void 0)return null;const c=t.map(u=>f("tr",{className:"border-gray-100 dark:border-dark-border2 border-b border-t",children:[f("td",{className:"pl-4 py-2 w-full",children:[f("span",{className:"text-[#585858] dark:text-[#C0C3C4] mr-2",children:[o(Mr,{className:"inline-block mr-2"}),o("span",{className:"text-xs uppercase font-bold",children:l("sidebar.sharing.public_link","Public Link")+" "})]}),o("span",{className:"text-sm",children:u.token})]}),f("td",{className:"pr-6 py-2 whitespace-nowrap text-[#5C6A7F] dark:text-[#7599ca] flex",children:[o("button",{className:"align-middle p-1 ml-2",title:l("sidebar.sharing.copy_link","Copy Link"),onClick:()=>{Zt(`${location.origin}/share/${u.token}`)},children:o(Pr,{})}),o("button",{onClick:()=>{i({variables:{token:u.token}})},className:"align-middle p-1 ml-2 hover:text-red-600 focus:text-red-600",title:l("sidebar.sharing.delete","Delete"),children:o(Fr,{})}),o(Ur,{share:u,id:n,query:s})]})]},u.token));return c.length==0&&c.push(o("tr",{className:"border-gray-100 dark:border-dark-border2 border-b border-t",children:o("td",{colSpan:2,className:"pl-4 py-2 italic text-gray-600 dark:text-gray-300",children:l("sidebar.sharing.no_shares_found","No shares found")})},"no-shares")),f(J,{children:[o(ee,{children:l("sidebar.sharing.title","Sharing options")}),o("div",{children:f("table",{className:"border-collapse w-full",children:[o("tbody",{children:c}),o("tfoot",{children:o("tr",{className:"text-left border-gray-100 dark:border-dark-border2 border-b border-t",children:o("td",{colSpan:2,className:"pl-4 py-2",children:f("button",{className:"text-green-500 font-bold uppercase text-xs",disabled:e,onClick:()=>{r({variables:{id:n}})},children:[o(Dr,{className:"inline-block mr-2"}),o("span",{children:l("sidebar.sharing.add_share","Add shares")})]})})})})]})})]})},Hr=e=>d.exports.createElement("svg",{width:"12px",height:"12px",viewBox:"0 0 12 12",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("path",{stroke:"currentColor",strokeWidth:2,fill:"none",d:"M1,11 L11,1 M11,11 L1,1"})),Wr=e=>d.exports.createElement("svg",{width:"16px",height:"17px",viewBox:"0 0 16 17",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("g",{stroke:"none",strokeWidth:1,fill:"none",fillRule:"evenodd"},d.exports.createElement("g",{transform:"translate(8, 8.5) rotate(-330) translate(-8, -8.5) translate(4, 2)",stroke:"currentColor"},d.exports.createElement("path",{d:"M6.76850191,0.434840963 C6.04896114,1.86467297 5.79279384,3.20120746 6,4.44444444 C6.17989677,5.52382504 6.63266329,6.48919052 7.35829958,7.34054088 C7.50158762,7.50867286 7.48146228,7.76112984 7.31333862,7.90442763 C7.24096223,7.96611661 7.14897427,8 7.053875,8 L0.922873286,8 C0.701919553,8.00007212 0.522801161,7.82095373 0.522801161,7.59999999 C0.522801161,7.50599929 0.555900634,7.41500023 0.616293108,7.34296634 C1.36567981,6.44912577 1.82691544,5.48295181 2,4.44444444 C2.19754347,3.25918364 1.93688287,1.92116674 1.21801821,0.430393743 C1.14595115,0.281147735 1.2086138,0.10177184 1.35789863,0.0297852978 C1.39855358,0.0101811032 1.44310701,8.29112604e-18 1.4882418,0 L6.50055336,0 C6.666219,3.58349532e-05 6.80051754,0.134334371 6.80051754,0.300000012 C6.80051754,0.346829895 6.78955305,0.393009303 6.76850191,0.434840963 Z",strokeLinejoin:"round"}),d.exports.createElement("line",{x1:4,y1:8,x2:4,y2:13,strokeLinecap:"round"})))),Qr=e=>d.exports.createElement("svg",{width:"16px",height:"17px",viewBox:"0 0 16 17",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("g",{stroke:"none",strokeWidth:1,fill:"none",fillRule:"evenodd"},d.exports.createElement("g",{transform:"translate(8, 8.5) rotate(-330) translate(-8, -8.5) translate(4, 2)",stroke:"currentColor",fill:"currentColor"},d.exports.createElement("path",{d:"M6.76850191,0.434840963 C6.04896114,1.86467297 5.79279384,3.20120746 6,4.44444444 C6.17989677,5.52382504 6.63266329,6.48919052 7.35829958,7.34054088 C7.50158762,7.50867286 7.48146228,7.76112984 7.31333862,7.90442763 C7.24096223,7.96611661 7.14897427,8 7.053875,8 L0.922873286,8 C0.701919553,8.00007212 0.522801161,7.82095373 0.522801161,7.59999999 C0.522801161,7.50599929 0.555900634,7.41500023 0.616293108,7.34296634 C1.36567981,6.44912577 1.82691544,5.48295181 2,4.44444444 C2.19754347,3.25918364 1.93688287,1.92116674 1.21801821,0.430393743 C1.14595115,0.281147735 1.2086138,0.10177184 1.35789863,0.0297852978 C1.39855358,0.0101811032 1.44310701,8.29112604e-18 1.4882418,0 L6.50055336,0 C6.666219,3.58349532e-05 6.80051754,0.134334371 6.80051754,0.300000012 C6.80051754,0.346829895 6.78955305,0.393009303 6.76850191,0.434840963 Z",strokeLinejoin:"round"}),d.exports.createElement("line",{x1:4,y1:8,x2:4,y2:13,strokeLinecap:"round"})))),gt=({title:e})=>{const{updateSidebar:t,setPinned:a,pinned:n}=d.exports.useContext(xe);return f("div",{className:"m-2 flex items-center",children:[o("button",{className:`${n?"lg:hidden":""}`,title:"Close sidebar",onClick:()=>t(null),children:o(Hr,{className:"m-2 text-[#1F2021] dark:text-[#abadaf]"})}),o("span",{className:"flex-grow -mt-1 ml-2",children:e}),o("button",{className:"hidden lg:block",title:"Pin sidebar",onClick:()=>a(!n),children:o(n?Qr:Wr,{className:"m-2 text-[#1F2021] dark:text-[#abadaf]"})})]})},Xr=M`
  mutation resetAlbumCover($albumID: ID!) {
    resetAlbumCover(albumID: $albumID) {
      id
      thumbnail {
        id
        thumbnail {
          url
        }
      }
    }
  }
`,qr=M`
  mutation setAlbumCover($coverID: ID!) {
    setAlbumCover(coverID: $coverID) {
      id
      thumbnail {
        id
        thumbnail {
          url
        }
      }
    }
  }
`,Yr=({cover_id:e})=>{const{t}=S(),[a]=T(qr,{variables:{coverID:e}}),[n,r]=d.exports.useState(!1);return d.exports.useEffect(()=>{r(!1)},[e]),V()?f(J,{children:[o(ee,{children:t("sidebar.album.album_cover","Album cover")}),o("div",{children:o("table",{className:"border-collapse w-full",children:o("tfoot",{children:o("tr",{className:"text-left border-gray-100 dark:border-dark-border2 border-b border-t",children:o("td",{colSpan:2,className:"pl-4 py-2",children:o("button",{className:"disabled:opacity-50 text-green-500 font-bold uppercase text-xs",disabled:n,onClick:()=>{r(!0),a({variables:{coverID:e}})},children:o("span",{children:t("sidebar.album.set_cover","Set as album cover photo")})})})})})})})]}):null},Zr=({id:e})=>{const{t}=S(),[a]=T(Xr,{variables:{albumID:e}}),[n,r]=d.exports.useState(!1);return d.exports.useEffect(()=>{r(!1)},[e]),f(J,{children:[o(ee,{children:t("sidebar.album.album_cover","Album cover")}),o("div",{children:o("table",{className:"border-collapse w-full",children:o("tfoot",{children:o("tr",{className:"text-left border-gray-100 dark:border-dark-border2 border-b border-t",children:o("td",{colSpan:2,className:"pl-4 py-2",children:o("button",{className:"disabled:opacity-50 text-red-500 font-bold uppercase text-xs",disabled:n,onClick:()=>{r(!0),a({variables:{albumID:e}})},children:o("span",{children:t("sidebar.album.reset_cover","Reset cover photo")})})})})})})})]})},Vr=y.table.attrs({className:"table-fixed w-full"})``,jr=y.thead.attrs({className:"bg-[#f9f9fb] dark:bg-[#2B3037]"})``,Kr=y.tr.attrs({className:"text-left uppercase text-xs border-gray-100 dark:border-dark-border2 border-b border-t"})``,Jr=y.tr.attrs({className:"cursor-pointer border-gray-100 dark:border-dark-border2 border-b hover:bg-gray-50 focus:bg-gray-50 dark:hover:bg-[#3c4759] dark:focus:bg-[#3c4759]"})``;var X={Table:Vr,Head:jr,HeadRow:Kr,Row:Jr};const ea=({albumID:e})=>{const{t}=S(),n=[{title:t("sidebar.album.download.thumbnails.title","Thumbnails"),description:t("sidebar.album.download.thumbnails.description","Low resolution images, no videos"),purpose:"thumbnail,video-thumbnail"},{title:t("sidebar.album.download.high-resolutions.title","High resolutions"),description:t("sidebar.album.download.high-resolutions.description","High resolution jpegs of RAW images"),purpose:"high-res"},{title:t("sidebar.album.download.originals.title","Originals"),description:t("sidebar.album.download.originals.description","The original images and videos"),purpose:"original"},{title:t("sidebar.album.download.web-videos.title","Converted videos"),description:t("sidebar.album.download.web-videos.description","Videos that have been optimized for web"),purpose:"video-web"}].map(r=>f(X.Row,{onClick:()=>location.href=`${Mt}/download/album/${e}/${r.purpose}`,tabIndex:0,children:[o("td",{className:"pl-4 py-2",children:`${r.title}`}),o("td",{className:"pr-4 py-2 text-sm text-gray-800 dark:text-gray-400 italic",children:`${r.description}`})]},r.purpose));return f(J,{children:[o(ee,{children:t("sidebar.download.title","Download")}),f(X.Table,{children:[o(X.Head,{children:o(X.HeadRow,{children:o("th",{className:"px-4 py-2",colSpan:2,children:t("sidebar.download.table_columns.name","Name")})})}),o("tbody",{children:n})]})]})},ta=M`
  query getAlbumSidebar($id: ID!) {
    album(id: $id) {
      id
      title
    }
  }
`,ra=({albumId:e})=>{var l;const{t}=S(),{loading:a,error:n,data:r}=de(ta,{variables:{id:e}});return a?o("div",{children:t("general.loading.default","Loading...")}):n?o("div",{children:n.message}):f("div",{children:[o(gt,{title:(l=r==null?void 0:r.album.title)!=null?l:t("sidebar.album.title_placeholder","Album title")}),o("div",{className:"mt-8",children:o(Gr,{id:e})}),o("div",{className:"mt-8",children:o(Zr,{id:e})}),o("div",{className:"mt-8",children:o(ea,{albumID:e})})]})};function aa(e,t=[]){const a=d.exports.useState(!1)[1],n=d.exports.useRef(!1);return d.exports.useLayoutEffect(()=>{const r=setTimeout(()=>{n.current=!0,a(l=>!l)},e);return()=>{n.current=!1,clearTimeout(r)}},t),n.current}const oa=e=>d.exports.createElement("svg",{viewBox:"0 0 15 15",width:"15px",height:"15px",fill:"currentColor",...e},d.exports.createElement("path",{d:"M7.33333333,0 C8.069713,0 8.66666667,0.596953667 8.66666667,1.33333333 L8.66666667,1.39333333 C8.66841527,1.83176195 8.93035293,2.22728781 9.33333333,2.4 C9.74472671,2.5815643 10.2252088,2.49444392 10.5466667,2.18 L10.5866667,2.14 C10.8367577,1.88963061 11.176121,1.74895112 11.53,1.74895112 C11.883879,1.74895112 12.2232423,1.88963061 12.4733333,2.14 C12.7237027,2.39009101 12.8643822,2.72945434 12.8643822,3.08333333 C12.8643822,3.43721233 12.7237027,3.77657566 12.4733333,4.02666667 L12.4333333,4.06666667 C12.1188894,4.38812449 12.031769,4.86860662 12.2133333,5.28 L12.2133333,5.33333333 C12.3860455,5.73631374 12.7815714,5.9982514 13.22,6 L13.3333333,6 C14.069713,6 14.6666667,6.59695367 14.6666667,7.33333333 C14.6666667,8.069713 14.069713,8.66666667 13.3333333,8.66666667 L13.2733333,8.66666667 C12.8349047,8.66841527 12.4393789,8.93035293 12.2666667,9.33333333 C12.0851024,9.74472671 12.1722227,10.2252088 12.4866667,10.5466667 L12.5266667,10.5866667 C12.7770361,10.8367577 12.9177155,11.176121 12.9177155,11.53 C12.9177155,11.883879 12.7770361,12.2232423 12.5266667,12.4733333 C12.2765757,12.7237027 11.9372123,12.8643822 11.5833333,12.8643822 C11.2294543,12.8643822 10.890091,12.7237027 10.64,12.4733333 L10.6,12.4333333 C10.2785422,12.1188894 9.79806005,12.031769 9.38666667,12.2133333 C8.98368626,12.3860455 8.7217486,12.7815714 8.72,13.22 L8.72,13.3333333 C8.72,14.069713 8.12304633,14.6666667 7.38666667,14.6666667 C6.650287,14.6666667 6.05333333,14.069713 6.05333333,13.3333333 L6.05333333,13.2733333 C6.04277107,12.8217805 5.75724785,12.4225768 5.33333333,12.2666667 C4.92193995,12.0851024 4.44145782,12.1722227 4.12,12.4866667 L4.08,12.5266667 C3.82990899,12.7770361 3.49054566,12.9177155 3.13666667,12.9177155 C2.78278767,12.9177155 2.44342434,12.7770361 2.19333333,12.5266667 C1.94296394,12.2765757 1.80228446,11.9372123 1.80228446,11.5833333 C1.80228446,11.2294543 1.94296394,10.890091 2.19333333,10.64 L2.23333333,10.6 C2.54777725,10.2785422 2.63489764,9.79806005 2.45333333,9.38666667 C2.28062114,8.98368626 1.88509528,8.7217486 1.44666667,8.72 L1.33333333,8.72 C0.596953667,8.72 0,8.12304633 0,7.38666667 C0,6.650287 0.596953667,6.05333333 1.33333333,6.05333333 L1.39333333,6.05333333 C1.84488612,6.04277107 2.24408988,5.75724785 2.4,5.33333333 C2.5815643,4.92193995 2.49444392,4.44145782 2.18,4.12 L2.14,4.08 C1.88963061,3.82990899 1.74895112,3.49054566 1.74895112,3.13666667 C1.74895112,2.78278767 1.88963061,2.44342434 2.14,2.19333333 C2.39009101,1.94296394 2.72945434,1.80228446 3.08333333,1.80228446 C3.43721233,1.80228446 3.77657566,1.94296394 4.02666667,2.19333333 L4.06666667,2.23333333 C4.38812449,2.54777725 4.86860662,2.63489764 5.28,2.45333333 L5.33333333,2.45333333 C5.73631374,2.28062114 5.9982514,1.88509528 6,1.44666667 L6,1.33333333 C6,0.596953667 6.59695367,0 7.33333333,0 Z M7.33333333,5.33333333 C6.22876383,5.33333333 5.33333333,6.22876383 5.33333333,7.33333333 C5.33333333,8.43790283 6.22876383,9.33333333 7.33333333,9.33333333 C8.43790283,9.33333333 9.33333333,8.43790283 9.33333333,7.33333333 C9.33333333,6.22876383 8.43790283,5.33333333 7.33333333,5.33333333 Z"})),vt=y.ol`
  &
    ${({hideLastArrow:e})=>e?"li:not(:last-child)::after":"li::after"} {
    content: '';
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='5px' height='6px' viewBox='0 0 5 6'%3E%3Cpolyline fill='none' stroke='%23979797' points='0.74 0.167710644 3.57228936 3 0.74 5.83228936' /%3E%3C/svg%3E");
    width: 5px;
    height: 6px;
    display: inline-block;
    margin: 6px;
    vertical-align: middle;
  }
`,na=M`
  query albumPathQuery($id: ID!) {
    album(id: $id) {
      id
      path {
        id
        title
      }
    }
  }
`,bo=({album:e,disableLink:t=!1})=>{const[a,{data:n}]=ue(na),{updateSidebar:r}=d.exports.useContext(xe);d.exports.useEffect(()=>{!e||V()&&t==!0&&a({variables:{id:e.id}})},[e]);const l=aa(200,[e]);if(!e)return f("div",{className:`flex mb-6 flex-col h-14 transition-opacity animate-pulse ${l?"opacity-100":"opacity-0"}`,children:[o("div",{className:"w-32 h-4 bg-gray-100 mb-2 mt-1"}),o("div",{className:"w-72 h-6 bg-gray-100"})]});let s=o("span",{children:e.title});const c=((n==null?void 0:n.album.path)||[]).slice().reverse().map(u=>o("li",{className:"inline-block hover:underline",children:o(ne,{to:`/album/${u.id}`,children:u.title})},u.id));return t||(s=o(ne,{to:`/album/${e.id}`,children:s})),f("div",{className:"flex mb-6 items-end h-14",children:[f("div",{className:"min-w-0",children:[o("nav",{"aria-label":"Album breadcrumb",children:o(vt,{children:c})}),o("h1",{className:"text-2xl truncate min-w-0",children:s})]}),V()&&o("button",{title:"Album options","aria-label":"Album options",className:we(Bt({}),"px-2 py-2 ml-2"),onClick:()=>{r(o(ra,{albumId:e.id}))},children:o(oa,{})})]})},la=e=>d.exports.createElement("svg",{width:"35px",height:"44px",viewBox:"0 0 35 44",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("path",{d:"M3.07677855,0.965719632 L33.3619943,20.3146075 C34.2928109,20.9092959 34.5652967,22.145962 33.9706083,23.0767786 C33.814312,23.3214163 33.606632,23.5290962 33.3619943,23.6853925 L3.07677855,43.0342804 C2.14596197,43.6289687 0.909295855,43.356483 0.31460748,42.4256664 C0.10916386,42.1041025 4.6731219e-17,41.7304772 0,41.3488878 L0,2.65111215 C-1.3527075e-16,1.54654265 0.8954305,0.651112151 2,0.651112151 C2.38158936,0.651112151 2.75521463,0.760276012 3.07677855,0.965719632 Z",fill:"currentColor"})),Re=y.div`
  flex-grow: 1;
  flex-basis: 0;
  height: 200px;
  margin: 4px;
  background-color: #eee;
  position: relative;
  overflow: hidden;
`,ia=y(Ce)`
  height: 200px;
  min-width: 100%;
  position: relative;
  object-fit: cover;

  transition: opacity 300ms;
`,sa=e=>o(ia,{...e,lazyLoading:!0}),ca=y.div`
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;

  ${({active:e})=>e&&`
      outline: 4px solid rgba(65, 131, 196, 0.6);
      outline-offset: -4px;
    `}
`,xt=y.button`
  font-size: 1.5em;
  margin: 160px 10px 0 10px;
  color: white;
  text-shadow: 0 0 4px black;
  opacity: 0;
  position: relative;

  border-radius: 50%;
  width: 34px;
  height: 34px;

  ${Re}:hover &, ${Re}:focus-within & {
    opacity: 1 !important;
  }

  &:hover,
  &:focus {
    background-color: rgba(0, 0, 0, 0.4);
  }

  transition: opacity 100ms, background-color 100ms;
`,da=({favorite:e,onClick:t})=>o(xt,{onClick:t,style:{opacity:e?"0.75":void 0},children:o("svg",{className:"text-white m-auto mt-1",width:"19px",height:"17px",viewBox:"0 0 19 17",version:"1.1",children:o("path",{d:"M13.999086,1 C15.0573371,1 16.0710089,1.43342987 16.8190212,2.20112483 C17.5765039,2.97781012 18,4.03198704 18,5.13009709 C18,6.22820714 17.5765039,7.28238406 16.8188574,8.05923734 L16.8188574,8.05923734 L15.8553647,9.04761889 L9.49975689,15.5674041 L3.14414912,9.04761889 L2.18065643,8.05923735 C1.39216493,7.2503776 0.999999992,6.18971057 1,5.13009711 C1.00000001,4.07048366 1.39216496,3.00981663 2.18065647,2.20095689 C2.95931483,1.40218431 3.97927681,1.00049878 5.00042783,1.00049878 C6.02157882,1.00049878 7.04154078,1.4021843 7.82019912,2.20095684 L7.82019912,2.20095684 L9.4997569,3.92390079 L11.1794784,2.20078881 C11.9271631,1.43342987 12.9408349,1 13.999086,1 L13.999086,1 Z",fill:e?"currentColor":"none",stroke:"currentColor",strokeWidth:e?"0":"2"})})}),ua=({onClick:e})=>o(pa,{onClick:e,children:o("svg",{width:"20px",height:"20px",viewBox:"0 0 20 20",version:"1.1",className:"m-auto",children:o("path",{d:"M10,0 C15.5228475,0 20,4.4771525 20,10 C20,15.5228475 15.5228475,20 10,20 C4.4771525,20 0,15.5228475 0,10 C0,4.4771525 4.4771525,0 10,0 Z M10,9 C9.44771525,9 9,9.44771525 9,10 L9,10 L9,14 L9.00672773,14.1166211 C9.06449284,14.6139598 9.48716416,15 10,15 C10.5522847,15 11,14.5522847 11,14 L11,14 L11,10 L10.9932723,9.88337887 C10.9355072,9.38604019 10.5128358,9 10,9 Z M10.01,5 L9.88337887,5.00672773 C9.38604019,5.06449284 9,5.48716416 9,6 C9,6.55228475 9.44771525,7 10,7 L10,7 L10.1266211,6.99327227 C10.6239598,6.93550716 11.01,6.51283584 11.01,6 C11.01,5.44771525 10.5622847,5 10.01,5 L10.01,5 Z",fill:"#FFFFFF"})})}),pa=y(xt)`
  margin: 10px !important;
  position: absolute;
  top: 0;
  right: 0;
`,ma=y(la)`
  color: rgba(255, 255, 255, 0.8);
  position: absolute;
  left: calc(50% - 17.5px);
  top: calc(50% - 22px);
`,ha=({media:e,active:t,selectImage:a,clickPresent:n,clickFavorite:r})=>{var c;let l=null;e.favorite!==void 0&&(l=o(da,{favorite:e.favorite,onClick:u=>{u.stopPropagation(),r()}}));let s=null;e.type==be.Video&&(s=o(ma,{}));let i=100;return e.thumbnail&&(i=Math.floor(e.thumbnail.width/e.thumbnail.height*200)),f(Re,{style:{cursor:"pointer",minWidth:`clamp(124px, ${i}px, 100% - 8px)`},onClick:()=>{n()},children:[o("div",{style:{minWidth:`${i}px`,height:"200px"},children:o(sa,{src:(c=e.thumbnail)==null?void 0:c.url,blurhash:e.blurhash})}),f(ca,{active:t,children:[s,o(ua,{onClick:u=>{u.stopPropagation(),a()}}),l]})]},e.id)},fa=y.div`
  flex-grow: 1;
  height: 200px;
  width: 300px;
  margin: 4px;
  background-color: #eee;
  position: relative;
`,ba=M`
  mutation markMediaFavorite($mediaId: ID!, $favorite: Boolean!) {
    favoriteMedia(mediaId: $mediaId, favorite: $favorite) {
      id
      favorite
    }
  }
`,ga=()=>T(ba),va=({media:e,markFavorite:t})=>t({variables:{mediaId:e.id,favorite:!e.favorite},optimisticResponse:{favoriteMedia:{id:e.id,favorite:!e.favorite,__typename:"Media"}}}),xa=y(ne)`
  box-shadow: inset 0 0 2px 1px rgba(0, 0, 0, 0.3), 0 0 0 1px rgb(255, 255, 255);
  border-radius: 50%;
  position: absolute;
  top: ${({$minY:e})=>e*100}%;
  bottom: ${({$maxY:e})=>(1-e)*100}%;
  left: ${({$minX:e})=>e*100}%;
  right: ${({$maxX:e})=>(1-e)*100}%;
`,wa=({face:e})=>o(xa,{to:`/people/${e.faceGroup.id}`,$minX:e.rectangle.minX,$maxX:e.rectangle.maxX,$minY:e.rectangle.minY,$maxY:e.rectangle.maxY}),Ca=y.div`
  position: absolute;
  width: ${({width:e})=>e*100}%;
  left: ${({width:e})=>(100-e*100)/2}%;
  height: 100%;
  top: 0;
  opacity: 0;

  user-select: none;
  transition: opacity ease 200ms;

  &:hover {
    opacity: 1;
  }
`,ya=({media:e})=>{var n;if(e.type!=be.Photo||e.thumbnail==null)return null;const t=(n=e.faces)==null?void 0:n.map(r=>o(wa,{face:r},r.id));let a=1;return e.thumbnail.width*.75<e.thumbnail.height&&(a=e.thumbnail.width*.75/e.thumbnail.height),o(Ca,{width:a,children:t})},_a=M`
  query sidebarDownloadQuery($mediaId: ID!) {
    media(id: $mediaId) {
      id
      downloads {
        title
        mediaUrl {
          url
          width
          height
          fileSize
        }
      }
    }
  }
`,Te=e=>t=>{if(t==0)return e("sidebar.download.filesize.byte","{{count}} Byte",{count:0});const a=Math.floor(Math.log(t)/Math.log(1024)),n=Math.round(t/Math.pow(1024,a));switch(a){case 0:return e("sidebar.download.filesize.byte","{{count}} Byte",{count:n});case 1:return e("sidebar.download.filesize.kilo_byte","{{count}} KB",{count:n});case 2:return e("sidebar.download.filesize.mega_byte","{{count}} MB",{count:n});case 3:return e("sidebar.download.filesize.giga_byte","{{count}} GB",{count:n});case 4:return e("sidebar.download.filesize.tera_byte","{{count}} TB",{count:n});default:throw new Error(`invalid byte value: ${t}`)}},Ia=e=>async t=>{const a=new URL(t,location.origin);if(V()==null){const i=location.pathname.match(/^\/share\/([\d\w]+)(\/?.*)$/);i&&a.searchParams.set("token",i[1])}const n=await fetch(a.href,{credentials:"include"});let r=null;if(n.headers.has("content-length")?r=await Sa(e)(n):r=await n.blob(),r==null){console.log("Blob is null canceling");return}const l=t.match(/[^/]*$/);if(l==null){console.error("Could not extract filename",t);return}const s=l[0];ka(r,s)},Sa=e=>async t=>{var h;const a=Number(t.headers.get("content-length")),n=(h=t.body)==null?void 0:h.getReader(),r=new Uint8Array(a);if(n==null)throw new Error("Download reader is null");let l=!1;const s=()=>{l=!0,n.cancel("Download canceled by user")},i=Math.random().toString(26);me.add({key:i,type:ke.Progress,onDismiss:s,props:{header:"Downloading photo",content:"Starting download",percent:0}});let c=0,u;do{if(u=await n.read(),l)break;u.value&&r.set(u.value,c),c+=u.value?u.value.length:0,me.add({key:i,type:ke.Progress,onDismiss:s,props:{header:"Downloading photo",percent:c/a*100,content:`${Te(e)(c)} of ${Te(e)(a)} bytes downloaded`}})}while(!u.done);return l?void 0:(me.add({key:i,type:ke.Progress,props:{header:"Downloading photo completed",content:"The photo has been downloaded",percent:100,positive:!0}}),setTimeout(()=>{me.removeKey(i)},2e3),new Blob([r.buffer],{type:t.headers.get("content-type")||void 0}))},ka=(e,t)=>{const a=window.URL.createObjectURL(e),n=document.createElement("a");document.body.appendChild(n),n.href=a,n.download=t,n.click(),n.remove(),window.URL.revokeObjectURL(a)},$a=({rows:e})=>{const{t}=S(),a=s=>{var c;const i=s.split(/[#?]/);if(i!=null)return(c=i[0].split(".").pop())==null?void 0:c.trim().toLowerCase()},n=Ia(t),r=Te(t),l=e.map(s=>f(X.Row,{onClick:()=>n(s.url),tabIndex:0,children:[o("td",{className:"pl-4 py-2",children:`${s.title}`}),o("td",{className:"py-2",children:`${s.width} x ${s.height}`}),o("td",{className:"py-2",children:`${r(s.fileSize)}`}),o("td",{className:"pr-4 py-2",children:a(s.url)})]},s.url));return f(X.Table,{children:[o(X.Head,{children:f(X.HeadRow,{children:[o("th",{className:"w-2/6 pl-4 py-2",children:t("sidebar.download.table_columns.name","Name")}),o("th",{className:"w-2/6 py-2",children:t("sidebar.download.table_columns.dimensions","Dimensions")}),o("th",{className:"w-1/6 py-2",children:t("sidebar.download.table_columns.file_size","Size")}),o("th",{className:"w-1/6 pr-4 py-2",children:t("sidebar.download.table_columns.file_type","Type")})]})}),o("tbody",{children:l})]})},La=({media:e})=>{const{t}=S();if(!e||!e.id)return null;const[a,{called:n,loading:r,data:l}]=ue(_a,{variables:{mediaId:e.id}});let s=[];n?r||(s=l&&l.media.downloads||[]):e.downloads?s=e.downloads:a();const i=s.map(c=>({title:c.title,url:c.mediaUrl.url,width:c.mediaUrl.width,height:c.mediaUrl.height,fileSize:c.mediaUrl.fileSize}));return f(J,{children:[o(ee,{children:t("sidebar.download.title","Download")}),o($a,{rows:i})]})},Je=({name:e,value:t})=>f("div",{children:[o("div",{className:"inline-block w-[100px] font-semibold text-sm text-[#888] text-right mr-2",children:e}),o("div",{className:"inline-block text-base",children:t})]}),et=y.div`
  margin-bottom: 1.5rem;
`,Ea=({media:e})=>{const{t}=S();let a=[];const n=Ma(t);if(e!=null&&e.exif){const l=e==null?void 0:e.exif,s=Object.keys(n).filter(h=>l[h]!==null&&h!="__typename"),i=s.reduce((h,w)=>{const x=l[w];return N(x)?h:{...h,[w]:x}},{});N(i.dateShot)||(i.dateShot=new Date(i.dateShot).toLocaleString()),typeof i.exposure=="number"&&i.exposure!==0&&(i.exposure=`1/${Math.round(1/i.exposure)}`);const c=e.exif.coordinates;N(c)||(i.coordinates=`${Math.round(c.latitude*1e6)/1e6}, ${Math.round(c.longitude*1e6)/1e6}`);const u=Pa(t);typeof i.exposureProgram=="number"&&u[i.exposureProgram]?i.exposureProgram=u[i.exposureProgram]:i.exposureProgram!==0&&delete i.exposureProgram,N(i.aperture)||(i.aperture=`f/${i.aperture}`),N(i.focalLength)||(i.focalLength=`${i.focalLength}mm`);const p=Fa(t);typeof i.flash=="number"&&p[i.flash]&&(i.flash=p[i.flash]),a=s.map(h=>o(Je,{name:n[h],value:i[h]},h))}let r=[];if(e!=null&&e.videoMetadata){const l=e.videoMetadata;let s=Object.keys(l).filter(i=>!["id","__typename","width","height"].includes(i)).reduce((i,c)=>{const u=l[c];return N(u)?i:{...i,[c]:u}},{});s={dimensions:`${e.videoMetadata.width}x${e.videoMetadata.height}`,...s},r=Object.keys(s).map(i=>o(Je,{name:i,value:s[i]},i))}return f("div",{children:[o(et,{children:r}),o(et,{children:a})]})},Ma=e=>({description:e("sidebar.media.exif.description","Description"),camera:e("sidebar.media.exif.name.camera","Camera"),maker:e("sidebar.media.exif.name.maker","Maker"),lens:e("sidebar.media.exif.name.lens","Lens"),exposureProgram:e("sidebar.media.exif.name.exposure_program","Program"),dateShot:e("sidebar.media.exif.name.date_shot","Date shot"),exposure:e("sidebar.media.exif.name.exposure","Exposure"),aperture:e("sidebar.media.exif.name.aperture","Aperture"),iso:e("sidebar.media.exif.name.iso","ISO"),focalLength:e("sidebar.media.exif.name.focal_length","Focal length"),flash:e("sidebar.media.exif.name.flash","Flash"),coordinates:e("sidebar.media.exif.name.coordinates","Coordinates")}),Pa=e=>({0:e("sidebar.media.exif.exposure_program.not_defined","Not defined"),1:e("sidebar.media.exif.exposure_program.manual","Manual"),2:e("sidebar.media.exif.exposure_program.normal_program","Normal program"),3:e("sidebar.media.exif.exposure_program.aperture_priority","Aperture priority"),4:e("sidebar.media.exif.exposure_program.shutter_priority","Shutter priority"),5:e("sidebar.media.exif.exposure_program.creative_program","Creative program"),6:e("sidebar.media.exif.exposure_program.action_program","Action program"),7:e("sidebar.media.exif.exposure_program.portrait_mode","Portrait mode"),8:e("sidebar.media.exif.exposure_program.landscape_mode","Landscape mode"),9:e("sidebar.media.exif.exposure_program.bulb","Bulb")}),Fa=e=>{const t={no_flash:e("sidebar.media.exif.flash.no_flash","No Flash"),fired:e("sidebar.media.exif.flash.fired","Fired"),did_not_fire:e("sidebar.media.exif.flash.did_not_fire","Did not fire"),on:e("sidebar.media.exif.flash.on","On"),off:e("sidebar.media.exif.flash.off","Off"),auto:e("sidebar.media.exif.flash.auto","Auto"),return_not_detected:e("sidebar.media.exif.flash.return_not_detected","Return not detected"),return_detected:e("sidebar.media.exif.flash.return_detected","Return detected"),no_flash_function:e("sidebar.media.exif.flash.no_flash_function","No flash function"),red_eye_reduction:e("sidebar.media.exif.flash.red_eye_reduction","Red-eye reduction")};return{0:t.no_flash,1:t.fired,5:`${t.fired}, ${t.return_not_detected}`,7:`${t.fired}, ${t.return_detected}`,8:`${t.on}, ${t.did_not_fire}`,9:`${t.on}, ${t.fired}`,13:`${t.on}, ${t.return_not_detected}`,15:`${t.on}, ${t.return_detected}`,16:`${t.off}, ${t.did_not_fire}`,20:`${t.off}, ${t.did_not_fire}, ${t.return_not_detected}`,24:`${t.auto}, ${t.did_not_fire}`,25:`${t.auto}, ${t.fired}`,29:`${t.auto}, ${t.fired}, ${t.return_not_detected}`,31:`${t.auto}, ${t.fired}, ${t.return_detected}`,32:`${t.no_flash_function}`,48:`${t.off}, ${t.no_flash_function}`,65:`${t.fired}, ${t.red_eye_reduction}`,69:`${t.fired}, ${t.red_eye_reduction}, ${t.return_not_detected}`,71:`${t.fired}, ${t.red_eye_reduction}, ${t.return_detected}`,73:`${t.on}, ${t.red_eye_reduction}`,77:`${t.on}, ${t.red_eye_reduction}, ${t.return_not_detected}`,79:`${t.on}, ${t.red_eye_reduction}, ${t.return_detected}`,80:`${t.off}, ${t.red_eye_reduction}`,88:`${t.auto}, ${t.did_not_fire}, ${t.red_eye_reduction}`,89:`${t.auto}, ${t.fired}, ${t.red_eye_reduction}`,93:`${t.auto}, ${t.red_eye_reduction}, ${t.return_not_detected}`,95:`${t.auto}, ${t.red_eye_reduction}, ${t.return_detected}`}},wt=y(({origin:e,selectable:t,scale:a,...n})=>o(Ce,{...n}))`
  position: absolute;
  transform-origin: ${({origin:e})=>`${e.x*100}% ${e.y*100}%`};
  object-fit: cover;

  transition: transform 250ms ease-out;
`,Na=y(wt)`
  width: 100%;
  top: 50%;
  transform: translateY(-50%)
    ${({origin:e,scale:t})=>`translate(${(.5-e.x)*100}%, ${(.5-e.y)*100}%) scale(${Math.max(t*.8,1)})`};

  ${({selectable:e,origin:t,scale:a})=>e?`
    &:hover {
      transform: translateY(-50%) translate(${(.5-t.x)*100}%, ${(.5-t.y)*100}%) scale(${Math.max(a*.85,1)})
      `:""}
`,tt=y(wt)`
  height: 100%;
  left: 50%;
  transform: translateX(-50%)
    ${({origin:e,scale:t})=>`translate(${(.5-e.x)*100}%, ${(.5-e.y)*100}%) scale(${Math.max(t*.8,1)})`};

  ${({selectable:e,origin:t,scale:a})=>e?`
    &:hover {
      transform: translateX(-50%) translate(${(.5-t.x)*100}%, ${(.5-t.y)*100}%) scale(${Math.max(a*.85,1)})
      `:""}
`,Da=y.div`
  background-color: #eee;
  position: relative;
  border-radius: 50%;
  width: ${({size:e})=>e};
  height: ${({size:e})=>e};
  object-fit: fill;
  overflow: hidden;
`,qe=({imageFace:e,selectable:t,size:a="150px"})=>{var i;if(!e)return null;const n=e.rectangle,r=Math.min(1/(n.maxX-n.minX),1/(n.maxY-n.minY)),l={x:(n.minX+n.maxX)/2,y:(n.minY+n.maxY)/2};let s=tt;return e.media.thumbnail&&(s=e.media.thumbnail.width>e.media.thumbnail.height?tt:Na),o(Da,{size:a,children:o(s,{selectable:t,scale:r,origin:l,src:(i=e.media.thumbnail)==null?void 0:i.url})})},Ra=e=>d.exports.createElement("svg",{width:"12px",height:"3px",viewBox:"0 0 8 2",xmlns:"http://www.w3.org/2000/svg",xmlnsXlink:"http://www.w3.org/1999/xlink",...e},d.exports.createElement("path",{d:"M1,0 C1.55228475,0 2,0.44771525 2,1 C2,1.55228475 1.55228475,2 1,2 C0.44771525,2 0,1.55228475 0,1 C0,0.44771525 0.44771525,0 1,0 Z M4,0 C4.55228475,0 5,0.44771525 5,1 C5,1.55228475 4.55228475,2 4,2 C3.44771525,2 3,1.55228475 3,1 C3,0.44771525 3.44771525,0 4,0 Z M7,0 C7.55228475,0 8,0.44771525 8,1 C8,1.55228475 7.55228475,2 7,2 C6.44771525,2 6,1.55228475 6,1 C6,0.44771525 6.44771525,0 7,0 Z",fill:"currentColor"})),Ct=({active:e,text:t})=>o(Pt,{active:!0,message:t,className:e?"opacity-100":"opacity-0"}),yt=({loading:e,fetchMore:t,data:a,getItems:n})=>{const r=d.exports.useRef(null),l=d.exports.useRef(null),[s,i]=d.exports.useState(!1),c=()=>{var h;const p={root:null,rootMargin:"-100% 0px 0px 0px",threshold:0};(h=r.current)==null||h.disconnect(),!s&&(r.current=new IntersectionObserver(w=>{if(w.find(x=>x.isIntersecting==!1)){const x=a!==void 0?n(a).length:0;t({variables:{offset:x}}).then(g=>{n(g.data).length==0&&i(!0)})}},p),l.current&&!e&&r.current.observe(l.current))},u=d.exports.useCallback(p=>{l.current=p,r.current!=null&&r.current.disconnect(),p!=null&&c()},[]);return d.exports.useEffect(()=>{r.current&&l.current&&(e?r.current.unobserve(l.current):r.current.observe(l.current))},[e]),d.exports.useEffect(()=>{c()},[t,a,s]),d.exports.useEffect(()=>{i(!1)},[a]),{containerElem:u,finished:s}},Ta=y(Ce)`
  max-width: 120px;
  max-height: 80px;
`,Aa=({imageFace:e,faceSelected:t,setFaceSelected:a})=>{var n;return f(ie,{children:[o(Fe,{children:o(Ta,{src:(n=e.media.thumbnail)==null?void 0:n.url,onClick:a})}),o(Fe,{className:"min-w-64 w-full",children:o(Be,{label:e.media.title,checked:t,onChange:a})})]},e.id)},_t=({imageFaces:e,selectedImageFaces:t,setSelectedImageFaces:a,title:n})=>{const{t:r}=S(),[l,s]=d.exports.useState(""),i=e.filter(c=>l==""||c.media.title.toLowerCase().includes(l.toLowerCase())).map(c=>o(Aa,{imageFace:c,faceSelected:t.includes(c),setFaceSelected:()=>a(u=>u.includes(c)?u.filter(p=>p!=c):[...u,c])},c.id));return f(G,{children:[o(ge,{className:"w-full",children:f(dt,{children:[o(ie,{children:o(ve,{colSpan:2,children:n})}),o(ie,{children:o(ve,{colSpan:2,children:o(se,{value:l,onChange:c=>s(c.target.value),placeholder:r("people_page.tableselect_image_faces.search_images_placeholder","Search images..."),fullWidth:!0})})})]})}),o("div",{className:"overflow-auto max-h-[500px] mt-2",children:o(ge,{children:o(ut,{children:i})})})]})},Oa=M`
  mutation detachImageFaces($faceIDs: [ID!]!) {
    detachImageFaces(imageFaceIDs: $faceIDs) {
      id
      label
    }
  }
`,It=e=>{const[t]=T(Oa,e);return async a=>{const n=a.map(l=>l.id);return await t({variables:{faceIDs:n}})}},Ba=({open:e,setOpen:t,faceGroup:a,selectedImageFaces:n})=>{var h;const{t:r}=S(),[l,s]=d.exports.useState([]),i=ye(),c=It({refetchQueries:[{query:K}]}),u=()=>{c(l).then(({data:w})=>{if(N(w))throw new Error("Expected data not to be null");t(!1),i(`/people/${w.detachImageFaces.id}`)})};if(d.exports.useEffect(()=>{N(n)||s(n)},[n]),d.exports.useEffect(()=>{e||s([])},[e]),e==!1)return null;const p=(h=a==null?void 0:a.imageFaces)!=null?h:[];return o(Ue,{title:r("people_page.modal.detach_image_faces.title","Detach Image Faces"),description:r("people_page.modal.detach_image_faces.description","Detach selected images of this face group and move them to a new face groups"),actions:[{key:"cancel",label:r("general.action.cancel","Cancel"),onClick:()=>t(!1)},{key:"detach",label:r("people_page.modal.detach_image_faces.action.detach","Detach image faces"),variant:"positive",onClick:()=>u()}],onClose:()=>t(!1),open:e,children:o(_t,{imageFaces:p,selectedImageFaces:l,setSelectedImageFaces:s,title:r("people_page.modal.detach_image_faces.action.select_images","Select images to detach")})})},Ua=y.div`
  display: inline-block;
  border-radius: 50%;
  border: 2px solid
    ${({$selected:e})=>e?"#2185c9":"rgba(255,255,255,0)"};
`,Ga=y(Fe)`
  display: flex;
  align-items: center;
`;y.span`
  ${({$selected:e})=>e&&"font-weight: bold;"}
  margin-left: 12px;
  width: 100%;
`;const za=({faceGroup:e,faceSelected:t,setFaceSelected:a})=>{var n;return o(ie,{onClick:a,children:f(Ga,{children:[o(Ua,{$selected:t,children:o(qe,{imageFace:e.imageFaces[0],size:"50px",selectable:!1})}),o("span",{className:`ml-3 ${t?"font-semibold":""} ${e.label?"":"text-gray-500 italic"}`,children:(n=e.label)!=null?n:"Unlabeled"})]})},e.id)},St=({faceGroups:e,selectedFaceGroup:t,setSelectedFaceGroup:a,title:n})=>{const{t:r}=S(),[l,s]=d.exports.useState(""),i=e.filter(c=>l==""||c.label&&c.label.toLowerCase().includes(l.toLowerCase())).map(c=>o(za,{faceGroup:c,faceSelected:(t==null?void 0:t.id)==c.id,setFaceSelected:()=>a(c)},c.id));return f(G,{children:[o(ge,{className:"w-full",children:f(dt,{children:[o(ie,{children:o(ve,{children:n})}),o(ie,{children:o(ve,{children:o(se,{fullWidth:!0,value:l,onChange:c=>s(c.target.value),placeholder:r("people_page.tableselect_face_group.search_faces_placeholder","Search faces...")})})})]})}),o("div",{className:"overflow-auto max-h-[500px] mt-2",children:o(ge,{className:"w-full",children:o(ut,{children:i})})})]})},Ha=M`
  mutation moveImageFaces($faceIDs: [ID!]!, $destFaceGroupID: ID!) {
    moveImageFaces(
      imageFaceIDs: $faceIDs
      destinationFaceGroupID: $destFaceGroupID
    ) {
      id
      imageFaces {
        id
      }
    }
  }
`,kt=({open:e,setOpen:t,faceGroup:a,preselectedImageFaces:n})=>{const{t:r}=S(),[l,s]=d.exports.useState([]),[i,c]=d.exports.useState(null),[u,p]=d.exports.useState(!1),h=ye(),[w]=T(Ha,{refetchQueries:[{query:K}]}),[x,{data:g}]=ue(K);if(d.exports.useEffect(()=>{N(n)||(s(n),p(!0))},[n]),d.exports.useEffect(()=>{u&&x()},[u]),d.exports.useEffect(()=>{e||(p(!1),s([]),c(null))},[e]),e==!1)return null;const b=()=>{const P=l.map(k=>k.id);if(N(i))throw new Error("Expected selectedFaceGroup not to be null");w({variables:{faceIDs:P,destFaceGroupID:i.id}}).then(()=>{t(!1),h(`/people/${i.id}`)})},C=a.imageFaces;let m=null;if(!u)m=o(_t,{imageFaces:C,selectedImageFaces:l,setSelectedImageFaces:s,title:r("people_page.modal.move_image_faces.image_select_table.title","Select images to move")});else if(g&&a){const P=g.myFaceGroups.filter(k=>k.id!=a.id);m=o(St,{title:r("people_page.modal.move_image_faces.destination_face_group_table.title","Select destination face group"),faceGroups:P,selectedFaceGroup:i,setSelectedFaceGroup:c})}else m=o("div",{children:r("general.loading.default","Loading...")});let I;return u?I={key:"move",label:r("people_page.modal.move_image_faces.destination_face_group_table.move_action","Move image faces"),onClick:()=>b(),variant:"positive"}:I={key:"next",label:r("people_page.modal.move_image_faces.image_select_table.next_action","Next"),onClick:()=>p(!0),variant:"positive"},o(Ue,{title:r("people_page.modal.move_image_faces.title","Move Image Faces"),description:r("people_page.modal.move_image_faces.description","Move selected images of this face group to another face group"),onClose:()=>t(!1),open:e,actions:[{key:"cancel",label:r("general.action.cancel","Cancel"),onClick:()=>t(!1)},I],children:m})},Wa=({faceGroup:e})=>{var P,k;const{t}=S(),[a,n]=d.exports.useState(!1),[r,l]=d.exports.useState((P=e==null?void 0:e.label)!=null?P:""),s=d.exports.createRef(),[i,c]=d.exports.useState(!1),[u,p]=d.exports.useState(!1),[h,w]=d.exports.useState(!1),[x,{loading:g}]=T(Ye),b=()=>{var _;l((_=e==null?void 0:e.label)!=null?_:""),n(!1)};d.exports.useEffect(()=>{s.current&&s.current.focus()},[s]),d.exports.useEffect(()=>{g||b()},[g]);const C=_=>{if(_.key=="Escape"){b();return}};let m;a?m=o(G,{children:o(se,{loading:g,ref:s,placeholder:t("people_page.face_group.label_placeholder","Label"),action:()=>{if(N(e))throw new Error("Expected faceGroup to be defined");x({variables:{groupID:e.id,label:r||null}})},value:r,onKeyDown:C,onChange:_=>l(_.target.value),onBlur:()=>{b()}})}):m=o(G,{children:o("h1",{className:`text-2xl font-semibold ${e!=null&&e.label?"":"text-gray-600 dark:text-gray-400"}`,children:(k=e==null?void 0:e.label)!=null?k:t("people_page.face_group.unlabeled_person","Unlabeled person")})});let I=null;return e&&(I=f(G,{children:[o(Lt,{open:i,setOpen:c,sourceFaceGroup:e,refetchQueries:[{query:K}]}),o(kt,{open:u,setOpen:p,faceGroup:e}),o(Ba,{open:h,setOpen:w,faceGroup:e})]})),f(G,{children:[f("div",{children:[o("div",{className:"mb-2",children:m}),f("ul",{className:"flex gap-2 flex-wrap mb-6",children:[o("li",{children:o(te,{onClick:()=>n(!0),children:t("people_page.action_label.change_label","Change label")})}),o("li",{children:o(te,{onClick:()=>c(!0),children:t("people_page.action_label.merge_people","Merge people")})}),o("li",{children:o(te,{onClick:()=>w(!0),children:t("people_page.action_label.detach_images","Detach images")})}),o("li",{children:o(te,{onClick:()=>p(!0),children:t("people_page.action_label.move_faces","Move faces")})})]})]}),I]})},Qa=M`
  query singleFaceGroup($id: ID!, $limit: Int!, $offset: Int!) {
    faceGroup(id: $id) {
      id
      label
      imageFaces(paginate: { limit: $limit, offset: $offset }) {
        id
        rectangle {
          minX
          maxX
          minY
          maxY
        }
        media {
          id
          type
          title
          blurhash
          thumbnail {
            url
            width
            height
          }
          highRes {
            url
          }
          favorite
        }
      }
    }
  }
`,Xa=({faceGroupID:e})=>{const{t}=S(),{data:a,error:n,loading:r,fetchMore:l}=de(Qa,{variables:{limit:200,offset:0,id:e}}),[s,i]=d.exports.useReducer(Rt,{presenting:!1,activeIndex:-1,media:[]}),{containerElem:c,finished:u}=yt({loading:r,fetchMore:l,data:a,getItems:h=>h.faceGroup.imageFaces});d.exports.useEffect(()=>{const h=(a==null?void 0:a.faceGroup.imageFaces.map(w=>w.media))||[];i({type:"replaceMedia",media:h})},[a]);const p=a==null?void 0:a.faceGroup;return n?o("div",{children:n.message}):f("div",{ref:c,children:[o(Wa,{faceGroup:p}),f("div",{children:[o(uo,{loading:r,dispatchMedia:i,mediaState:s}),o(Ct,{active:!u&&!r,text:t("general.loading.paginate.media","Loading more media")})]})]})},K=M`
  query myFaces($limit: Int, $offset: Int) {
    myFaceGroups(paginate: { limit: $limit, offset: $offset }) {
      id
      label
      imageFaceCount
      imageFaces(paginate: { limit: 1 }) {
        id
        rectangle {
          minX
          maxX
          minY
          maxY
        }
        media {
          id
          title
          thumbnail {
            url
            width
            height
          }
        }
      }
    }
  }
`,Ye=M`
  mutation setGroupLabel($groupID: ID!, $label: String) {
    setFaceGroupLabel(faceGroupID: $groupID, label: $label) {
      id
      label
    }
  }
`,qa=M`
  mutation recognizeUnlabeledFaces {
    recognizeUnlabeledFaces {
      id
    }
  }
`,Ya=({labeled:e,children:t,className:a,...n})=>o("span",{...n,className:Nt(a,`${e?"":"text-gray-400 dark:text-gray-500"}`),children:t}),rt=y(Ya)`
  &:hover,
  &:focus-visible {
    color: #2683ca;
  }
`,Ze=({group:e,className:t,textFieldClassName:a,editLabel:n,setEditLabel:r})=>{var g,b;const{t:l}=S(),[s,i]=d.exports.useState((g=e.label)!=null?g:""),c=d.exports.createRef(),[u,{loading:p}]=T(Ye,{variables:{groupID:e.id}}),h=()=>{var C;i((C=e.label)!=null?C:""),r(!1)};d.exports.useEffect(()=>{var C;(C=c.current)==null||C.focus()},[c]),d.exports.useEffect(()=>{p||h()},[p]);const w=C=>{if(C.key=="Escape"){h();return}};let x;return n?x=o(rt,{className:t,labeled:!!e.label,children:o(se,{className:a,loading:p,ref:c,placeholder:l("people_page.face_group.label_placeholder","Label"),value:s,action:()=>u({variables:{groupID:e.id,label:s==""?null:s}}),onKeyDown:w,onChange:C=>i(C.target.value),onBlur:()=>{h()}})}):x=f(rt,{className:we(t,"whitespace-nowrap inline-block overflow-hidden overflow-clip"),labeled:!!e.label,onClick:()=>r(!0),children:[o(Za,{children:e.imageFaceCount}),o("button",{className:"",children:(b=e.label)!=null?b:l("people_page.face_group.unlabeled","Unlabeled")})]}),x},Za=y.span.attrs({className:"bg-gray-100 text-gray-900 dark:bg-gray-400 dark:text-black text-sm px-1 mr-2 rounded-md"})``,$t=({group:e})=>{const t=e.imageFaces[0],[a,n]=d.exports.useState(!1);return f("div",{className:"m-3",children:[o(ne,{to:`/people/${e.id}`,children:o(qe,{imageFace:t,selectable:!0})}),o(Ze,{className:"block cursor-pointer text-center w-full mt-3",textFieldClassName:"w-[140px]",group:e,editLabel:a,setEditLabel:n})]})},Va=y.div`
  display: flex;
  flex-wrap: wrap;
  margin-top: 24px;
`,ja=()=>{const{t:e}=S(),{data:t,error:a,loading:n,fetchMore:r}=de(K,{variables:{limit:50,offset:0}}),[l,{loading:s}]=T(qa),{containerElem:i,finished:c}=yt({loading:n,fetchMore:r,data:t,getItems:p=>p.myFaceGroups});if(a)return o("div",{children:a.message});let u=null;return t&&(u=t.myFaceGroups.map(p=>o($t,{group:p},p.id))),f(at,{title:e("title.people","People"),children:[o(te,{disabled:s,onClick:()=>{l()},children:e("people_page.recognize_unlabeled_faces_button","Recognize unlabeled faces")}),o(Va,{ref:i,children:u}),o(Ct,{active:!c&&!n,text:e("general.loading.paginate.faces","Loading more people")})]})},Ka=()=>{const{t:e}=S(),{person:t}=Ft();if(N(t))throw new Error("Expected `person` parameter to be defined");return o(at,{title:e("title.people","People"),children:o(Xa,{faceGroupID:t})})};var go=Object.freeze({__proto__:null,[Symbol.toStringTag]:"Module",MY_FACES_QUERY:K,SET_GROUP_LABEL_MUTATION:Ye,FaceDetails:Ze,FaceGroup:$t,PeoplePage:ja,PersonPage:Ka});const Ja=M`
  mutation combineFaces($destID: ID!, $srcID: ID!) {
    combineFaceGroups(
      destinationFaceGroupID: $destID
      sourceFaceGroupID: $srcID
    ) {
      id
    }
  }
`,Lt=({open:e,setOpen:t,sourceFaceGroup:a,refetchQueries:n})=>{var w;const{t:r}=S(),[l,s]=d.exports.useState(null),i=ye(),{data:c}=de(K),[u]=T(Ja,{refetchQueries:n});if(e==!1)return null;const p=(w=c==null?void 0:c.myFaceGroups.filter(x=>x.id!=(a==null?void 0:a.id)))!=null?w:[],h=()=>{if(N(l))throw new Error("No selected face group");u({variables:{srcID:a.id,destID:l.id}}).then(()=>{t(!1),i(`/people/${l.id}`)})};return o(Ue,{title:r("people_page.modal.merge_face_groups.title","Merge Face Groups"),description:r("people_page.modal.merge_face_groups.description","All images within this face group will be merged into the selected face group."),actions:[{key:"cancel",label:r("general.action.cancel","Cancel"),onClick:()=>t(!1)},{key:"merge",label:r("people_page.modal.action.merge","Merge"),onClick:()=>h(),variant:"positive"}],onClose:()=>t(!1),open:e,children:o(St,{title:r("people_page.modal.merge_face_groups.destination_table.title","Select the destination face"),faceGroups:p,selectedFaceGroup:l,setSelectedFaceGroup:s})})},he=({label:e,className:t,onClick:a})=>o(fe.Item,{children:({active:n})=>o("button",{onClick:a,className:we(`whitespace-normal w-full block py-1 cursor-pointer ${n?"bg-gray-50 text-black":"text-gray-700"}`,t),children:e})}),eo=({menuFlipped:e,face:t,setChangeLabel:a,className:n})=>{const{t:r}=S(),[l,s]=d.exports.useState(!1),[i,c]=d.exports.useState(!1),u=[{query:Et,variables:{id:t.media.id}}],p=ye(),h=It({refetchQueries:u}),w=f(G,{children:[o(Lt,{sourceFaceGroup:t.faceGroup,open:l,setOpen:s,refetchQueries:u}),o(kt,{faceGroup:{imageFaces:[],...t.faceGroup},open:i,setOpen:c,preselectedImageFaces:[t]})]}),x=()=>{!confirm(r("sidebar.people.confirm_image_detach","Are you sure you want to detach this image?"))||h([t]).then(({data:g})=>{if(N(g))throw new Error("Expected data not to be null");p(`/people/${g.detachImageFaces.id}`)})};return f(G,{children:[f(fe,{as:"div",className:we("relative inline-block",n),children:[o(fe.Button,{as:te,className:"px-1.5 py-1.5 align-middle ml-1",children:o(Ra,{className:"text-gray-500"})}),o(fe.Items,{className:"",children:f(ft,{width:120,flipped:e,children:[o(he,{onClick:()=>a(!0),className:"border-b",label:r("people_page.action_label.change_label","Change label")}),o(he,{onClick:()=>s(!0),className:"border-b",label:r("sidebar.people.action_label.merge_face","Merge face")}),o(he,{onClick:()=>x(),className:"border-b",label:r("sidebar.people.action_label.detach_image","Detach image")}),o(he,{onClick:()=>c(!0),label:r("sidebar.people.action_label.move_face","Move face")})]})})]}),w]})},to=({face:e,menuFlipped:t})=>{const[a,n]=d.exports.useState(!1);return f("li",{className:"inline-block",children:[o(ne,{to:`/people/${e.faceGroup.id}`,children:o(qe,{imageFace:e,selectable:!0,size:"92px"})}),f("div",{className:"mt-1 whitespace-nowrap",children:[o(Ze,{className:"text-sm max-w-[80px] align-middle",textFieldClassName:"w-[100px]",group:e.faceGroup,editLabel:a,setEditLabel:n}),!a&&o(eo,{menuFlipped:t,className:"pl-0.5",face:e,setChangeLabel:n})]})]})},ro=({media:e})=>{var n;const{t}=S(),a=((n=e.faces)!=null?n:[]).map((r,l)=>o(to,{face:r,menuFlipped:l==0},r.id));return a.length==0?null:f(J,{children:[o(ee,{children:t("sidebar.people.title","People")}),f("div",{className:"overflow-x-auto mb-[-200px]",style:{scrollbarWidth:"none"},children:[o("ul",{className:"flex gap-4 mx-4",children:a}),o("div",{className:"h-[200px]"})]})]})},ao=({coordinates:e})=>{const{t}=S(),{mapContainer:a,mapboxToken:n}=Tt({mapboxOptions:{interactive:!1,zoom:12,center:{lat:e.latitude,lng:e.longitude}},configureMapbox:(r,l)=>{r.addControl(new l.NavigationControl({showCompass:!1}));const s=new l.Marker({color:"red",scale:.8});s.setLngLat({lat:e.latitude,lng:e.longitude}),s.addTo(r)}});return N(n)?null:f(J,{children:[o(ee,{children:t("sidebar.location.title","Location")}),o("div",{className:"w-full h-64",children:a})]})},Et=M`
  query sidebarMediaQuery($id: ID!) {
    media(id: $id) {
      id
      title
      type
      highRes {
        url
        width
        height
      }
      thumbnail {
        url
        width
        height
      }
      videoWeb {
        url
        width
        height
      }
      videoMetadata {
        id
        width
        height
        duration
        codec
        framerate
        bitrate
        colorProfile
        audio
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
      album {
        id
        title
        path {
          id
          title
        }
      }
      faces {
        id
        rectangle {
          minX
          maxX
          minY
          maxY
        }
        faceGroup {
          id
          label
          imageFaceCount
        }
        media {
          id
          title
          thumbnail {
            url
            width
            height
          }
        }
      }
    }
  }
`,oo=y(Ce)`
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  object-fit: contain;
`,no=y(Dt)`
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
`,lo=({media:e,previewImage:t})=>e.type===be.Photo?o(oo,{src:t==null?void 0:t.url}):e.type===be.Video?o(no,{media:e}):f("div",{children:["ERROR: Unknown media type: ",e.type]}),Ee=({media:e,hidePreview:t})=>{var p,h,w,x;const{updateSidebar:a}=d.exports.useContext(xe),{t:n}=S();let r=null;e.highRes?r=e.highRes:e.thumbnail&&(r=e.thumbnail);const l=(r==null?void 0:r.width)&&(r==null?void 0:r.height)?r.height/r.width:3/2;let s=null;const i=(p=e.exif)==null?void 0:p.coordinates;i&&(s=o(ao,{coordinates:i}));let c=null;const u=e.album;if(!N(u)){console.log("PATH reversed",(h=u.path)!=null?h:[]);const g=[...[...(w=u.path)!=null?w:[]].reverse(),u].map(b=>o("li",{className:"inline-block hover:underline",children:o(ne,{className:"text-blue-900 dark:text-blue-200 hover:underline",to:`/album/${b.id}`,onClick:()=>a(null),children:b.title})},b.id));c=f("div",{className:"mx-4 my-4",children:[o("h2",{className:"uppercase text-xs text-gray-900 dark:text-gray-300 font-semibold",children:n("sidebar.media.album_path","Album path")}),o(vt,{hideLastArrow:!0,children:g})]})}return f("div",{children:[o(gt,{title:(x=e.title)!=null?x:"Loading..."}),o("div",{className:"lg:mx-4",children:!t&&f("div",{className:"w-full h-0 relative",style:{paddingTop:`${Math.min(l,.75)*100}%`},children:[o(lo,{previewImage:r||void 0,media:e}),o(ya,{media:e})]})}),o(Ea,{media:e}),c,o(ro,{media:e}),s,o(La,{media:e}),o(zr,{id:e.id}),o("div",{className:"mt-8",children:o(Yr,{cover_id:e.id})})]})},io=({media:e,hidePreview:t})=>{const[a,{loading:n,error:r,data:l}]=ue(Et);return d.exports.useEffect(()=>{e!=null&&V()&&a({variables:{id:e.id}})},[e]),e?V()?r?o("div",{children:r.message}):n||l==null?o(Ee,{media:e,hidePreview:t}):o(Ee,{media:l.media,hidePreview:t}):o(Ee,{media:e,hidePreview:t}):null},so=y.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  min-height: 200px;
  position: relative;
  margin: -4px;

  @media (max-width: 1000px) {
    /* Compensate for tab bar on mobile */
    margin-bottom: 76px;
  }
`,co=y.div`
  height: 200px;
  flex-grow: 999999;
`,vo=M`
  fragment MediaGalleryFields on Media {
    id
    type
    blurhash
    thumbnail {
      url
      width
      height
    }
    highRes {
      url
    }
    videoWeb {
      url
    }
    favorite
  }
`,uo=({mediaState:e,dispatchMedia:t})=>{const[a]=ga(),{media:n,activeIndex:r,presenting:l}=e,{updateSidebar:s}=d.exports.useContext(xe);let i=[];if(n)i=n.map((c,u)=>{const p=r==u;return o(ha,{media:c,active:p,selectImage:()=>{t({type:"selectImage",index:u}),s(o(io,{media:e.media[u]}))},clickFavorite:()=>{va({media:c,markFavorite:a})},clickPresent:()=>{At({dispatchMedia:t,activeIndex:u})}},c.id)});else for(let c=0;c<6;c++)i.push(o(fa,{},c));return f(G,{children:[f(so,{"data-testid":"photo-gallery-wrapper",children:[i,o(co,{})]}),l&&o(Ot,{activeMedia:e.media[e.activeIndex],dispatchMedia:t})]})};export{bo as A,vo as M,Ct as P,uo as a,ga as b,ha as c,io as d,co as e,go as f,va as t,yt as u};
