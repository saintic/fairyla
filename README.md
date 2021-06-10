<p align="center">
  <a href="https://is.fairyla.com">
    <img src="material/logo.png" alt="fairyla" width="360">
  </a>
</p>
<p align="center">
  Fairyla - 谁还不是小仙女啦！
</p>
<p align="center">
  <a href="https://github.com/staugur/fairyla/actions/workflows/test.yml">
      <img src="https://github.com/staugur/fairyla/actions/workflows/test.yml/badge.svg" alt="Go test">
  </a>
  <a href="https://hub.docker.com/r/staugur/fairyla">
      <img src="https://github.com/staugur/fairyla/actions/workflows/docker.yml/badge.svg" alt="Publish Docker Image">
  </a>
  <a href="https://github.com/staugur/fairyla/actions/workflows/asset.yml">
      <img src="https://github.com/staugur/fairyla/actions/workflows/asset.yml/badge.svg" alt="Publish Release Asset">
  </a>
  <a href="https://codecov.io/gh/staugur/fairyla">
    <img src="https://codecov.io/gh/staugur/fairyla/branch/master/graph/badge.svg?token=FXV9VCEVLP"/>
  </a>
</p>

---

这是一个前后端分离的应用，前端基于Vue3+Element-Plus，后端基于Golang+Echo，
照片/视频存放到[sapic](https://github.com/sapicd/sapic)，数据存放到Redis~~

这个项目来源于一个 `fairy` 女生，起初设计为男女（情侣）双方相对私密相册空间，不过经过功能更新，配置logo、slogan等，可作为花瓣网、堆糖网等精简化小众形态。

## 功能

- 私有或公开的限量专辑
- 以专辑为单位的照片（视频）集合
- 共享专辑（双方共同维护专辑中照片）
