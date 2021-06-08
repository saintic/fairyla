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
	"bytes"
	"text/template"
)

const forgotTpl = `
<table
    style="width: 600px; border: 1px solid #ddd; border-radius: 3px; color: #555; font-size: 12px; height: auto; margin: auto; overflow: hidden; text-align: left; word-break: break-all; word-wrap: break-word;"
    cellspacing="0" cellpadding="0" border="0">
    <tbody style="margin: 0; padding: 0;">
        <tr style="background-color: #393D49; height: 60px; margin: 0; padding: 0;">
            <td style="margin: 0; padding: 0;">
                <div style="color: #5EB576; margin: 0; margin-left: 30px; padding: 0;">
                    <a style="font-size: 14px; margin: 0; padding: 0; color: #5EB576; text-decoration: none;"
                        href="{{ .SiteURL }}" target="_blank" rel="noopener">
                        {{ .SiteName }}
                    </a>
                </div>
            </td>
        </tr>
        <tr style="margin: 0; padding: 0;">
            <td style="margin: 0; padding: 30px;">
                <p style="line-height: 20px; margin: 0; margin-bottom: 10px; padding: 0;">
                    Hi，<em style="font-weight: 700;">{{ .User }}</em>，请完成以下操作（10分钟内有效）：
                </p>
                <div>
                    <a href="{{ .URL }}"
                        style="background-color: #009E94; color: #fff; display: inline-block; height: 32px; line-height: 32px; margin: 0 15px 0 0; padding: 0 15px; text-decoration: none;"
                        target="_blank" rel="noopener">
                        点击重置密码
                    </a>
                </div>
                <p
                    style="line-height: 20px; margin-top: 20px; padding: 10px; background-color: #f2f2f2; font-size: 12px;">
                    点击无效？手动复制链接到浏览器访问：{{ .URL }}
                </p>
            </td>
        </tr>
        <tr style="background-color: #fafafa; color: #999; height: 35px; margin: 0; padding: 0; text-align: center;">
            <td style="margin: 0; padding: 0;">如非本人，请勿操作！系统邮件，请勿回复！</td>
        </tr>
    </tbody>
</table>
`

type forgotDat struct {
	SiteURL, SiteName, User, URL string
}

func NewForgot(siteURL, siteName, user, verifyURL string) (content string, err error) {
	t, err := template.New("forgot").Parse(forgotTpl)
	if err != nil {
		return
	}
	fd := forgotDat{siteURL, siteName, user, verifyURL}
	var tpl bytes.Buffer
	if e := t.Execute(&tpl, fd); e != nil {
		err = e
		return
	}
	content = tpl.String()
	return
}
