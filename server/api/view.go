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
	c.JSON(code, vars.ResErr(getLocale(c), msg))
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

func createAlbumView(c echo.Context) error {
	name := c.FormValue("name")
	labels := c.FormValue("labels")
	a, err := album.NewAlbum(getUser(c), name)
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

func dropAlbumView(c echo.Context) error {
	user := getUser(c)
	albumID := autoAlbumID(c)
	w := album.New(rc)
	err := w.DropAlbum(user, albumID)
	if err != nil {
		return err
	}
	return c.JSON(200, vars.ResOK())
}

func createFairyView(c echo.Context) error {
	w := album.New(rc)
	user := getUser(c)
	albumID := c.FormValue("album_id")
	albumName := c.FormValue("album")
	src := c.FormValue("src")
	desc := c.FormValue("desc")
	var a *album.Album
	if albumID == "" {
		if albumName == "" {
			return errors.New("invalid album_id or album")
		}
		a, _ = album.NewAlbum(user, albumName)
		albumID = a.ID
	} else {
		ta, err := w.GetAlbum(user, albumID)
		if err != nil {
			return err
		}
		a = &ta
	}
	f, err := album.NewFairy(user, albumID, src, desc)
	if err != nil {
		return err
	}
	a.SetLatest(f)
	err = w.WriteAlbum(a)
	if err != nil {
		return err
	}
	err = w.WriteFairy(f)
	if err != nil {
		return err
	}
	return c.JSON(200, vars.NewResData(f))
}

func dropFairyView(c echo.Context) error {
	user := getUser(c)
	albumID := autoAlbumID(c)
	fairyID := c.FormValue("fairy_id")
	if fairyID == "" {
		return errors.New("invalid param")
	}
	w := album.New(rc)
	err := w.DropFairy(user, albumID, fairyID)
	if err != nil {
		return err
	}
	return c.JSON(200, vars.ResOK())
}

func listAlbumView(c echo.Context) error {
	w := album.New(rc)
	data, err := w.ListAlbums(getUser(c))
	if err != nil {
		return err
	}
	return c.JSON(200, vars.NewResData(data))
}

func getAlbumView(c echo.Context) error {
	w := album.New(rc)
	user := getUser(c)
	albumID := autoAlbumID(c)
	if gtc.IsTrue(c.QueryParam("fairy")) {
		data, err := w.GetAlbumFairies(user, albumID)
		if err != nil {
			return err
		}
		return c.JSON(200, vars.NewResData(data))
	} else {
		data, err := w.GetAlbum(user, albumID)
		if err != nil {
			return err
		}
		return c.JSON(200, vars.NewResData(data))
	}
}

func getAlbumFairyView(c echo.Context) error {
	w := album.New(rc)
	user := getUser(c)
	albumID := autoAlbumID(c)
	data, err := w.GetFairies(user, albumID)
	if err != nil {
		return err
	}
	return c.JSON(200, vars.NewResData(data))
}

func listFairyView(c echo.Context) error {
	w := album.New(rc)
	data, err := w.ListFairies(getUser(c))
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
	post.Form.Add("album", "fairyla")

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
