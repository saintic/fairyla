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

package album

import "fairyla/internal/db"

// 专辑可选系统内置标签
var Labels = []string{"亲人", "爱人", "暗恋", "偶像", "动漫"}

// 专辑属性
type Album struct {
	ID     string   `json:"id"`    // 专辑ID，唯一性，索引，由Name而来
	Owner  string   `json:"owner"` // 所属用户
	Name   string   `json:"name"`  // 专辑名称，不具备唯一性
	CTime  uint64   `json:"ctime"`
	Public bool     `json:"public"`
	Label  []string `json:"label"`
}

// 照片属性
type Fairy struct {
	Album string `json:"album"` // 所属专辑
	CTime uint64 `json:"ctime"`
	Desc  string `json:"desc"`
	Src   string `json:"src"`
}

type wrap struct {
	*db.Conn
}

func New(redis_url string) (w wrap, err error) {
	c, err := db.New(redis_url)
	if err != nil {
		return
	}
	return wrap{c}, nil
}

func (w wrap) CreateAlbum(a Album) {
	//
}
