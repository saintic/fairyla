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

package sys

import (
	"errors"
	"fairyla/pkg/util"
	"fmt"
	"os"
	"strconv"
	"strings"

	"tcw.im/gtc"
)

type Setting struct {
	Redis string
	Host  string
	Port  uint
	Sapic Sapic
}

type Sapic struct {
	Api   string `json:"api"`
	Field string `json:"field"`
	Token string `json:"token"`
}

var ep = "/api/upload"

func parsePort(sport string) (dport uint, err error) {
	if sport != "" {
		iport, e := strconv.Atoi(sport)
		if e != nil {
			err = e
			return
		}
		dport = uint(iport)
	}
	err = errors.New("empty port")
	return
}

// New from cli options first
func New(host string, port uint, redis, baseURL, token, field string) *Setting {
	if !util.IsValidURL(baseURL) {
		panic("invalid sapic url")
	}
	baseURL = strings.TrimSuffix(strings.TrimSuffix(baseURL, ep), "/")
	c := &Setting{
		Redis: redis, Host: host, Port: port,
		Sapic: Sapic{
			baseURL + ep, field, token,
		},
	}
	c.parseEnv()
	return c
}

// 从环境变量读取配置，优先级高，会覆盖参数
func (s *Setting) parseEnv() {
	redis := os.Getenv("fairyla_redis_url")
	host := os.Getenv("fairyla_host")
	port := os.Getenv("fairyla_port")
	api := os.Getenv("fairyla_sapic_api")     // upload api url
	field := os.Getenv("fairyla_sapic_field") // upload file field name
	token := os.Getenv("fairyla_sapic_token")
	if redis != "" {
		s.Redis = redis
	}
	if host != "" {
		s.Host = host
	}
	dport, err := parsePort(port)
	if err == nil && dport > 1024 {
		s.Port = dport
	}
	if util.IsValidURL(api) && strings.HasSuffix(api, ep) {
		s.Sapic.Api = api
	}
	if field != "" {
		s.Sapic.Field = field
	}
	if token != "" {
		s.Sapic.Token = token
	}
    // other options
    sep := "fairyla_"
    for _, e := rage os.Environ() {
        env := e.split("=")
        if strings.HasPrefix(env[0], sep) {
            field := strings.TrimPrefix(env[0], sep)
            switch field {
            case "icp":
                ...
            case val2:
                ...
            default:
                ...
        }
        }

    }
}

func (s *Setting) String() string {
	token := "<No Token>"
	if s.Sapic.Token != "" {
		token = fmt.Sprintf("<%s>", s.Sapic.Token)
	}
	return fmt.Sprintf(
		"Host: %s\nPort: %d\nRedis: %s\nSapic:\n Api: %s\n Field: %s\n LinkToken: %s",
		s.Host, s.Port, s.Redis, s.Sapic.Api, s.Sapic.Field, token,
	)
}

// 检查是否缺少必须项
func (s *Setting) Check() {
	sa := s.Sapic
	err := ""
	if s.Redis == "" {
		err = "redis"
	} else if sa.Api == "" {
		err = "sapic-url"
	} else if sa.Token == "" {
		err = "sapic-token"
	}
	if err != "" {
		panic(fmt.Sprintf("miss required option: %s\n", err))
	}
}

// 站点公共配置项
func (s *Setting) SitePublic() map[string]interface{} {
	cfg := make(map[string]interface{})
	cfg["sapic"] = s.Sapic
	return cfg
}
