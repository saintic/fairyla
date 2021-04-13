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

	"fairyla/internal/db"

	"github.com/labstack/echo/v4"
)

var rc *db.Conn

func StartApi(redis_url, host string, port uint) {
	c, err := db.New(redis_url)
	if err != nil {
		panic(err)
	}
	rc = c

	e := echo.New()
    e.Debug = true
	e.HTTPErrorHandler = customHTTPErrorHandler

	auth := e.Group("/auth")
	auth.POST("/signup", signUpView)
	auth.POST("/signin", signInView)

	test := e.Group("/test", loginRequired)
	test.POST("/check", testView)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}
