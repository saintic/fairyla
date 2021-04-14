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

package auth

import (
	"errors"
	"fairyla/pkg/db"
	"fairyla/pkg/util"
	"fairyla/vars"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const module = "auth"

func Register(c *db.Conn, username, password string) error {
	if !util.IsName(username) {
		return errors.New("invalid name")
	}
	if len(password) < 6 {
		return errors.New("password is too short")
	}
	pwhash, err := util.GeneratePasswordHash(password)
	if err != nil {
		return err
	}
	has, err := c.SIsMember(vars.UserIndex, username)
	if err != nil {
		return err
	}
	if has {
		return errors.New("username already exists")
	}
	pipe := c.Pipeline()
	pipe.SAdd(vars.UserIndex, username)
	pipe.HSet(vars.GenUserKey(username), module, pwhash)
	_, err = pipe.Execute()
	return err
}

// Login Check username & password, if ok, generate a jwt token
func Login(c *db.Conn, username, password string) (token string, err error) {
	if !util.IsName(username) {
		err = errors.New("invalid name")
		return
	}
	has, err := c.SIsMember(vars.UserIndex, username)
	if err != nil {
		return
	}
	if !has {
		err = errors.New("not found username")
		return
	}
	pwhash, err := c.HGet(vars.GenUserKey(username), module)
	if err != nil {
		return
	}
	if !util.CheckPasswordHash(pwhash, password) {
		err = errors.New("wrong password")
		return
	}
	// generate jwt token
	claims := jwt.MapClaims{
		"name": username,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	jt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jt.SignedString([]byte(pwhash))
}

// Verify jwt token
func ParseToken(c *db.Conn, token string) (claims jwt.MapClaims, err error) {
	if strings.Count(token, ".") != 2 {
		err = errors.New("illegal token")
		return
	}

	jt, err := jwt.Parse(token, func(jt *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := jt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jt.Header["alg"])
		}
		claims = jt.Claims.(jwt.MapClaims)
		username := claims["name"].(string)
		has, err := c.SIsMember(vars.UserIndex, username)
		if err != nil {
			return nil, err
		}
		if !has {
			return nil, errors.New("not found username")
		}
		pwhash, err := c.HGet(vars.GenUserKey(username), module)
		if err != nil {
			return nil, err
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(pwhash), nil
	})

	if jt.Valid { // if illegal token, maybe panic
		return claims, nil
	} else {
		return nil, err
	}
}
