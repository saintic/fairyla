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
	"fmt"
	"sort"

	"fairyla/pkg/db"
	"fairyla/pkg/util"
	"fairyla/vars"

	"tcw.im/gtc"
)

var (
	albumLimitNum uint64 = 9 // 用户专辑数量限制
)

// 专辑属性
type album struct {
	ID          string   `json:"id"`    // 专辑ID，唯一性，索引，由Name而来
	Owner       string   `json:"owner"` // 所属用户
	Name        string   `json:"name"`  // 专辑名称，不具备唯一性
	CTime       int64    `json:"ctime"`
	Public      bool     `json:"public"`
	Label       []string `json:"label"`
	LatestFairy *fairy   `json:"latest_fairy"` // 最近上传的fairy
	SteadyFairy *fairy   `json:"steady_fairy"` // 固定设置的fairy
}

// 照片属性
type fairy struct {
	ID      string `json:"id"`       // 专辑ID，唯一性，索引，由Src而来
	AlbumID string `json:"album_id"` // 所属专辑
	Owner   string `json:"owner"`    // 所属用户
	CTime   int64  `json:"ctime"`
	Desc    string `json:"desc"`
	Src     string `json:"src"` // 照片存储地址，理论上要求唯一
}

// 专辑及其包含的照片
type albumFairy struct {
	album
	Fairy []fairy `json:"fairy"`
}

// 专辑名转成ID
func AlbumName2ID(owner, name string) string {
	return vars.AlbumPreID + gtc.MD5(owner+name)
}

func NewAlbum(owner, name string) (a *album, err error) {
	if owner == "" || name == "" {
		err = errors.New("invalid fairy param")
		return
	}
	// Name在用户中唯一，即ID唯一
	a = &album{
		ID: AlbumName2ID(owner, name), Owner: owner, Name: name, Public: true,
		CTime: util.Now(),
	}
	return
}

func (a *album) AddLabel(label string) {
	a.Label = append(a.Label, label)
}

func (a *album) SetLatest(f *fairy) {
	a.LatestFairy = f
}

func (a *album) SetSteady(f *fairy) {
	a.SteadyFairy = f
}

func (a *album) Exist(rc *db.Conn) (bool, error) {
	return rc.HExists(vars.GenAlbumKey(a.Owner), a.ID)
}

func NewFairy(owner, albumID, src, desc string) (f *fairy, err error) {
	if albumID == "" || owner == "" {
		err = errors.New("invalid fairy param")
		return
	}
	if !util.IsValidURL(src) {
		err = errors.New("illegal fairyl src url")
		return
	}
	now := util.Now()
	ID := fmt.Sprintf("%s-%s-%d", albumID, src, now)
	f = &fairy{vars.FairyPreID + gtc.MD5(ID), albumID, owner, now, desc, src}
	return
}

// 对外接口
type wrap struct {
	*db.Conn
}

func New(c *db.Conn) wrap {
	return wrap{c}
}

// Only check the basic parameters and (overwrite) write
func (w wrap) WriteAlbum(a *album) error {
	// check param
	if a.Owner == "" || a.ID == "" {
		return errors.New("invalid album param")
	}
	// check db
	index := vars.GenAlbumKey(a.Owner)
	length, err := w.HLen(index)
	if err != nil {
		return err
	}
	if length > albumLimitNum {
		return errors.New("the number of albums exceeds the limit")
	}
	// write db, if exists(=update)
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

// 删除用户所有专辑及其下照片
func (w wrap) DropAlbums(owner string) error {
	index := vars.GenAlbumKey(owner)
	albumIDs, err := w.HKeys(index)
	if err != nil {
		return err
	}
	pipe := w.Pipeline()
	for _, aid := range albumIDs {
		pipe.Del(vars.GenFairyKey(owner, aid))
	}
	_, err = pipe.Execute()
	return err
}

// 删除用户某个专辑及其下照片
func (w wrap) DropAlbum(owner, albumID string) error {
	_, err := w.HDel(vars.GenAlbumKey(owner), albumID)
	if err != nil {
		return err
	}
	_, err = w.Del(vars.GenFairyKey(owner, albumID))
	if err != nil {
		return err
	}
	return nil
}

// Only check the basic parameters and (overwrite) write
func (w wrap) WriteFairy(f *fairy) error {
	// check param
	if f.ID == "" || f.Owner == "" || f.AlbumID == "" || f.Src == "" {
		return errors.New("invalid fairy param")
	}
	// write db
	index := vars.GenFairyKey(f.Owner, f.AlbumID)
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

// 删除用户某个专辑下某个照片
func (w wrap) DropFairy(owner, albumID, fairyID string) error {
	index := vars.GenFairyKey(owner, albumID)
	_, err := w.HDel(index, fairyID)
	return err
}

// 列出用户所有专辑数据（不包含专辑下照片）
func (w wrap) ListAlbums(user string) (albums []album, err error) {
	data, err := w.HGetAll(vars.GenAlbumKey(user))
	if err != nil {
		return
	}
	albums = make([]album, 0, len(data))
	for _, v := range data {
		var a album
		e := json.Unmarshal([]byte(v), &a)
		if e == nil {
			albums = append(albums, a)
		}
	}
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].CTime > albums[j].CTime
	})
	return
}

// 获取用户某张专辑数据（不包含专辑下照片）
func (w wrap) GetAlbum(user, albumID string) (a album, err error) {
	val, err := w.HGet(vars.GenAlbumKey(user), albumID)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(val), &a)
	if err != nil {
		return
	}
	return
}

// 获取用户某张专辑数据（包含专辑下照片）
func (w wrap) GetAlbumFairies(user, albumID string) (af albumFairy, err error) {
	a, err := w.GetAlbum(user, albumID)
	if err != nil {
		return
	}
	f, err := w.GetFairies(user, albumID)
	if err != nil {
		return
	}
	return albumFairy{a, f}, nil
}

// 列出用户所有专辑ID及其下照片数据
func (w wrap) ListFairies(user string) (out map[string][]fairy, err error) {
	data, err := w.HGetAll(vars.GenAlbumKey(user))
	if err != nil {
		return
	}
	out = make(map[string][]fairy)
	for albumID := range data {
		fs, e := w.GetFairies(user, albumID)
		if e != nil {
			err = e
			return
		}
		out[albumID] = fs
	}
	return
}

// 获取用户某个专辑下所有照片数据（不包含专辑数据）
func (w wrap) GetFairies(user, albumID string) (fairies []fairy, err error) {
	data, err := w.HGetAll(vars.GenFairyKey(user, albumID))
	if err != nil {
		return
	}
	fairies = make([]fairy, 0, len(data))
	for _, v := range data {
		var f fairy
		e := json.Unmarshal([]byte(v), &f)
		if e == nil {
			fairies = append(fairies, f)
		}
	}
	sort.Slice(fairies, func(i, j int) bool {
		return fairies[i].CTime > fairies[j].CTime
	})
	return
}

// 获取用户某个专辑某个照片数据
func (w wrap) GetFairy(user, albumID, fairyID string) (f fairy, err error) {
	val, err := w.HGet(vars.GenFairyKey(user, albumID), albumID)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(val), &f)
	if err != nil {
		return
	}
	return
}
