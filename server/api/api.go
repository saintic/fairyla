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
	"github.com/labstack/echo/v4/middleware"
)

var rc *db.Conn

func StartApi(config *sys.Setting) {
	c, err := db.New(config.Redis)
	if err != nil {
		panic(err)
	}
	rc = c

	e := echo.New()
	e.Debug = true
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Use(middleware.Logger())

	api := e.Group("/api")

	api.GET("/config", func(c echo.Context) error {
		data := config.SitePublic()
		data["isLogin"] = false
		user, err := checkJWT(c)
		if err == nil {
			data["isLogin"] = true
		}
		data["user"] = user

		return c.JSON(200, vars.NewResData(data))
	})

	auth := api.Group("/auth")
	auth.POST("/signup", signUpView)
	auth.POST("/signin", signInView)

	test := api.Group("/test", loginRequired)
	test.POST("/check", testView)

	user := api.Group("/user", loginRequired)
	user.GET("/album", listAlbumView)
	user.POST("/album", createAlbumView)
	user.GET("/fairy", listFairyView)
	user.POST("/fairy", createFairyView)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", config.Host, config.Port)))
}
