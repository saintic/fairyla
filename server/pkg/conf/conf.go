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

package conf

import (
	"errors"

	"gopkg.in/ini.v1"
	"tcw.im/ufc"
)

type Config struct {
	path string
	obj  *ini.File
}

func New(cfgfile string) (cfg Config, err error) {
	if !ufc.IsCommonFile(cfgfile) {
		err = errors.New("not found config file")
		return
	}
	obj, err := ini.Load(cfgfile)
	if err != nil {
		return
	}
	return Config{cfgfile, obj}, nil
}

func (c Config) getKey(section, key string) string {
	return c.obj.Section(section).Key(key).String()
}

func (c Config) mustKey(section, key, defaults string) string {
	v := c.getKey(section, key)
	if v == "" {
		v = defaults
	}
	return v
}

// GetKey 获取默认分区下某个键的值
func (c Config) GetKey(key string) string {
	return c.getKey(ini.DefaultSection, key)
}

// MustKey 获取默认分区下某个键的值，可设置默认值
func (c Config) MustKey(key, defaults string) string {
	return c.mustKey(ini.DefaultSection, key, defaults)
}

// GetSecKey 获取某个分区下某个键的值
func (c Config) GetSecKey(section, key string) string {
	return c.getKey(section, key)
}

// MustSecKey 获取某个分区下某个键的值，可设置默认值
func (c Config) MustSecKey(section, key, defaults string) string {
	return c.mustKey(section, key, defaults)
}
