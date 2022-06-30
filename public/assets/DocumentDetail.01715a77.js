import{w as C,v as M}from"./element-plus.274d913c.js";import{b as L,B as j,D as A,V as B,W as f,X as N,Y as T,Z as y,$ as P,a0 as z}from"./index.f06094a1.js";import{m as q,u as H}from"./useHighlight.6bca46f2.js";import{r as O,bu as V,aT as Z,t as R,p as b,bv as u,bw as k,bx as p,by as h,b7 as U,R as D,f as a,h as s,$ as _,Z as d,X as m,T as v,S as G,o as c,ai as g}from"./vendor.cec8f34c.js";var W="/assets/icon-empty.9fe6949e.png";function S(t,e){document.querySelectorAll('[data-pid="'+t+'"]').forEach(function(i){let n=i.querySelector("i.editor-arrow-right");n&&!u(n,"expand")&&p(n,"expand");let o=k(i,"data-id");i.style.display=e?null:"none",o&&S(o,e)})}const X={watch:{"$route.params.node_id":function(){this.getDocumentDetail()}},setup(){const t=O(null),{top:e}=V(t),i=j(),{isGuest:n}=Z(i);R(e,()=>T.emit(y,e.value<25),{immediate:!0});const{initHighlight:o}=H(),r=b("documentImportModal"),l=b("setDocumentTitle");return{isGuest:n,title:t,initHighlight:o,documentImportModal:r,setDocumentTitle:l}},data(){return{zoomTemplate:`<template id="template-zoom-image">
                            <div class="zoom-image-wrapper">
                                <div class="zoom-image-container" data-zoom-container></div>
                            </div>
                          </template>`,zoomImageOption:{template:"#template-zoom-image",container:"[data-zoom-container]"},hasDocument:!0,isLoading:!0,DOCUMENT_TYPES:A,document:{},project_id:null}},methods:{onEditBtnClick(){this.$router.push({name:"document.api.edit",params:{project_id:this.$route.params.project_id,node_id:this.$route.params.node_id}})},initTableToggle(){document.querySelectorAll(".ac-param-table .editor-arrow-right").forEach(function(t){t.onclick=function(){S(k(this,"data-id"),!u(this,"expand")),p(this,"expand")}}),document.querySelectorAll("div.collapse-title .response_body_title").forEach(function(t){t.onclick=function(){let e=this.parentElement,i=e.parentElement,n=u(i,"close");h(e.nextElementSibling,n),h(i.nextElementSibling,n),p(i,"close")}}),document.querySelectorAll("h3.collapse-title >span").forEach(function(t){t.onclick=function(){let e=this.parentElement,i=u(e,"close");h(e.nextElementSibling,i),p(e,"close")}})},initMediumZoom(){q(".ProseMirror .image-view img",this.zoomImageOption)},initTippy(){U("[data-tippy-content]",{theme:"light",appendTo:document.querySelector(".ProseMirror")})},initCodeBlockToClipboard(){document.querySelectorAll(".code-block button").forEach(t=>{t.setAttribute("data-text",t.parentElement.querySelector("code").innerText)})},initStaticDocInteractive(){this.$nextTick(()=>{this.initTableToggle(),this.initTippy(),this.initMediumZoom(),this.initCodeBlockToClipboard(),this.initHighlight(document.querySelectorAll("pre code"))})},onImportBtnCLick(){this.documentImportModal.show({project_id:this.$route.params.project_id},B)},getDocumentDetail(){const t=parseInt(this.$route.params.node_id,10);if(isNaN(t)){f(),this.isLoading=!1,this.hasDocument=!1;return}this.doc_id=t,this.project_id=parseInt(this.$route.params.project_id,10),this.isLoading=!0,this.hasDocument=!0,N(this.project_id,this.doc_id,"html").then(e=>{this.document=this.transferDoc(e.data||{}),this.setDocumentTitle(this.document.title),this.initStaticDocInteractive()}).catch(e=>{}).finally(()=>{this.$nextTick(()=>{this.isLoading=!1,f()})})},transferDoc(t){return t}},mounted(){T.emit(y,!1),this.getDocumentDetail()},unmounted(){this.document={},this.isLoading=!0,this.hasDocument=!0}},Y={class:"ac-document"},F={class:"ac-document__desc"},J=s("i",{class:"iconfont iconIconPopoverUser"},null,-1),K=s("i",{class:"iconfont icontime"},null,-1),Q=["innerHTML"],$={key:0},tt=s("img",{src:W,alt:""},null,-1),et={key:0,style:{width:"470px",display:"block",margin:"auto"}},ot=g(" \u60A8\u5F53\u524D\u5C1A\u672A\u521B\u5EFA\u6587\u6863\uFF0C\u8BF7\u4ECE\u5DE6\u4FA7\u76EE\u5F55\u680F\u70B9\u51FB\u6DFB\u52A0\uFF0C\u5F00\u59CB\u5728\u7EBF\u7EF4\u62A4 API \u6587\u6863\u3002\u60A8\u8FD8\u53EF\u4EE5\u5C06\u672C\u5730\u9879\u76EE "),it={key:1,style:{width:"470px",display:"block",margin:"auto"}},nt=["innerHTML"];function st(t,e,i,n,o,r){const l=C,x=P,E=z,I=M;return D((c(),a("div",Y,[D(s("div",null,[s("h1",{class:"ac-document__title",ref:"title"},_(o.document.title),513),s("p",F,[d(l,{effect:"dark",content:o.document.last_updated_by+" \u6700\u540E\u7F16\u8F91",placement:"bottom"},{default:m(()=>[s("span",null,[J,g(_(o.document.last_updated_by),1)])]),_:1},8,["content"]),d(l,{effect:"dark",content:"\u66F4\u65B0\u4E8E "+o.document.updated_time,placement:"bottom"},{default:m(()=>[s("span",null,[K,g(_(o.document.updated_time),1)])]),_:1},8,["content"])]),o.document.content?(c(),a("div",{key:0,class:"ProseMirror readonly",innerHTML:o.document.content},null,8,Q)):v("",!0)],512),[[G,o.hasDocument&&o.document.id]]),o.hasDocument?v("",!0):(c(),a("div",$,[d(x,{styles:{width:"260px",height:"auto","margin-bottom":"26px"}},{icon:m(()=>[tt]),title:m(()=>[n.isGuest?(c(),a("div",it,"\u60A8\u5F53\u524D\u5C1A\u672A\u521B\u5EFA\u6587\u6863")):(c(),a("div",et,[ot,s("a",{class:"text-blue-600",href:"javascript:void(0);",onClick:e[0]||(e[0]=(...w)=>r.onImportBtnCLick&&r.onImportBtnCLick(...w))},"\u5BFC\u5165")]))]),_:1})])),d(E,{bottom:100,right:100}),s("div",{innerHTML:o.zoomTemplate},null,8,nt)])),[[I,o.isLoading]])}var dt=L(X,[["render",st]]);export{dt as default};
