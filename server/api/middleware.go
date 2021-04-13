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
	"fairyla/internal/user/auth"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
)

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/v4")
		return next(c)
	}
}

func loginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
        if err := next(c); err != nil {
			c.Error(err)
		}

		fmt.Println("parse jwt")

        fmt.Println(next, next == nil)
        scheme := "Bearer "
		field := c.Request().Header.Get(echo.HeaderAuthorization)
		token := strings.TrimPrefix(field, scheme)
		if !strings.HasPrefix(field, scheme) || token == "" {
			fmt.Println("err1")
			return echo.NewHTTPError(400, "missing or malformed jwt")
		}
		fmt.Println("check jwt")
		_, err := auth.ParseToken(rc, token)
		if err != nil {
			fmt.Println("err2")
			return echo.NewHTTPError(401, "invalid or expired jwt")
		}
		fmt.Println("end jwt")
		return nil
	}
}

func testMD(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        log.Println("into test jwt")
        scheme := "Bearer "
		field := c.Request().Header.Get(echo.HeaderAuthorization)
		token := strings.TrimPrefix(field, scheme)
        log.Println(token)
	    //auth.ParseToken(rc, token)
        //if err != nil {
            //return err
        //}
        ok := check(token)
        if ok {
            log.Println("into next")
            return next(c)
        }

        return errors.New("xxxx")
	}
}

func  check(token string) bool {
	xxxx, err := auth.ParseToken(rc, token)
    log.Println(xxxx)
    return err == nil
}