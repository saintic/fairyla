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

	"github.com/labstack/echo/v4"
)

var (
	rc  *db.Conn
	cfg *sys.Setting
)

func StartApi(config *sys.Setting) {
	c, err := db.New(config.Redis)
	if err != nil {
		panic(err)
	}
	rc = c
	cfg = config

	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Static("/", cfg.Dir)

	api := e.Group("/api")
	api.GET("/config", configView)
	api.GET("/album", pubAlbumView)

	auth := api.Group("/auth")
	auth.POST("/signup", signUpView)
	auth.POST("/signin", signInView)

	user := api.Group("/user", loginRequired)
	user.POST("/upload", uploadView)
	user.GET("/album", listAlbumView)               // 获取用户所有专辑信息
	user.GET("/album/:id", getAlbumView)            // 获取用户某个专辑信息
	user.GET("/album/:id/fairy", getAlbumFairyView) // 仅获取专辑下所有照片信息
	user.POST("/album", createAlbumView)
	user.DELETE("/album/:id", dropAlbumView)
	user.GET("/fairy", listFairyView) // 获取用户所有照片信息
	user.POST("/fairy", createFairyView)
	user.DELETE("/fairy/:id", dropFairyView)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)))
}
