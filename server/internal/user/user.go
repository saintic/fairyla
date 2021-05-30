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
	Name  string // 用户名，唯一
	Alias string // 昵称，别名
	Bio   string
}

type Setting struct {
	AlbumDefaultPublic bool
}

type User struct {
	Profile
	Setting
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
	has, err := w.HasUser(user)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("not found user")
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
		return errors.New("not found user")
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

func (w wrap) UserData(user string) (u User, err error) {
	data, err := w.HGetAll(vars.GenUserKey(user))
	if err != nil {
		return
	}
	p := Profile{}
	s := Setting{}
	err = json.Unmarshal([]byte(data[pmodule]), &p)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(data[smodule]), &s)
	if err != nil {
		return
	}
	u = User{p, s}
	return
}
