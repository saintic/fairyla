# fairyla api

Fairyla是前后端分离项目，后端使用Golang（echo）编写，所有接口统一在 `/api` 路径下。

接口均返回 `JSON` 类型（Object），基本结构是：

```json
{"message": "{string}", "success": "{boolean}"}
```

success表明请求是否成功处理，成功时额外返回 `data` 字段，其类型可能是 object 或 array

message是错误消息，当success未成功时有意义。

## Misc（杂项）

### GET+HEAD /ready

k8s准备就绪探针，应用正常连通redis即响应`ok`，状态200，否则返回`err`，状态503，
虽是JSON类型，但无基本结构，仅此一个特例。

### GET /config

前端全局状态、配置

### GET /album

公共专辑列表（Ta是）

### GET /album/:owner/:id

公共专辑详情

- `:owner`：专辑属主
- `:id`：专辑ID或Name

## Auth（认证）

### POST /auth/signup

注册

### POST /auth/signin

登录

### POST /auth/forgot

忘记密码

#### POST /auth/reset_passwd

重置密码

## User（用户相关接口）

此路径下所有接口均需登录态

### UserSelf（用户数据）

#### PUT /user/profile

更新用户资料

#### PUT /user/passwd

更新用户密码

#### PUT /user/setting

更新用户设置

#### POST /user/upload

用户上传照片

### UserAlbum（用户专辑）

#### POST /user/album

创建专辑

#### PUT /user/album/:id

更新专辑属性

#### DELETE /user/album/:id

删除专辑

#### GET /user/album

列出专辑列表

#### GET /user/album/names

列出个人专辑及认领的专辑名称

#### GET /user/album/:id"

专辑详情

#### GET /user/album/:id/fairy

仅列出专辑下照片

### UserFairy（用户专辑照片）

#### POST /user/fairy

创建照片

#### DELETE /user/fairy/:id

删除照片

### UserCliam（用户认领专辑）

所谓认领，即用户主动申请共享其他用户的某张专辑，认领的专辑，认领者可以上传及修改大部分属性。

专辑只能共享给一个人，认领后可以取消。

#### POST /user/claim/:owner/:id

创建一个认领申请

#### DELETE /user/claim/:owner/:id

删除认领的专辑

#### GET /user/claim

列出认领的专辑

#### GET /user/claim/:owner/:id

获取认领专辑详情
