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
	"strings"

	"fairyla/internal/user/auth"

	"github.com/labstack/echo/v4"
)

var (
	ErrJWTMissing = echo.NewHTTPError(400, "missing or malformed jwt")
	ErrJWTInvalid = echo.NewHTTPError(401, "invalid or expired jwt")
)

func getJWT(c echo.Context) (token string, err error) {
	scheme := "Bearer "
	field := c.Request().Header.Get(echo.HeaderAuthorization)
	token = strings.TrimPrefix(field, scheme)
	if !strings.HasPrefix(field, scheme) || token == "" {
		err = ErrJWTMissing
		return
	}
	return
}

func checkJWT(c echo.Context) (user string, err error) {
	token, err := getJWT(c)
	if err != nil {
		return
	}
	claims, err := auth.ParseToken(rc, token)
	if err != nil {
		err = ErrJWTInvalid
		return
	}
	user = claims["name"].(string)
	return
}

// API登录拦截器
func loginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := checkJWT(c)
		if err != nil {
			if err == ErrJWTMissing || err == ErrJWTInvalid {
				return err
			}
			return echo.NewHTTPError(400, err.Error())
		}
		c.Set("user", user)
		return next(c)
	}
}
