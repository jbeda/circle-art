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
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

const (
	dpi       = float64(96) // CSS standard. Good enough.
	mmPerInch = float64(25.4)
	mmPerCm   = float64(10)
	ptPerInch = float64(72)
	ptPerPc   = float64(12)
)

var valueRE *regexp.Regexp = regexp.MustCompile(`^([-+]?[0-9]*\.?[0-9]+(?:[eE][-+]?[0-9]+)?)(.*)$`)

// parseValue parses a string value and returns the equivalent in pixels.
func parseValue(s string) (float64, error) {
	caps := valueRE.FindStringSubmatch(s)
	if len(caps) == 0 {
		return 0, fmt.Errorf("Unparsable value: %s", s)
	}
	v, err := strconv.ParseFloat(caps[1], 64)
	if err != nil {
		return 0.0, errors.Wrapf(err, "Error parsing value: %s", s)
	}

	switch u := caps[2]; u {
	case "":
		fallthrough
	case "px":
		return v, nil
	case "in":
		return v * dpi, nil
	case "mm":
		return v * (dpi / mmPerInch), nil
	case "cm":
		return v * (mmPerCm * dpi / mmPerInch), nil
	case "pt":
		return v * (dpi / ptPerInch), nil
	case "pc":
		return v * (ptPerPc * dpi / ptPerInch), nil
	default:
		return 0, fmt.Errorf("Unknown unit: %s", u)
	}
}
