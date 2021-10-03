// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vsftp

import (
	"bytes"
	"fmt"
	"log"

	"github.com/vikingo-project/vsat/utils"
)

type listFormatter []FileInfo

// Short returns a string that lists the collection of files by name only,
// one per line
func (formatter listFormatter) Short() []byte {
	var buf bytes.Buffer
	log.Println("im short", formatter)
	for _, file := range formatter {
		utils.PrintDebug("format file... %+v", file)
		fmt.Fprintf(&buf, "%s\r\n", file.Name())
	}
	return buf.Bytes()
}
