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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const SALT_LENGTH = 8
const SALT_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const DEFAULT_PBKDF2_ITERATIONS = 150000
const DEFAULT_PBKDF2_METHOD = "pbkdf2:sha256"

func genSalt() string {
	var bytes = make([]byte, SALT_LENGTH)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = SALT_CHARS[v%byte(len(SALT_CHARS))]
	}
	return string(bytes)
}

func rendeMethod(iter int) string {
	return fmt.Sprintf("%s:%d", DEFAULT_PBKDF2_METHOD, iter)
}

func hashInternal(password, salt string, iter int) string {
	t := pbkdf2.Key([]byte(password), []byte(salt), iter, 32, sha256.New)
	pwh := hex.EncodeToString(t)
	return fmt.Sprintf("%s$%s$%s", rendeMethod(iter), salt, pwh)
}

func GeneratePasswordHash(password string) string {
	salt := genSalt()
	return hashInternal(password, salt, DEFAULT_PBKDF2_ITERATIONS)
}

func CheckPasswordHash(pwhash string, password string) bool {
	if strings.Count(pwhash, "$") < 2 {
		return false
	}
	pws := strings.Split(pwhash, "$")
	method := pws[0]
	salt := pws[1]
	if !strings.HasPrefix(method, DEFAULT_PBKDF2_METHOD) {
		return false
	}
	if strings.Count(method, ":") < 2 {
		return false
	}
	mds := strings.Split(method, ":")
	iter, err := strconv.Atoi(mds[2])
	if err != nil {
		return false
	}

	return strings.EqualFold(pwhash, hashInternal(password, salt, iter))
}
