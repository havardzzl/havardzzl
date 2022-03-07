package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func copyFileContents(src, dst string) (err error) {
	_, err = os.Stat(dst)
	if err != nil {
		fmt.Println("目标文件不存在，跳过")
		return
	}
	in, err := os.Open(src)
	if err != nil {
		fmt.Println("打开源文件失败:" + err.Error())
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		fmt.Println("创建目标文件失败:" + err.Error())
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func main() {
	src := "/opt/app/calico-ipam"
	if os.Getenv("SRC_PATH") != "" {
		src = os.Getenv("SRC_PATH")
	}
	dst := "/opt/cni/bin/calico-ipam"
	if os.Getenv("DST_PATH") != "" {
		dst = os.Getenv("DST_PATH")
	}

	fmt.Println("源文件：" + src + " 目标文件：" + dst)
	for i := 0; i < 10; i++ {
		err := copyFileContents(src, dst)
		fmt.Println("拷贝结果:", err)
		if err == nil {
			break
		}
		fmt.Println("拷贝失败，等待5秒后重试")
		time.Sleep(time.Second * 5)
	}
	time.Sleep(time.Hour * 3650)
}
