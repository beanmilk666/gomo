package gomo

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	LINUX        = "linux"
	WINDOWS      = "windows"
	UTIME_INDEX  = 13 //任务在用户态运行时间的下标
	STIME_INDEX  = 14 //任务在核心态运行时间的下标
	CUTIME_INDEX = 15 //任务的所有已死线程在用户态运行时间下标
	CSTIME_INDEX = 16 //任务的所有已死线程在核心态运行时间下标
)

/**
 * 根据pid获取进程使用时间，单位为jiffies
 * 进程的总CPU时间:time=utime+stime+cutime+sutime（详细解释见常亮定义）
 */
func getProcessCpuTime(pid string) int {
	fileName := fmt.Sprintf("/proc/%s/stat", pid) //进程的状态信息文件（包含进程CPU使用时间）
	file, err := os.Open(fileName)
	defer file.Close()
	if nil != err {
		fmt.Println("Open file err:", err)
		return 0
	}
	con, err := ioutil.ReadAll(file)
	if nil != err {
		fmt.Println("Read file err:", err)
		return 0
	}
	//[]byte直接转化成string会在最后加上一个换行
	cons := strings.Replace(string(con), "\n", "", 1)
	rsSlice := strings.Split(cons, " ")
	uTime, _ := strconv.Atoi(rsSlice[UTIME_INDEX])
	sTime, _ := strconv.Atoi(rsSlice[STIME_INDEX])
	cuTime, _ := strconv.Atoi(rsSlice[CUTIME_INDEX])
	csTime, _ := strconv.Atoi(rsSlice[CSTIME_INDEX])
	return uTime + sTime + cuTime + csTime
}

/**
 * 获取CPU总使用时间，单位为jiffies
 * [cpu  4852929 335 3980828 1018525134 359190 0 186774 0 0 0]总长度是12个（cpu后有一个空格）
 */
func getTotalCpuTime() int {
	file, err := os.Open("/proc/stat")
	defer file.Close()
	if nil != err {
		fmt.Println("Open file err:", err)
		return 0
	}
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n') //读取第一行，总CPU耗时信息在第一行
	line = strings.Replace(line, "\n", "", 1)
	if err != nil {
		if err == io.EOF {
			return 0
		}
		return 0
	}
	rsSlice := strings.Split(line, " ")
	totalTime := 0
	for _, sTime := range rsSlice[2:] {
		t, _ := strconv.Atoi(sTime)
		totalTime += t
	}
	return totalTime
}
