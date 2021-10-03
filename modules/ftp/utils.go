package vsftp

import (
	"bytes"
	"fmt"
	"time"
)

func formatFiles(items []FileInfo) []byte {
	var buf bytes.Buffer
	for _, file := range items {
		fmt.Fprintf(&buf, "%s	1	user	user	%d", file.Mode().String(), file.Size())
		if file.ModTime().Before(time.Now().AddDate(-1, 0, 0)) {
			fmt.Fprint(&buf, file.ModTime().Format("	Jan _2  2006	"))
		} else {
			fmt.Fprint(&buf, file.ModTime().Format("	Jan _2 15:04	"))
		}
		fmt.Fprintf(&buf, "%s\r\n", file.Name())
	}
	b := buf.Bytes()
	return b
}
