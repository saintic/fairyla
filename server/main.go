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
	"flag"
	"fmt"
	"log"

	"fairyla/api"
	"fairyla/internal/sys"
)

// fairy server version
const version = "0.1.0"

var (
	v bool
	s bool

	host        string
	port        uint
	rawurl      string
	sapic_url   string
	sapic_token string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&s, "print-config", false, "show config info and exit")

	flag.StringVar(&host, "host", "0.0.0.0", "http listen host")
	flag.UintVar(&port, "port", 10210, "http listen port")

	flag.StringVar(&rawurl, "redis", "", "redis url, format: redis://[:<password>@]<host>[:<port>/<db>]")

	flag.StringVar(&sapic_url, "sapic-url", "", "Sapic Base URL (API & SDK)")
	flag.StringVar(&sapic_token, "sapic-token", "", "Sapic Api LinkToken")
}

func main() {
	flag.Parse()
	config := sys.New(rawurl, host, port, sapic_url, sapic_token)
	if v {
		fmt.Println(version)
	} else if s {
		fmt.Println(config)
	} else {
		config.Check()
		api.StartApi(config)
	}
}
