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
	"regexp"
)

// API返回字段msg国际化翻译
func msgTranslator(locale, msg string) string {
	if locale == "" || msg == "" || msg == "ok" {
		return msg
	}
	if trans, has := precise[locale]; has {
		if newMsg, has := trans[msg]; has {
			return newMsg
		}
	}
	fmt.Printf("msg translator into fuzzy: %s\n", msg)
	if trans, has := fuzzy[locale]; has {
		for pat, repl := range trans {
			if pat.MatchString(msg) {
				return pat.ReplaceAllString(msg, repl)
			}
		}
	}
	return msg
}

var precise = map[string]map[string]string{
	"zh": {
		"Hello world": "你好，世界",
		// http status message
		"Not Found":             "页面未发现",
		"Internal Server Error": "服务器错误",
		"Method Not Allowed":    "方法不允许",

		// common message
		"invalid param":            "无效参数",
		"missing or malformed jwt": "Token丢失或格式不正确",
		"invalid or expired jwt":   "Token无效或已过期",

		// project message
		"album already exists":      "专辑已存在",
		"invalid album param":       "无效的album参数",
		"invalid album_id or album": "无效的album_id或album参数",
		"invalid fairy param":       "无效的fairy参数",
		"illegal fairyl src url":    "非法的fairy照片地址",
		"invalid username":          "无效用户名",
		"password is too short":     "密码太短",
		"username already exists":   "用户已存在",
		"not found username":        "未发现用户",
		"wrong password":            "密码错误",
		"illegal token":             "非法Token",

		// project message(too long)
		"the number of albums exceeds the limit": "专辑数量限制",
	},
}

var fuzzy = map[string]map[*regexp.Regexp]string{
	"zh": {
		regexp.MustCompile("hello world"): "你好，世界",
	},
}
