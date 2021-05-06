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
	"fairyla/internal/album"
	"fairyla/vars"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func getUser(c echo.Context) string {
	return c.Get("user").(string)
}

func autoAlbumID(c echo.Context) string {
	albumID := c.Param("id")
	return getAlbumID(getUser(c), albumID)
}

func getAlbumID(user, albumID string) string {
	if albumID == "" {
		return albumID
	}
	if !strings.HasPrefix(albumID, vars.AlbumPreID) {
		// id is name
		albumID = album.AlbumName2ID(user, albumID)
	}
	return albumID
}

type LangQ struct {
	Lang string
	Q    float64
}

func parseAcceptLanguage(acptLang string) []LangQ {
	var lqs []LangQ

	langQStrs := strings.Split(acptLang, ",")
	for _, langQStr := range langQStrs {
		langQ := strings.Split(strings.Trim(langQStr, " "), ";")
		if len(langQ) == 1 {
			lq := LangQ{langQ[0], 1}
			lqs = append(lqs, lq)
		} else {
			qp := strings.Split(langQ[1], "=")
			q, err := strconv.ParseFloat(qp[1], 64)
			if err != nil {
				panic(err)
			}
			lq := LangQ{langQ[0], q}
			lqs = append(lqs, lq)
		}
	}
	return lqs
}
