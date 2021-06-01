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
	"strings"

	"fairyla/internal/user/event"
	"fairyla/pkg/db"
	"fairyla/pkg/util"
	"fairyla/vars"

	"tcw.im/gtc"
)

// 专辑属性
type Album struct {
	ID     string   `json:"id"`    // 专辑ID，唯一性，索引，由User和Name而来
	Name   string   `json:"name"`  // 专辑名称，不具备唯一性
	Owner  string   `json:"owner"` // 所属用户
	Ta     string   `json:"ta"`    // 认领用户
	CTime  int64    `json:"ctime"`
	Public bool     `json:"public"`
	Label  []string `json:"label"`
	// 专辑封面：自动设置为最近上传的或固定的fairy
	LatestFairy *Fairy `json:"latest_fairy"`
	SteadyFairy *Fairy `json:"steady_fairy"`
	// 专辑额外配置项
	Opt *AlbumOption `json:"opt"`
}

type AlbumOption struct {
	ClaimingBy []string `json:"claiming_by"`
}

// 照片属性
type Fairy struct {
	ID      string `json:"id"`       // 专辑ID，唯一性，索引，由Src而来
	AlbumID string `json:"album_id"` // 所属专辑ID
	Creator string `json:"creator"`  // 创建者（可能是专辑属主或认领者）
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
	// such as: md.a.<User>.<MD5>
	return fmt.Sprintf("%s%s.%s", vars.AlbumPre, owner, gtc.MD5(owner+name))
}

// 从专辑ID解析出用户名
func AlbumID2User(albumID string) string {
	aID := strings.TrimPrefix(albumID, vars.AlbumPre)
	return strings.Split(aID, ".")[0]
}

func NewAlbum(owner, name string) (a *Album, err error) {
	if owner == "" || name == "" {
		err = errors.New("invalid album param")
		return
	}
	// Name在用户中唯一，即ID唯一
	a = &Album{
		ID: AlbumName2ID(owner, name), Owner: owner, Name: name, Public: true,
		CTime: util.Now(), Opt: &AlbumOption{},
	}
	return
}

func (a *Album) AddLabel(label string) {
	a.Label = append(a.Label, label)
}

func (a *Album) RemoveLabel(label string) {
	a.Label = util.DeleteSlice(a.Label, label)
}

func (a *Album) SetLatest(f *Fairy) {
	a.LatestFairy = f
}

func (a *Album) SetSteady(f *Fairy) {
	a.SteadyFairy = f
}

// 判断专辑ID是否存在，完全匹配专辑属主和ID
func (a *Album) Exist(rc *db.Conn) (bool, error) {
	return rc.HExists(vars.GenAlbumKey(a.Owner), a.ID)
}

// 专辑认领入库或清除
func (a *Album) Claim(rc *db.Conn, isRemove bool) (err error) {
	if a.Ta != "" {
		if isRemove {
			_, err = rc.SRem(vars.GenClaimKey(a.Ta), a.ID)
		} else {
			_, err = rc.SAdd(vars.GenClaimKey(a.Ta), a.ID)
			// 分享后发送通知
			tpl := fmt.Sprintf(
				"您好，用户【%s】已将专辑【%s】分享给您共同维护。", a.Owner, a.Name,
			)
			m, _ := event.NewMessage(a.Ta, tpl, "message")
			m.Opt.Theme = "success"
			m.Write(rc)
		}
	}
	return
}

