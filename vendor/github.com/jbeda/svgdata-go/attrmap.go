// Copyright 2018 Joe Beda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package svgdata

import (
	"fmt"

	"github.com/pkg/errors"
)

type AttrMap map[string]string

// ExtractValue will pull a value out of the AttrMap. It will convert it to a float and delete it from the map.  If it
// is not in the map 0.0 will be returned.  An error will be returned if this is not a parsable value.
func (am AttrMap) ExtractValue(k string) (float64, error) {
	var v float64
	var err error
	str, ok := am[k]
	if ok {
		v, err = parseValue(str)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		delete(am, k)
		return v, nil
	}
	return v, nil
}

// ExtractValue will pull a value out of the AttrMap. It will convert it to a float and delete it from the map.  If it
// is not in the map an error will be returned.  An error will be returned if this is not a parsable value.
func (am AttrMap) ExtractValueNoDefault(k string) (float64, error) {
	var v float64
	var err error
	str, ok := am[k]
	if !ok {
		return 0, errors.New(fmt.Sprintf("no value for %s specified", k))
	}
	if ok {
		v, err = parseValue(str)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		delete(am, k)
		return v, nil
	}
	return 0, nil
}
