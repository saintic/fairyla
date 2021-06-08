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
	"log"
	"regexp"
)

// API返回字段msg国际化翻译
func msgTranslator(locale, msg string) string {
	if locale == "" || locale == "en" || msg == "" || msg == "ok" {
		return msg
	}
	if trans, has := precise[locale]; has {
		if newMsg, has := trans[msg]; has {
			return newMsg
		}
	}
	log.Printf("msg translator into fuzzy(%s): %s\n", locale, msg)
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
		"no data":               "没有数据",

		// common message
		"invalid param":            "无效参数",
		"missing or malformed jwt": "Token丢失或格式不正确",
		"invalid or expired jwt":   "Token无效或已过期",
		RedigoNil:                  "暂无数据",
		"Invalid Api Token":        "无效的Token",

		// project message
		"album already exists":     "专辑已存在",
		"invalid album param":      "无效的album参数",
		"invalid fairy param":      "无效的fairy参数",
		"not found album":          "未发现专辑",
		"not found fairy":          "未发现照片",
		"invalid username":         "无效用户名",
		"password is too short":    "密码太短",
		"username already exists":  "用户已存在",
		"not found username":       "未发现用户",
		"wrong password":           "密码错误",
		"illegal token":            "非法Token",
		"unsupported file type":    "不支持的文件类型",
		"invalid action param":     "无效的action参数",
		"not found claim":          "未发现此认领专辑",
		"cannot share to yourself": "不能分享给自己",
		"cannot claim your album":  "不能认领自己的专辑",
		"invalid classify":         "无效的classify参数",
		"pending for approval":     "等待审批",
		"invalid email":            "邮箱错误",
		"refuse to delete":         "拒绝删除",

		// project message(too long)
		"the number of albums exceeds the limit": "专辑数量限制",
		"the uploaded file exceeds the limit":    "上传文件超出限制",
		"invalid album_id or album_name":         "无效的专辑id或name",
		"illegal fairyl src url":                 "非法的fairy照片地址",
		"album id does not match user":           "专辑id与用户不匹配",
	},
}

var fuzzy = map[string]map[*regexp.Regexp]string{
	"zh": {
		regexp.MustCompile("(.*)connection refused"): "${1}连接拒绝",
		regexp.MustCompile("already belong (.*)"):    "已经属于${1}",
	},
}
