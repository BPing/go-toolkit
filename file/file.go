package file

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

// 下载文件
func DownLoadVideo(httpSrc, dst string) (n int64, err error) {
	os.MkdirAll(path.Dir(dst), os.ModeDir)
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	resp, err := http.Get(httpSrc)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	n, err = io.Copy(out, bytes.NewReader(pix))
	if err != nil {
		DeletePath(dst)
	}
	return
}

// 复制文件
func CopyFile(src, dst string) (n int64, err error) {
	os.MkdirAll(path.Dir(dst), os.ModeDir)
	os.MkdirAll(path.Dir(src), os.ModeDir)
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	pix, err := ioutil.ReadAll(in)
	if err != nil {
		return
	}

	n, err = io.Copy(out, bytes.NewReader(pix))
	if err != nil {
		DeletePath(dst)
	}
	return
}

// 获取本地文件路径
func LocalMapVideo(rootDir, httpUrl string) (localPath string) {
	path := strings.Split(httpUrl, "/")
	if len(path) > 1 {
		localPath = rootDir + "/" + path[len(path)-1]
	}
	return
}

// 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 删除文件或者整个目录，谨慎使用
func DeletePath(path string) error {
	return os.RemoveAll(path)
}
