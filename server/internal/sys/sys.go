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
	"fmt"
	"os"
	"strconv"
)

type Setting struct {
	Redis string
	Host  string
	Port  uint
	Sapic Sapic
}

type Sapic struct {
	Api   string
	Token string
}

// New from cli options first
func New(redis, host string, port uint, api, token string) *Setting {
	c := &Setting{
		Redis: redis, Host: host, Port: port,
		Sapic: Sapic{api, token},
	}
	c.parseEnv()
	return c
}

// override option
func (s *Setting) parseEnv() {
	redis := os.Getenv("fairyla_redis_url")
	host := os.Getenv("fairyla_host")
	port := os.Getenv("fairyla_port")
	api := os.Getenv("fairyla_sapic_api")
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
	if api != "" {
		s.Sapic.Api = api
	}
	if token != "" {
		s.Sapic.Token = token
	}
}

func (s *Setting) String() string {
	token := "<No Token>"
	if s.Sapic.Token != "" {
		token = fmt.Sprintf("<%s>", s.Sapic.Token)
	}
	return fmt.Sprintf(
		"Host: %s\nPort: %d\nRedis: %s\nSapic: %s %s",
		s.Host, s.Port, s.Redis, s.Sapic.Api, token,
	)
}

func (s *Setting) Check() {
	sa := s.Sapic
	if s.Redis == "" || sa.Api == "" || sa.Token == "" {
		panic("miss required options")
	}
}

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
