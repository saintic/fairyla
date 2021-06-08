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
	api.Match([]string{"GET", "HEAD"}, "/ready", readyView)
	api.GET("/config", configView)
	api.GET("/album", listPublicAlbumView)
	api.GET("/album/:owner/:id", getPublicAlbumView)

	auth := api.Group("/auth")
	auth.POST("/signup", signUpView)
	auth.POST("/signin", signInView)
	auth.POST("/forgot", forgotView)
	auth.POST("/reset_passwd", resetPasswdView)

	user := api.Group("/user", loginRequired) // 用户接口，需登录

	user.GET("/", getUserView)
	user.PUT("/profile", updateUserProfileView)
	user.PUT("/passwd", updateUserPasswdView)
	user.PUT("/setting", updateUserSettingView)

	user.POST("/upload", uploadView)

	user.POST("/album", createAlbumView)
	user.PUT("/album/:id", updateAlbumView)
	user.DELETE("/album/:id", dropAlbumView)
	user.GET("/album", listAlbumView)
	user.GET("/album/names", listAlbumNamesView)
	user.GET("/album/:id", getAlbumView)
	user.GET("/album/:id/fairy", listFariesOfAlbumView)

	user.POST("/fairy", createFairyView)
	user.DELETE("/fairy/:id", dropFairyView)

	user.POST("/claim", createClaimView)
	user.GET("/claim", listClaimView)
	user.GET("/claim/:owner/:id", getClaimView)

	user.POST("/_event", createEventView)
	user.GET("/event", listEventView)
	user.DELETE("/event/:id", dropEventView)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)))
}
