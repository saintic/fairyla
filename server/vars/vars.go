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
	// :fairy:<AlbumID> Hash类型 Key是照片ID Value是照片属性，JSON格式
	GenFairyKey = func(albumID string) string {
		return "fairy:" + albumID
	}
)

type (
	Res struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	ResToken struct {
		Res
		Token string `json:"token"`
	}
)

func ResOK() Res {
	return Res{true, "ok"}
}

func ResErr(msg string) Res {
	return Res{false, msg}
}

func NewResToken(token string) ResToken {
	return ResToken{Res{true, "ok"}, token}
}