func NewFairy(creator, albumID, src, desc string) (f *Fairy, err error) {
	if albumID == "" {
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
	ID := vars.FairyPre + gtc.MD5(fmt.Sprintf("%s-%s-%d", albumID, src, now))
	f = &Fairy{ID, albumID, creator, now, desc, src, isVideo}
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

// 删除用户某个专辑及其下照片
func (w wrap) DropAlbum(owner, albumID string) error {
	a, err := w.GetAlbum(owner, albumID)
	if err != nil {
		return err
	}
	pipe := w.Pipeline()
	pipe.HDel(vars.GenAlbumKey(owner), albumID)
	pipe.Del(vars.GenFairyKey(albumID))
	if a.Ta != "" {
		pipe.SRem(vars.GenClaimKey(a.Ta), albumID)
	}
	_, err = pipe.Execute()
	return err
}

// 删除用户所有专辑及其下照片
func (w wrap) DropAlbums(owner string) error {
	index := vars.GenAlbumKey(owner)
	albumIDs, err := w.HKeys(index)
	if err != nil {
		return err
	}
	for _, aid := range albumIDs {
		err = w.DropAlbum(owner, aid)
		if err != nil {
			return err
		}
	}
	return nil
}

// Only check the basic parameters and (overwrite) write
func (w wrap) WriteFairy(f *Fairy) error {
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
	return err
}

// 删除某个专辑下某个照片
func (w wrap) DropFairy(albumID, fairyID string) error {
	_, err := w.HDel(vars.GenFairyKey(albumID), fairyID)
	return err
}

// 列出用户所有专辑数据（不包含专辑下照片）
func (w wrap) ListAlbums(owner string) (albums []Album, err error) {
	data, err := w.HGetAll(vars.GenAlbumKey(owner))
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
func (w wrap) GetAlbum(owner, albumID string) (a Album, err error) {
	// 可不校验，owner和albumID中的user不对应，获取Nil为空错误直接返回
	if owner != AlbumID2User(albumID) {
		err = errors.New("not found album")
		return
	}
	val, err := w.HGet(vars.GenAlbumKey(owner), albumID)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(val), &a)
	if err != nil {
		return
	}
	if a.Opt == nil {
		a.Opt = &AlbumOption{}
	}
	return
}

// 获取用户某张专辑数据（包含专辑下照片）
func (w wrap) GetAlbumFairies(owner, albumID string) (af AlbumFairy, err error) {
	a, err := w.GetAlbum(owner, albumID)
	if err != nil {
		return
	}
	f, err := w.GetFairies(albumID)
	if err != nil {
		return
	}
	return AlbumFairy{a, f}, nil
}

// 仅获取某个专辑下所有照片数据（即：不包含专辑数据）
func (w wrap) GetFairies(albumID string) (fairies []Fairy, err error) {
	data, err := w.HGetAll(vars.GenFairyKey(albumID))
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

// 列出用户所有专辑ID及其下照片数据（也许要屏蔽此接口）
func (w wrap) ListFairies(user string) (out map[string][]Fairy, err error) {
	albumIDs, err := w.HKeys(vars.GenAlbumKey(user))
	if err != nil {
		return
	}
	out = make(map[string][]Fairy)
	for _, albumID := range albumIDs {
		fs, e := w.GetFairies(albumID)
		if e != nil {
			err = e
			return
		}
		out[albumID] = fs
	}
	return
}

// 获取用户某个专辑某个照片数据
func (w wrap) GetFairy(albumID, fairyID string) (f Fairy, err error) {
	val, err := w.HGet(vars.GenFairyKey(albumID), fairyID)
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

// 是否有此用户
func (w wrap) HasUser(user string) (bool, error) {
	return w.SIsMember(vars.UserIndex, user)
}

// 用户是否有此认领的专辑
func (w wrap) HasClaim(user, albumID string) (bool, error) {
	return w.SIsMember(vars.GenClaimKey(user), albumID)
}

// 列出用户名下所有认领的专辑数据（不包含专辑下照片）
func (w wrap) ListClaimAlbums(owner string) (data []Album, err error) {
	albumIDs, err := w.SMembers(vars.GenClaimKey(owner))
	if err != nil {
		return
	}
	for _, albumID := range albumIDs {
		a, e := w.GetAlbum(AlbumID2User(albumID), albumID)
		if e == nil {
			data = append(data, a)
		}
	}
	return
}

// 认领他人专辑，by是认领者，to、albumID分别是专辑属主和ID
func (w wrap) CreateClaim(by, to, albumID string) error {
	if by == to {
		return errors.New("cannot claim your album")
	}
	has, err := w.HasUser(to)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("not found username")
	}
	a, err := w.GetAlbum(to, albumID)
	if err != nil {
		return err
	}
	exist, err := a.Exist(w.Conn)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("invalid album param")
	}
	if a.Ta != "" && a.Ta == by {
		return errors.New("already belong you")
	}
	bies := a.Opt.ClaimingBy
	if gtc.StrInSlice(by, bies) {
		return errors.New("pending for approval")
	}
	bies = append(bies, by)
	a.Opt.ClaimingBy = bies
	err = w.WriteAlbum(&a)
	if err != nil {
		return err
	}
	tpl := fmt.Sprintf(
		"您好，用户【%s】认领了您的专辑【%s】，希望共同维护，点此处理。", by, a.Name,
	)
	m, _ := event.NewMessage(to, tpl, "notify")
	m.Opt.Title = "专辑认领通知"
	m.Opt.Theme = "info"
	m.Opt.NotifyJump = "/album/" + a.Name
	return m.Write(w.Conn)
}
