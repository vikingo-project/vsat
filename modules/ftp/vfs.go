package vsftp

import (
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/vikingo-project/vsat/files"
	"github.com/vikingo-project/vsat/models"
)

type VirtualFS struct {
	nodes sync.Map
}

type vfsNode struct {
	name  string
	dir   bool
	size  int64
	value string // file hash
}

func (n *vfsNode) Name() string {
	return n.name
}

func (m *vfsNode) Size() int64 {
	return m.size
}

func (n *vfsNode) Mode() os.FileMode {
	if n.dir {
		return os.FileMode(0777) | os.ModeDir
	}
	return os.FileMode(0777)
}

func (m *vfsNode) ModTime() time.Time {
	return time.Now()
}

func (n *vfsNode) IsDir() bool {
	return n.dir // todo: m.Mode().IsDir()
}

func (n *vfsNode) Sys() interface{} {
	return nil
}

func (vfs VirtualFS) Clone() VirtualFS {
	var n sync.Map
	vfs.nodes.Range(func(k, v interface{}) bool {
		vm, ok := v.(sync.Map)
		if ok {
			n.Store(k, copyVFS(vm))
		} else {
			n.Store(k, v)
		}
		return true
	})
	return VirtualFS{nodes: n}
}
func fileInCurrentDir(k, path string) bool {
	p := strings.Replace(k, path, "", 1)
	if strings.HasPrefix(p, "/") {
		if len(p) > 1 {
			return false
		}
	}
	if len(p) == 0 {
		return false
	}
	return true
}
func (vfs VirtualFS) ListDir(path string, callback func(FileInfo) error) error {
	vfs.nodes.Range(func(key, value interface{}) bool {
		k := key.(string)
		if strings.HasPrefix(k, path) {
			if fileInCurrentDir(k, path) {
				node := value.(*vfsNode)
				err := callback(node)
				if err != nil {
					return false
				}
			}
		}
		return true
	})
	return nil
}

func (vfs VirtualFS) Stat(path string) (*vfsNode, error) {
	if nodeI, ok := vfs.nodes.Load(path); ok {
		if node, ok := nodeI.(*vfsNode); ok {
			log.Println("Found path in the tree", nodeI)
			return node, nil
		} else {
			return nil, errors.New("system error")
		}
	} else {
		return nil, errors.New("path not found")
	}
}

func (vfs VirtualFS) MkDir(path string) {
	log.Println("MkDir", path)
	if _, ok := vfs.nodes.Load(path); !ok {
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			vfs.nodes.Store(path, &vfsNode{name: parts[len(parts)-1]})
		}
	}
}

func (vfs VirtualFS) Exists(path string) bool {
	if _, ok := vfs.nodes.Load(path); ok {
		return true
	} else {
		return false
	}
}

func (vfs VirtualFS) ReadFile(path string) []byte {
	if nodeI, ok := vfs.nodes.Load(path); ok {
		if node, ok := nodeI.(*vfsNode); ok {
			// todo: get from db
			file := files.GetFileByHash(node.value, true)
			return file.Data
		} else {
			return []byte{}
		}
	}
	return []byte{}
}

func (vfs VirtualFS) RemoveFile(path string) {
	// todo: check for folder
	vfs.nodes.Delete(path)
}

func (vfs VirtualFS) PutFile(path string, f models.File) {
	vfs.nodes.Store(path, &vfsNode{name: f.FileName, size: f.Size})
}
