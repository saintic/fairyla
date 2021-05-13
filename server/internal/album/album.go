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
	"log"
	"net/url"
	"sort"

	"fairyla/pkg/db"
	"fairyla/pkg/util"
	"fairyla/vars"

	"tcw.im/gtc"
)

// 专辑属性
type Album struct {
	ID          string   `json:"id"`    // 专辑ID，唯一性，索引，由Name而来
	Name        string   `json:"name"`  // 专辑名称，不具备唯一性
	Owner       string   `json:"owner"` // 所属用户
	Ta          string   `json:"ta"`    // 认领用户
	CTime       int64    `json:"ctime"`
	Public      bool     `json:"public"`
	Label       []string `json:"label"`
	LatestFairy *Fairy   `json:"latest_fairy"` // 最近上传的fairy
	SteadyFairy *Fairy   `json:"steady_fairy"` // 固定设置的fairy
}

// 照片属性
type Fairy struct {
	ID      string `json:"id"`       // 专辑ID，唯一性，索引，由Src而来
	AlbumID string `json:"album_id"` // 所属专辑
	Owner   string `json:"owner"`    // 所属用户（上传者，专辑属主或认领者）
	CTime   int64  `json:"ctime"`
	Desc    string `json:"desc"`
	Src     string `json:"src"`      // 照片存储地址，理论上要求唯一
	IsVideo bool   `json:"is_video"` // 是否为视频
}

// 专辑及其包含的照片
type AlbumFairy struct {
	Album
	Fairy []Fairy `json:"fairy"`
}

// 专辑名转成ID
func AlbumName2ID(owner, name string) string {
	return vars.AlbumPreID + gtc.MD5(owner+name)
}

func NewAlbum(owner, name string) (a *Album, err error) {
	if owner == "" || name == "" {
		err = errors.New("invalid fairy param")
		return
	}
	// Name在用户中唯一，即ID唯一
	a = &Album{
		ID: AlbumName2ID(owner, name), Owner: owner, Name: name, Public: true,
		CTime: util.Now(),
	}
	return
}

func (a *Album) AddLabel(label string) {
	a.Label = append(a.Label, label)
}

func (a *Album) SetLatest(f *Fairy) {
	a.LatestFairy = f
}

func (a *Album) SetSteady(f *Fairy) {
	a.SteadyFairy = f
}

func (a *Album) Exist(rc *db.Conn) (bool, error) {
	return rc.HExists(vars.GenAlbumKey(a.Owner), a.ID)
}

