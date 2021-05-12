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

package util

import (
	"net/url"
	"path/filepath"
	"regexp"
	"time"

	"fairyla/vars"

	"tcw.im/gtc"
)

var (
	namePat = regexp.MustCompile(`^[a-z][0-9a-z\_\-]{1,31}$`)
)

func IsName(name string) bool {
	if name != "" {
		return namePat.MatchString(name)
	}
	return false
}

func Now() int64 {
	return time.Now().Unix()
}

func IsValidURL(toTest string) bool {
	u, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return false
	}
	return true
}

func IsImage(filename string) bool {
	ext := filepath.Ext(filepath.Base(filename))
	return gtc.StrInSlice(ext, vars.AllowImage)
}

func IsVideo(filename string) bool {
	ext := filepath.Ext(filepath.Base(filename))
	return gtc.StrInSlice(ext, vars.AllowVideo)
}
