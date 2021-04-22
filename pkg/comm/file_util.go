package comm

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
)

// CreateDir 创建目录, 不存在则创建
func CreateDir(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		cmd := exec.Command("mkdir", "-p", dir)
		err = cmd.Run()
	}
	return
}

// WriteFile 如果存在文件，比较不同则写入
func WriteFile(name string, content []byte) (same bool, err error) {
	dirName := path.Dir(name)
	if err = CreateDir(dirName); err != nil {
		err = errors.Wrap(err, "createDir")
		return
	}

	if _, err = os.Stat(name); os.IsNotExist(err) {
		err = errors.Wrap(ioutil.WriteFile(name, content, 0644), "write_file")
		return
	}

	old, _ := ioutil.ReadFile(name)
	if !bytes.Equal(content, old) {
		err = ioutil.WriteFile(name, content, 0644)
	} else {
		same = true
	}
	return
}