func NewFairy(owner, albumID, src, desc string) (f *Fairy, err error) {
	if albumID == "" || owner == "" {
		err = errors.New("invalid fairy param")
		return
	}
	if !util.IsValidURL(src) {
		err = errors.New("illegal fairyl src url")
		return
	}
	u, _ := url.ParseRequestURI(src)
	isVideo := util.IsVideo(u.Path)
	isImage := util.IsImage(u.Path)
	if !isVideo && !isImage {
		err = errors.New("unsupported file type")
		return
	}
	now := util.Now()
	ID := fmt.Sprintf("%s-%s-%d", albumID, src, now)
	f = &Fairy{
		vars.FairyPreID + gtc.MD5(ID), albumID, owner, now, desc, src, isVideo,
	}
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
func (w wrap) WriteAlbum(a *Album) error {
	// check param
	if a.Owner == "" || a.ID == "" {
		return errors.New("invalid album param")
	}
	// check db
	index := vars.GenAlbumKey(a.Owner)
	aIDs, err := w.HKeys(index)
	if err != nil {
		return err
	}
	if len(aIDs) >= vars.AlbumLimitNum && !gtc.StrInSlice(a.ID, aIDs) {
		return errors.New("the number of albums exceeds the limit")
	}
	// write db, if exists(=update)
	val, err := json.Marshal(a)
	if err != nil {
		return err
	}
	_, err = w.HSet(index, a.ID, string(val))
	return err
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
	pipe := w.Pipeline()
	pipe.HDel(vars.GenAlbumKey(owner), albumID)
	pipe.Del(vars.GenFairyKey(owner, albumID))
	_, err := pipe.Execute()
	return err
}

// Only check the basic parameters and (overwrite) write
func (w wrap) WriteFairy(f *Fairy) error {
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
	return err
}

// 删除用户某个专辑下某个照片
func (w wrap) DropFairy(owner, albumID, fairyID string) error {
	index := vars.GenFairyKey(owner, albumID)
	_, err := w.HDel(index, fairyID)
	return err
}

// 列出用户所有专辑数据（不包含专辑下照片）
func (w wrap) ListAlbums(user string) (albums []Album, err error) {
	data, err := w.HGetAll(vars.GenAlbumKey(user))
	if err != nil {
		return
	}
	albums = make([]Album, 0, len(data))
	for _, v := range data {
		var a Album
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
func (w wrap) GetAlbum(user, albumID string) (a Album, err error) {
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
func (w wrap) GetAlbumFairies(user, albumID string) (af AlbumFairy, err error) {
	a, err := w.GetAlbum(user, albumID)
	if err != nil {
		return
	}
	f, err := w.GetFairies(user, albumID)
	if err != nil {
		return
	}
	return AlbumFairy{a, f}, nil
}

// 列出用户所有专辑ID及其下照片数据
func (w wrap) ListFairies(user string) (out map[string][]Fairy, err error) {
	data, err := w.HGetAll(vars.GenAlbumKey(user))
	if err != nil {
		return
	}
	out = make(map[string][]Fairy)
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

// 仅获取用户某个专辑下所有照片数据（不包含专辑数据）
func (w wrap) GetFairies(user, albumID string) (fairies []Fairy, err error) {
	data, err := w.HGetAll(vars.GenFairyKey(user, albumID))
	if err != nil {
		return
	}
	fairies = make([]Fairy, 0, len(data))
	for _, v := range data {
		var f Fairy
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
func (w wrap) GetFairy(user, albumID, fairyID string) (f Fairy, err error) {
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

// 列出所有用户所有公开专辑数据（不包含专辑下照片）
func (w wrap) ListPublicAlbums() (data []Album, err error) {
	users, err := w.SMembers(vars.UserIndex)
	if err != nil {
		return
	}
	pipe := w.Pipeline()
	for _, user := range users {
		pipe.Send("HVALS", vars.GenAlbumKey(user))
	}
	rs, err := pipe.Execute()
	if err != nil {
		return
	}
	for _, r := range rs {
		for _, d := range r.([]interface{}) {
			var a Album
			e := json.Unmarshal(d.([]byte), &a)
			if e != nil {
				log.Println(e)
				continue
			}
			if a.Public {
				data = append(data, a)
			}
		}
	}
	return
}

// 列出所有用户所有公开专辑数据（包含专辑下照片）
// 当参数 albumIDs 不为空时，返回指定的专辑ID对应数据
// 当参数 albumNames 不为空时，返回切指定的专辑名对应数据
func (w wrap) ListPublicAlbumFaries(albumIDs, albumNames []string) (data []AlbumFairy, err error) {
	albums, err := w.ListPublicAlbums()
	if err != nil {
		return
	}
	for _, a := range albums {
		isUse := false
		if len(albumIDs) > 0 || len(albumNames) > 0 {
			// 说明仅返回两个参数指定的专辑即可
			if gtc.StrInSlice(a.ID, albumIDs) {
				isUse = true
			}
			if gtc.StrInSlice(a.Name, albumNames) {
				isUse = true
			}
		} else {
			isUse = true
		}
		if isUse {
			f, e := w.GetFairies(a.Owner, a.ID)
			if e != nil {
				err = e
				return
			}
			data = append(data, AlbumFairy{a, f})
		}
	}
	return
}
