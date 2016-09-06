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

package followers

import "time"

// Follower represents a follower
type Follower struct {
	Username string    `bson:"username" json:"username" valid:"required,alphanum,length(1|30)"`
	Follower string    `bson:"follower" json:"username" valid:"required,alphanum,length(1|30)"`
	FollowAt time.Time `bson:"followAt" json:"followAt" valid:"-"`
}
