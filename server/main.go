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

package main

import (
	"fairyla/api"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// fairy server version
const version = "0.1.0"

var (
	v bool

	host   string
	port   uint
	rawurl string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&v, "version", false, "show version and exit")

	flag.StringVar(&host, "host", "0.0.0.0", "")
	flag.UintVar(&port, "port", 10210, "")

	flag.StringVar(&rawurl, "redis", "", "redis url, format: redis://[:<password>@]<host>[:<port>/<db>]")
}

func main() {
	flag.Parse()
	if v {
		fmt.Println(version)
	} else {
		handle()
	}
}

func handle() {
	if rawurl == "" {
		rawurl = os.Getenv("fairyla_redis_url")
	}
	envhost := os.Getenv("fairyla_host")
	envport := os.Getenv("fairyla_port")
	if envhost != "" {
		host = envhost
	}
	if envport != "" {
		envport, err := strconv.Atoi(envport)
		if err != nil {
			fmt.Println("Invalid environment fairyla_port")
			return
		}
		port = uint(envport)
	}

	api.StartApi(rawurl, host, port)
}
