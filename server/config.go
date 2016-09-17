// Copyright 2016 Thanh Ngo <felix.infinite@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

// DefaultConfig is server default configuration
var DefaultConfig = Config{
	Domain: "127.0.0.1",
	Port:   5088,
	DBURL:  "localhost",
	DBName: "donit",
}

// Config represents a server configuration structure
type Config struct {
	// DomainName is the domain of this server without the scheme
	Domain string
	Port   int
	// The database url as defined in
	// https://godoc.org/gopkg.in/mgo.v2#Dial
	DBURL string
	// The database name
	DBName string
}
