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
	api.Match([]string{"GET", "HEAD"}, "/healthz", func(c echo.Context) error {
		ping, err := rc.Ping()
		if err != nil || !ping {
			return c.JSONBlob(503, []byte(`"ERR"`))
		}
		return c.JSONBlob(200, []byte(`"OK"`))
	})
	api.GET("/config", configView)
	api.GET("/album", pubAlbumView)

	auth := api.Group("/auth")
	auth.POST("/signup", signUpView)
	auth.POST("/signin", signInView)

	user := api.Group("/user", loginRequired) // 用户接口，需登录
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

	user.GET("/claim", listClaimView)
	user.POST("/claim", createClaimView)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)))
}
