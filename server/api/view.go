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
	"fairyla/internal/user/auth"
	"fmt"

	"github.com/labstack/echo/v4"
)

type res struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type resToken struct {
	res
	Token string `json:"token"`
}

func resOK() res {
	return res{true, "ok"}
}

func resErr(msg string) res {
	return res{false, msg}
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := 400
	msg := err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}
	c.JSON(code, resErr(msg))
}

func signUpView(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	err := auth.Register(rc, username, password)
	if err != nil {
		return err
	}
	return c.JSON(200, resOK())
}

func signInView(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	token, err := auth.Login(rc, username, password)
	if err != nil {
		return err
	}
	return c.JSON(200, resToken{resOK(), token})
}

func signCheckView(c echo.Context) error {
	token := c.FormValue("token")
	data, err := auth.ParseToken(rc, token)
	fmt.Println("sign view:", data, err)
	if err != nil {
		return err
	}
	return c.JSON(200, resOK())
}

func testView(c echo.Context) error {
    return c.String(200, "ok")
}