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
	"github.com/flatcar-linux/ignition/v2/config/shared/errors"
	"github.com/flatcar-linux/ignition/v2/config/util"

	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
)

func (f File) Validate(c path.ContextPath) (r report.Report) {
	r.Merge(f.Node.Validate(c))
	r.AddOnError(c.Append("mode"), validateMode(f.Mode))
	r.AddOnError(c.Append("overwrite"), f.validateOverwrite())
	return
}

func (f File) validateOverwrite() error {
	if util.IsTrue(f.Overwrite) && f.Contents.Source == nil {
		return errors.ErrOverwriteAndNilSource
	}
	return nil
}

func (f FileEmbedded1) IgnoreDuplicates() map[string]struct{} {
	return map[string]struct{}{
		"Append": {},
	}
}

func (fc FileContents) Validate(c path.ContextPath) (r report.Report) {
	r.AddOnError(c.Append("compression"), fc.validateCompression())
	r.AddOnError(c.Append("verification", "hash"), fc.validateVerification())
	r.AddOnError(c.Append("source"), validateURLNilOK(fc.Source))
	return
}

func (fc FileContents) validateCompression() error {
	if fc.Compression != nil {
		switch *fc.Compression {
		case "", "gzip":
		default:
			return errors.ErrCompressionInvalid
		}
	}
	return nil
}

func (fc FileContents) validateVerification() error {
	if fc.Verification.Hash != nil && fc.Source == nil {
		return errors.ErrVerificationAndNilSource
	}
	return nil
}
