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
	"strings"

	"github.com/labstack/echo/v4"
)

func getUser(c echo.Context) string {
	return c.Get("user").(string)
}

func autoAlbumID(c echo.Context) string {
	albumID := c.Param("id")
	if albumID == "" {
		return albumID
	}
	if !strings.HasPrefix(albumID, vars.AlbumPreID) {
		// id is name
		albumID = album.AlbumName2ID(getUser(c), albumID)
	}
	return albumID
}
