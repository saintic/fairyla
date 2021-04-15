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
	"log"
	"strings"

	"fairyla/internal/user/auth"

	"github.com/labstack/echo/v4"
)

func loginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		scheme := "Bearer "
		field := c.Request().Header.Get(echo.HeaderAuthorization)
		token := strings.TrimPrefix(field, scheme)
		if !strings.HasPrefix(field, scheme) || token == "" {
			return echo.NewHTTPError(400, "missing or malformed token")
		}
		claims, err := auth.ParseToken(rc, token)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(401, "invalid or expired token")
		}
		c.Set("user", claims["name"])
		return next(c)
	}
}
