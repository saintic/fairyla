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

import (
	"encoding/json"
	"errors"

	"fairyla/pkg/db"
	"fairyla/pkg/util"
	"fairyla/vars"

	"tcw.im/ufc"
)

var (
	albumLimitNum uint64 = 5

	// 专辑可选系统内置标签
	SysLabels = []string{"亲人", "爱人", "暗恋", "偶像", "动漫"}
)

// 专辑属性
type album struct {
	ID     string   `json:"id"`    // 专辑ID，唯一性，索引，由Name而来
	Owner  string   `json:"owner"` // 所属用户
	Name   string   `json:"name"`  // 专辑名称，不具备唯一性
	CTime  int64    `json:"ctime"`
	Public bool     `json:"public"`
	Label  []string `json:"label"`
}

// 照片属性
type fairy struct {
	ID      string `json:"id"`       // 专辑ID，唯一性，索引，由Src而来
	AlbumID string `json:"album_id"` // 所属专辑
	CTime   int64  `json:"ctime"`
	Desc    string `json:"desc"`
	Src     string `json:"src"` // 照片存储地址，理论上要求唯一
}

func NewAlbum(owner, name string) (a *album, err error) {
	if owner == "" || name == "" {
		err = errors.New("invalid fairy param")
		return
	}
	a = &album{
		ID: ufc.MD5(owner + name), Owner: owner, Name: name, Public: true,
		CTime: util.Now(),
	}
	return
}

func (a *album) AddLabel(label string) {
	a.Label = append(a.Label, label)
}

func NewFairy(albumID, src, desc string) (f *fairy, err error) {
	if albumID == "" || src == "" {
		err = errors.New("invalid fairy param")
		return
	}
	f = &fairy{ufc.MD5(albumID + src), albumID, util.Now(), desc, src}
	return
}

// 对外接口
type wrap struct {
	*db.Conn
}

func New(c *db.Conn) wrap {
	return wrap{c}
}

func (w wrap) CreateAlbum(a *album) error {
	// check param
	if a.Owner == "" || a.ID == "" {
		return errors.New("invalid album param")
	}
	// check db
	index := vars.GenAlbumKey(a.Owner) // redis type hash
	length, err := w.HLen(index)
	if err != nil {
		return err
	}
	if length > albumLimitNum {
		return errors.New("the number of albums exceeds the limit")
	}
	// write db
	val, err := json.Marshal(a)
	if err != nil {
		return err
	}
	_, err = w.HSet(index, a.ID, string(val))
	if err != nil {
		return err
	}
	return nil
}

func (w wrap) CreateFairy(f *fairy) error {
	// check param
	if f.ID == "" || f.AlbumID == "" || f.Src == "" {
		return errors.New("invalid fairy param")
	}
	// write db
	index := vars.GenFairyKey(f.AlbumID)
	val, err := json.Marshal(f)
	if err != nil {
		return err
	}
	_, err = w.HSet(index, f.ID, string(val))
	if err != nil {
		return err
	}
	return nil
}

func (w wrap) ShowAlbums(user string) error {
	index := vars.GenAlbumKey(user)
	w.HGetAll(index)
	return nil
}
