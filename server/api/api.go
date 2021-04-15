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
	"fmt"

	"fairyla/internal/sys"
	"fairyla/pkg/db"
	"fairyla/vars"

	"github.com/labstack/echo/v4"
)

var rc *db.Conn

type status struct {
	Logged bool   `json:"logged"`
	User   string `json:"user"`
}

func StartApi(config *sys.Setting) {
	c, err := db.New(config.Redis)
	if err != nil {
		panic(err)
	}
	rc = c

	e := echo.New()
	e.Debug = true
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.GET("/config", func(c echo.Context) error {
		data := config.SitePublic()
		user, err := checkJWT(c)
		userStatus := status{}
		if err == nil {
			userStatus.Logged = true
			userStatus.User = user
		}
		data["status"] = userStatus
		return c.JSON(200, vars.NewResData(data))
	})

	auth := e.Group("/auth")
	auth.POST("/signup", signUpView)
	auth.POST("/signin", signInView)

	test := e.Group("/test", loginRequired)
	test.POST("/check", testView)

	user := e.Group("/user", loginRequired)
	user.GET("/album", listAlbumView)
	user.POST("/album", createAlbumView)
	user.GET("/fairy", listFairyView)
	user.POST("/fairy", createFairyView)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", config.Host, config.Port)))
}
