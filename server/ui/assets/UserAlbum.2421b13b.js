import{_ as a}from"./Album.7be08414.js";import{X as e,V as s,c as t}from"./vendor.884d25e1.js";import"./index.7bcebd71.js";const r={Name:"UserAlbum",components:{Album:a},data:()=>({albums:[]}),created(){this.$http.get("/user/album").then((a=>{for(let e of a.data)e.steady_fairy?e.fairy=e.steady_fairy:e.fairy=e.latest_fairy,e.to={name:"UserAlbumFairy",params:{name:e.name}},this.albums.push(e)}))}};r.render=function(a,r,m,u,l,n){const o=e("Album");return s(),t(o,{albums:l.albums},null,8,["albums"])};export default r;
