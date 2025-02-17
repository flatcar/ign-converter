// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"fmt"
	"net/url"

	"github.com/flatcar-linux/ignition/config/shared/errors"
	"github.com/flatcar-linux/ignition/config/validate/report"
)

func (f File) Validate() report.Report {
	if f.Overwrite != nil && *f.Overwrite && f.Append {
		return report.ReportFromError(errors.ErrAppendAndOverwrite, report.EntryError)
	}
	return report.Report{}
}

func (f File) ValidateMode() report.Report {
	r := report.Report{}
	if err := validateMode(f.Mode); err != nil {
		r.Add(report.Entry{
			Message: err.Error(),
			Kind:    report.EntryError,
		})
	}
	if f.Mode == nil {
		r.Add(report.Entry{
			Message: errors.ErrPermissionsUnset.Error(),
			Kind:    report.EntryWarning,
		})
	}
	return r
}

func (fc FileContents) ValidateCompression() report.Report {
	r := report.Report{}
	switch fc.Compression {
	case "", "gzip":
	default:
		r.Add(report.Entry{
			Message: errors.ErrCompressionInvalid.Error(),
			Kind:    report.EntryError,
		})
	}
	return r
}

func (fc FileContents) ValidateSource() report.Report {
	r := report.Report{}
	err := validateURL(fc.Source)
	if err != nil {
		r.Add(report.Entry{
			Message: fmt.Sprintf("invalid url %q: %v", fc.Source, err),
			Kind:    report.EntryError,
		})
	}
	return r
}

func (fc FileContents) ValidateHTTPHeaders() report.Report {
	r := report.Report{}

	if len(fc.HTTPHeaders) < 1 {
		return r
	}

	u, err := url.Parse(fc.Source)
	if err != nil {
		r.Add(report.Entry{
			Message: errors.ErrInvalidUrl.Error(),
			Kind:    report.EntryError,
		})
		return r
	}

	switch u.Scheme {
	case "http", "https":
	default:
		r.Add(report.Entry{
			Message: errors.ErrUnsupportedSchemeForHTTPHeaders.Error(),
			Kind:    report.EntryError,
		})
	}

	return r
}
