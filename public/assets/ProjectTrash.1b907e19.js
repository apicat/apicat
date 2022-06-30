import{E as C,S as F,g as S,h as x,f as T,i as L,K as B,L as V,H as $,v as I}from"./element-plus.85f938a9.js";import{am as p,o as f,W as D,X as a,Z as i,av as P,ai as g,G as N,aT as R,r as h,t as U,f as K,R as M,h as b,U as A,bf as w,bt as O}from"./vendor.b27f3889.js";import{b as q,I as E,J as H,x as G,K as y,M as J,N as W}from"./index.a22b0d3c.js";const X={emits:["on-ok"],data(){return{project_id:this.$route.params.project_id||"",isShow:!1,isLoading:!1,dir:[],form:{dir_id:[]},rules:{dir_id:{required:!0,min_len:1,message:"\u8BF7\u9009\u62E9\u5206\u7C7B",trigger:"change",type:"array"}}}},watch:{isShow:function(){!this.isShow&&this.reset()},"form.dir_id":function(){let t=this.form.dir_id.slice(-1)[0];t!==void 0&&(this.document.node_id=t)}},methods:{show(t){if(!t)throw new Error("\u6587\u6863\u4FE1\u606F\u4E0D\u80FD\u4E3A\u7A7A\uFF01");this.document=t,this.isShow=!0},hide(){this.isShow=!1},onCloseBtnClick(){this.isShow=!1,this.reset()},reset(){this.$refs.teamForm.resetFields()},handleSubmit(t){this.$refs[t].validate(e=>e&&this.submit())},submit(){this.isLoading=!0,E(this.document).then(t=>{this.onCloseBtnClick(),C({type:"success",closable:!0,message:p("span",null,["\u6587\u6863\u6062\u590D\u6210\u529F\uFF0C",p("a",{class:"text-blue-600",href:`/editor/${this.project_id}/doc/${this.document.doc_id}`},"\u67E5\u770B\u8BE6\u60C5")])}),this.$emit("on-ok")}).finally(()=>{this.isLoading=!1})},transferDir(t){let e=[],u=(c,o)=>{(c||[]).forEach(s=>{let m={value:s.id,label:s.title,children:[]};o.push(m),s.sub_nodes&&s.sub_nodes.length&&u(s.sub_nodes,m.children)})};return u(t,e),[{value:0,label:"\u6839\u76EE\u5F55"}].concat(e)},async getDocumentDirList(t){H(t).then(({data:e})=>{this.dir=this.transferDir(e)})}},mounted(){this.getDocumentDirList(this.project_id)}},Z=g("\u6682\u65E0\u6570\u636E"),z=g(" \u53D6\u6D88 "),Q=g(" \u786E\u5B9A ");function Y(t,e,u,c,o,s){const m=F,_=S,n=x,r=T,l=L;return f(),D(l,{modelValue:o.isShow,"onUpdate:modelValue":e[4]||(e[4]=d=>o.isShow=d),width:400,"custom-class":"show-footer-line vertical-center-modal","close-on-click-modal":!1,"append-to-body":"",title:"\u539F\u6587\u6863\u6240\u5728\u5206\u7C7B\u5DF2\u88AB\u5220\u9664\uFF0C\u8BF7\u9009\u62E9\u5176\u4ED6\u5206\u7C7B"},{footer:a(()=>[i(r,{onClick:e[2]||(e[2]=d=>s.onCloseBtnClick())},{default:a(()=>[z]),_:1}),i(r,{loading:o.isLoading,type:"primary",onClick:e[3]||(e[3]=d=>s.handleSubmit("teamForm"))},{default:a(()=>[Q]),_:1},8,["loading"])]),default:a(()=>[i(n,{ref:"teamForm",model:o.form,rules:o.rules,"label-position":"top",style:{"margin-bottom":"-19px"},onKeyup:e[1]||(e[1]=P(d=>s.handleSubmit("teamForm"),["enter"]))},{default:a(()=>[i(_,{label:"",prop:"dir_id",class:"hide_required"},{default:a(()=>[i(m,{modelValue:o.form.dir_id,"onUpdate:modelValue":e[0]||(e[0]=d=>o.form.dir_id=d),class:"w-full",options:o.dir,props:{checkStrictly:!0},placeholder:"\u8BF7\u9009\u62E9\u5206\u7C7B"},{empty:a(()=>[Z]),_:1},8,["modelValue","options"])]),_:1})]),_:1},8,["model","rules"])]),_:1},8,["modelValue"])}var ee=q(X,[["render",Y]]);const te=b("span",null,"\u56DE\u6536\u7AD9",-1),oe=["href"],se=["onClick"],ne=N({setup(t){const e=G(),{projectInfo:u}=R(e),c=h(),o=h([]),s=h(!1),m=n=>{w.start();const r={project_id:e.projectInfo.id,doc_id:n.id};E(r).then(({status:l})=>{if(l===y.NO_PARENT_DIR){c.value&&c.value.show(r);return}l===y.OK&&(C({type:"success",showClose:!0,message:()=>p("span",null,["\u6587\u6863\u6062\u590D\u6210\u529F\uFF0C",p("a",{class:"text-blue-600",href:`/editor/${e.projectInfo.id}/doc/${n.id}`},"\u67E5\u770B\u8BE6\u60C5")])}),_())}).catch(l=>{}).finally(()=>{w.done()})},_=async()=>{s.value=!0;try{const{data:n}=await J(e.projectInfo.id);o.value=(n||[]).map(r=>(r.remaining=O(r.deleted_at)+"\u5929",r.previewUrl=W({project_id:e.projectInfo.id,doc_id:r.id}),r))}catch{}finally{s.value=!1}};return U(()=>u.value,async()=>{u.value&&u.value.id&&await _()},{immediate:!0}),(n,r)=>{const l=B,d=V,k=$,j=I;return f(),K(A,null,[M((f(),D(k,{shadow:"never","body-style":{padding:0}},{header:a(()=>[te]),default:a(()=>[i(d,{data:o.value,"empty-text":"\u6682\u65E0\u6570\u636E"},{default:a(()=>[i(l,{prop:"title",label:"\u6587\u6863\u540D\u79F0","show-overflow-tooltip":""}),i(l,{prop:"deleted_at",label:"\u5220\u9664\u65F6\u95F4"}),i(l,{prop:"remaining",label:"\u5269\u4F59"}),i(l,{label:"\u64CD\u4F5C"},{default:a(({row:v})=>[b("a",{class:"cursor-pointer mr-3 text-blue-600",target:"_blank",href:v.previewUrl},"\u9884\u89C8",8,oe),b("span",{class:"cursor-pointer mr-3 text-blue-600",href:"javascript:void(0)",onClick:re=>m(v)},"\u6062\u590D",8,se)]),_:1})]),_:1},8,["data"])]),_:1})),[[j,s.value]]),i(ee,{ref_key:"restoreDocumentModal",ref:c,onOnOk:_},null,512)],64)}}});export{ne as default};
