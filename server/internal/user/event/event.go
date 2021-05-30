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

package event

import (
	"encoding/json"
	"errors"
	"fmt"

	"fairyla/internal/user"
	"fairyla/pkg/db"
	"fairyla/pkg/util"
	"fairyla/vars"

	"tcw.im/gtc"
)

var classifies = []string{"message", "notify", "alert"}

type Message struct {
	Owner string `json:"owner"` // 所属用户
	// 不同消息类型所需配置选项
	Opt *MessageOption `json:"opt"`

	ID    string `json:"id"`
	CTime int64  `json:"ctime"`
}

type MessageOption struct {
	// 消息类型，不同类型在 Element-Plus 有不同触发动作，
	// 即 message notify alert（messagebox）
	Classify string `json:"classify"`
	Title    string `json:"title"`
	Msg      string `json:"msg"`
	Theme    string `json:"theme"` //消息类型（主题样式），用于显示图标
	IsHtml   bool   `json:"is_html"`
	// 特定类型专用参数
	NotifyJump string `json:"notify_jump"` // 跳转路径
}

func NewMessage(owner, msg, classify string) (m *Message, err error) {
	if owner == "" || msg == "" {
		err = errors.New("invalid param")
		return
	}
	if classify == "" {
		classify = "message"
	}
	if !gtc.StrInSlice(classify, classifies) {
		err = errors.New("invalid classify")
		return
	}
	now := util.Now()
	id := fmt.Sprintf("%s:%s:%d", owner, classify, now)
	opt := &MessageOption{
		Classify: classify, Msg: msg, Theme: "info", IsHtml: false,
	}
	m = &Message{owner, opt, gtc.MD5(id), now}
	return
}

func (m *Message) Write(rc *db.Conn) error {
	if m.Owner == "" || m.ID == "" || m.Opt == nil || m.Opt.Msg == "" {
		return errors.New("invalid param")
	}
	// check owner
	w := user.New(rc)
	has, err := w.HasUser(m.Owner)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("not found username")
	}
	val, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = rc.HSet(vars.GenEventKey(m.Owner), m.ID, string(val))
	return err
}

// 对外接口
type wrap struct {
	*db.Conn
}

func New(c *db.Conn) wrap {
	return wrap{c}
}

func (w wrap) ListMessages(owner, classify string) (ms []Message, err error) {
	data, err := w.HVals(vars.GenEventKey(owner))
	if err != nil {
		return
	}
	ms = make([]Message, 0, len(data))
	for _, v := range data {
		var m Message
		e := json.Unmarshal([]byte(v), &m)
		if e == nil {
			if classify == "" || (classify != "" && classify == m.Opt.Classify) {
				ms = append(ms, m)
			}
		}
	}
	return
}

func (w wrap) DropMessages(owner, id string) error {
	_, err := w.HDel(vars.GenEventKey(owner), id)
	return err
}
