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

package concreteachieving

import "github.com/iocat/donit/internal/achieving/internal/achievable"

// Achievable rerpesents a concreate achieving.Achievable
type Achievable struct {
	achievable.Achievable `valid:"required"`
}

// HasAchieved implements achieving.Achievable's HadAchieved
func (a *Achievable) HasAchieved() bool {
	return a.Achievable.HasAchieved()
}

// IsRepetitive returns whether this is a repetitive task or not
func (a *Achievable) IsRepetitive() bool {
	return a.Achievable.IsHabit()
}
