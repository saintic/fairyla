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
	"strings"

	"fairyla/internal/user/auth"

	"github.com/labstack/echo/v4"
)

func getJWT(c echo.Context) (token string, err error) {
	scheme := "Bearer "
	field := c.Request().Header.Get(echo.HeaderAuthorization)
	token = strings.TrimPrefix(field, scheme)
	if !strings.HasPrefix(field, scheme) || token == "" {
		err = errors.New("missing or malformed token")
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
		err = errors.New("invalid or expired token")
		return
	}
	user = claims["name"].(string)
	return
}

func loginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := checkJWT(c)
		if err != nil {
			return echo.NewHTTPError(400, err.Error())
		}
		c.Set("user", user)
		return next(c)
	}
}
