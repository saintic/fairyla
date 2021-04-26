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
	"errors"
	"fairyla/internal/album"
	"fairyla/internal/user/auth"
	"fairyla/vars"
	"fmt"

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
	a, err := album.NewAlbum(c.Get("user").(string), name)
	if err != nil {
		return err
	}
	w := album.New(rc)
	err = w.CreateAlbum(a)
	if err != nil {
		return err
	}
	return c.JSON(200, vars.NewResData(a))
}

func createFairyView(c echo.Context) error {
	user := c.Get("user").(string)
	albumID := c.FormValue("album_id")
	albumName := c.FormValue("album_name")
	src := c.FormValue("src")
	desc := c.FormValue("desc")
	w := album.New(rc)
	if albumID == "" {
		// auto create album
		if albumName == "" {
			return errors.New("invalid album_id or album_name")
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
			err := w.CreateAlbum(a)
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
	err = w.CreateFairy(f)
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
