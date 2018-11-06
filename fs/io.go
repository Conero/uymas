package fs

import (
	"io"
	"io/ioutil"
)

// @Date：   2018/11/6 0006 14:33
// @Author:  Joshua Conero
// @Name:    读写

type FsReaderWriter struct {
	content []byte
	dstFile string
	srcFile string
	errorMsg string
}

func (f *FsReaderWriter) Read(p []byte) (n int, err error) {
	if f.content != nil{
		return len(f.content), nil
	} else if f.srcFile != ""{
		content, err := ioutil.ReadFile(f.srcFile)
		f.content = content
		return len(content), err
	}
	return 0, f
}

func (f *FsReaderWriter) Write(p []byte) (n int, err error)  {
	if f.content != nil{
		if f.dstFile != ""{
			err := ioutil.WriteFile(f.dstFile, f.content, 0755)
			return len(f.content), err
		}else {
			f.errorMsg = "未设置目标文件，文件写入失败！"
			return len(f.content), f
		}
	}
	return 0, nil
}

func (f *FsReaderWriter) Error() string {
	return f.errorMsg
}

// 文件复制
func Copy(dstFile, srcFile string) (bool, error) {
	frw := &FsReaderWriter{
		dstFile: dstFile,
		srcFile: srcFile,
	}
	if _, err := io.Copy(frw, frw); err != nil {
		return false, err
	}
	return true, nil
}

