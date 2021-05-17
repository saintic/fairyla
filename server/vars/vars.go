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
	UserIndex  = "users"                    // redis type set
	GenUserKey = func(user string) string { // redis type hash
		// :user:auth 密码
		// :user:profile 资料
		// :user:setting 配置
		return "user:" + user
	}

	// :album:<User> Hash类型 Key是专辑ID Value是专辑属性，JSON格式
	GenAlbumKey = func(user string) string {
		return "album:" + user
	}
	// :fairy:<User>:<AlbumID> Hash类型 Key是照片ID Value是照片属性，JSON格式
	GenFairyKey = func(user, albumID string) string {
		return fmt.Sprintf("fairy:%s:%s", user, albumID)
	}
	// :claim:<User> Set类型 Key格式是 owner:album_id
	GenClaimKey = func(user string) string {
		return fmt.Sprintf("claim:%s", user)
	}

	AllowImage = []string{".png", ".jpg", ".jpeg", ".gif", ".webp"}
	AllowVideo = []string{".mp4", ".ogg", ".ogv", ".webm", ".3gp", ".mov"}
)

const (
	AlbumPreID = "md.a."
	FairyPreID = "md.f."

	// 上传限制，单位MB
	UploadLimitSize int64 = 20
	// 用户专辑数量限制
	AlbumLimitNum = 9

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

func ResErr(lang, msg string) res {
	return res{false, msgTranslator(lang, msg)}
}

func NewResData(data interface{}) resData {
	return resData{res{true, "ok"}, data}
}
