/*
   Copyright 2021 Hiroshi.tao

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package vars

import (
	"fmt"
)

var (
	// 用户索引 Set类型
	UserIndex = "users"
	// :user:<User> Hash类型 auth密码 profile资料 setting配置
	GenUserKey = func(user string) string { // redis type hash
		return "user:" + user
	}

	// :album:<User> Hash类型 Key是专辑ID Value是专辑属性，JSON格式
	GenAlbumKey = func(user string) string {
		return "album:" + user
	}
	// :fairy:<AlbumID> Hash类型 Key是照片ID Value是照片属性，JSON格式
	GenFairyKey = func(albumID string) string {
		return fmt.Sprintf("fairy:%s", albumID)
	}
	// 用户已认领的专辑列表
	// :claim:<User> Set类型 Key是album_id
	GenClaimKey = func(user string) string {
		return fmt.Sprintf("claim:%s", user)
	}
	// :tmp:event:<User> Hash类型 Key是event_id Value是事件内容，JSON格式
	GenEventKey = func(user string) string {
		return fmt.Sprintf("tmp:event:%s", user)
	}

	AllowImage = []string{".png", ".jpg", ".jpeg", ".gif", ".webp"}
	AllowVideo = []string{".mp4", ".ogg", ".ogv", ".webm", ".3gp", ".mov"}
)

const (
	DefaultSiteName = "Fairyla - 是小仙女啦"

	AlbumPre = "md.a."
	FairyPre = "md.f."

	// 上传限制，单位MB
	UploadLimitSize int64 = 20
	// 用户专辑数量限制
	AlbumLimitNum = 99

	RedigoNil = "redigo: nil returned"

	// 允许多种方式获取JWT认证值
	AllowExtraJWT = true
)

type (
	res struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	resData struct {
		res
		Data interface{} `json:"data"`
	}
)

func ResOK() res {
	return res{true, "ok"}
}

func ResErr(msg string) res {
	return res{false, msg}
}

func ResErrLocale(lang, msg string) res {
	return res{false, msgTranslator(lang, msg)}
}

func NewResData(data interface{}) resData {
	return resData{res{true, "ok"}, data}
}
