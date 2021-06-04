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

package user

import (
	"encoding/json"
	"errors"

	"fairyla/pkg/db"
	"fairyla/pkg/util"
	"fairyla/vars"
)

const (
	pmodule = "profile"
	smodule = "setting"
)

type Profile struct {
	Name  string `json:"name"`  // 用户名，唯一
	Alias string `json:"alias"` // 昵称，别名
	Bio   string `json:"bio"`   // 宣言
	Email string `json:"email"`
}

type Setting struct {
	AlbumDefaultPublic bool   `json:"album_default_public"` //默认专辑状态
	Slogan             string `json:"slogan"`               // 覆盖系统Slogan
}

type User struct {
	Profile
	Setting
}

func NewProfile(user string) Profile {
	return Profile{Name: user}
}

func NewSetting() Setting {
	return Setting{AlbumDefaultPublic: true}
}

// 子模块名称
func Module(m string) string {
	switch m {
	case "p", "profile":
		return pmodule
	case "s", "setting":
		return smodule
	default:
		return ""
	}
}

// 对外接口
type wrap struct {
	*db.Conn
}

func New(c *db.Conn) wrap {
	return wrap{c}
}

func (w wrap) HasUser(user string) (bool, error) {
	return w.SIsMember(vars.UserIndex, user)
}

// 更新用户资料
func (w wrap) UpdateProfile(p Profile) error {
	user := p.Name
	if !util.IsName(user) {
		return errors.New("invalid name")
	}
	if p.Email != "" && !util.IsEmail(p.Email) {
		return errors.New("invalid email")
	}
	has, err := w.HasUser(user)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("not found username")
	}
	val, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = w.HSet(vars.GenUserKey(user), pmodule, string(val))
	if err != nil {
		return err
	}
	return nil
}

// 更新用户设置
func (w wrap) UpdateSetting(user string, s Setting) error {
	has, err := w.HasUser(user)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("not found username")
	}
	val, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = w.HSet(vars.GenUserKey(user), smodule, string(val))
	if err != nil {
		return err
	}
	return nil
}

func (w wrap) UserProfile(user string) (p Profile, err error) {
	index := vars.GenUserKey(user)
	exist, err := w.HExists(index, pmodule)
	if err != nil {
		return
	}
	if exist {
		data, e := w.HGet(index, pmodule)
		if e != nil {
			err = e
			return
		}
		err = json.Unmarshal([]byte(data), &p)
		if err != nil {
			return
		}
	} else {
		p = NewProfile(user)
	}
	return
}

func (w wrap) UserSetting(user string) (s Setting, err error) {
	index := vars.GenUserKey(user)
	exist, err := w.HExists(index, smodule)
	if err != nil {
		return
	}
	if exist {
		data, e := w.HGet(index, smodule)
		if e != nil {
			err = e
			return
		}
		err = json.Unmarshal([]byte(data), &s)
		if err != nil {
			return
		}
	} else {
		s = NewSetting()
	}
	return
}

func (w wrap) UserData(user string) (u User, err error) {
	p, err := w.UserProfile(user)
	if err != nil {
		return
	}
	s, err := w.UserSetting(user)
	if err != nil {
		return
	}
	u = User{p, s}
	return
}
