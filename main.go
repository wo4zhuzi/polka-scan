package main

import (
	"log"
	"os"
	"polka-scan/conf"
	"polka-scan/polka"
	"strconv"
	"syscall"
	"time"
)

func main()  {
	/*防止程序多次执行*/
	lockFile := "./lock.pid"
	lock, err := os.Create(lockFile)
	if err != nil {
		log.Fatal("创建文件锁失败", err)
	}
	defer os.Remove(lockFile)
	defer lock.Close()

	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		log.Println("上一个任务未执行完成，暂停执行")
		os.Exit(1)
	}

	polkadotConf := conf.IniFile.Section("polkadot")
	_concurrent := polkadotConf.Key("concurrent").String()
	concurrent, _ := strconv.Atoi(_concurrent)
	for  {
		polka.RunScan(concurrent)
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}
