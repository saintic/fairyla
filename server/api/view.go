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

package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"fairyla/internal/album"
	"fairyla/internal/user/auth"
	"fairyla/vars"

	"github.com/labstack/echo/v4"
	"tcw.im/gtc"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := 400
	msg := err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}
	c.JSON(code, vars.ResErr(msg))
}

func signUpView(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	err := auth.Register(rc, username, password)
	if err != nil {
		return err
	}
	return c.JSON(200, vars.ResOK())
}

func signInView(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	remember := gtc.IsTrue(c.FormValue("remember"))
	fmt.Println(username, password, remember)
	expire := 60 * 60 * 2 // 2h
	if remember {
		expire = 60 * 60 * 24 * 7 // 7d
	}
	token, err := auth.Login(rc, username, password, remember)
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	data["token"] = token
	data["expire"] = expire
	return c.JSON(200, vars.NewResData(data))
}

func testView(c echo.Context) error {
	return c.JSON(200, vars.ResOK())
}

func createAlbumView(c echo.Context) error {
	name := c.FormValue("name")
	labels := c.FormValue("labels")
	a, err := album.NewAlbum(c.Get("user").(string), name)
	if err != nil {
		return err
	}
	has, err := a.Exist(rc)
	if err != nil {
		return err
	}
	if has {
		return errors.New("album already exists")
	}
	if labels != "" {
		for _, label := range strings.Split(labels, ",") {
			a.AddLabel(strings.TrimSpace(label))
		}
	}
	w := album.New(rc)
	err = w.WriteAlbum(a)
	if err != nil {
		return err
	}
	return c.JSON(200, vars.NewResData(a))
}

func createFairyView(c echo.Context) error {
	user := c.Get("user").(string)
	albumID := c.FormValue("album_id")
	albumName := c.FormValue("album")
	src := c.FormValue("src")
	desc := c.FormValue("desc")
	w := album.New(rc)
	if albumID == "" {
		// auto create album
		if albumName == "" {
			return errors.New("invalid album_id or album")
		}
		a, err := album.NewAlbum(user, albumName)
		if err != nil {
			return err
		}
		has, err := a.Exist(rc)
		if err != nil {
			return err
		}
		if !has {
			err := w.WriteAlbum(a)
			if err != nil {
				return err
			}
		}
		albumID = a.ID
	}
	f, err := album.NewFairy(user, albumID, src, desc)
	if err != nil {
		return err
	}
	err = w.WriteFairy(f)
	if err != nil {
		return err
	}
	return c.JSON(200, vars.NewResData(f))
}

func listAlbumView(c echo.Context) error {
	w := album.New(rc)
	data, err := w.ListAlbum(c.Get("user").(string))
	if err != nil {
		return err
	}
	return c.JSON(200, vars.NewResData(data))
}

func listFairyView(c echo.Context) error {
	w := album.New(rc)
	data, err := w.ListFairy(c.Get("user").(string))
	if err != nil {
		return err
	}
	return c.JSON(200, vars.NewResData(data))
}

func configView(c echo.Context) error {
	data := cfg.SitePublic()
	data["isLogin"] = false
	user, err := checkJWT(c)
	if err == nil {
		data["isLogin"] = true
	}
	data["user"] = user
	return c.JSON(200, vars.NewResData(data))
}

func uploadView(c echo.Context) error {
	album := c.FormValue("album")
	title := c.FormValue("title")
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	fd, err := file.Open()
	if err != nil {
		return err
	}
	defer fd.Close()

	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}
	pic := base64.StdEncoding.EncodeToString(content)

	var post http.Request
	post.ParseForm()
	post.Form.Add(cfg.Sapic.Field, pic)
	post.Form.Add("album", album)
	post.Form.Add("title", title)

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(
		"POST", cfg.Sapic.Api, strings.NewReader(post.Form.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", echo.MIMEApplicationForm)
	req.Header.Set("User-Agent", "fairyla/v1")
	req.Header.Set("Authorization", "LinkToken "+cfg.Sapic.Token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	ret := struct {
		Code int    `json:"-"`
		Msg  string `json:"-"`
		Src  string `json:"src"`
	}{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return err
	}
	if ret.Code == 0 {
		return c.JSON(200, vars.NewResData(ret))
	} else {
		return errors.New(ret.Msg)
	}
}
