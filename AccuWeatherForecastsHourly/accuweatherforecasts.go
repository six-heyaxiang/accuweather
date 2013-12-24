// AccuWeatherForecastsHourly project main.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	. "net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//配置文件名称
const setting_file_name = "./settings.properties"

//go程数量
var complicate_count int = 0

//日志文件
var logger *log.Logger = nil
var savePath string = ""
var logFileName string = ""

//抓取数据保存路径
var dataSavePath string = ""

//apikey
var apikey string = ""

//城市请求key值
var cityInfo string = ""

//管道
var quit chan int
var end chan int

var taskCount int

//City类
type City struct {
	Id       string
	AccuKey  string
	Response string
}

func main() {

	//读取配置文件
	settings, _ := readSetting(setting_file_name)
	savePath = settings["logPath"]
	dataSavePath = settings["savePath"]
	complicate_count, _ = strconv.Atoi(settings["complicateNum"])
	apikey = settings["apiKey"]
	logFileName = settings["logFileName"]
	cityInfo = settings["cityInfo"]

	//配置日志保存文件
	t := time.Now()
	logger, _ = setLoggerSaveFile(savePath, savePath+logFileName+"."+strconv.Itoa(t.Year())+"-"+strconv.Itoa(int(t.Month()))+"-"+strconv.Itoa(t.Day()))

	logger.Println("核心数：" + strconv.Itoa(runtime.NumCPU()) + "协程数：" + strconv.Itoa(complicate_count))
	//设置核心数
	runtime.GOMAXPROCS(runtime.NumCPU())
	cities, _ := readFileArray(cityInfo)
	taskCount = len(cities)
	city := make(chan City, complicate_count*3)
	end = make(chan int)
	quit = make(chan int)
	result := make(chan City, complicate_count)
	go writeResponseToFile(result)
	go writeCitiesToChannel(city, cities)
	for i := 0; i < complicate_count; i++ {
		go startRequest(city, result, quit)
	}
	for {
		if <-end > 0 {
			logger.Println("任务执行完成一次")
		}
		time.Sleep(10 * 1000)
		logger.Println("请求任务再次执行")
		go writeCitiesToChannel(city, cities)
	}
}
func writeCitiesToChannel(city chan City, cities []City) {
	for i := 0; i < len(cities); i++ {
		city <- cities[i]
	}
	logger.Println("城市信息写入channel完成")
}
func writeResponseToFile(result chan City) {
	count := 0
	os.MkdirAll(dataSavePath, 0700)
	for {
		city := <-result
		count++
		if count == taskCount {
			quit <- 1
			count = 0
			continue
		}
		if len(city.Response) != 0 {
			path := dataSavePath + city.Id + ".json"
			file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0660)
			if nil != err {
				logger.Println(city.Id + ".json保存失败！")
				continue
			} else {
				_, err := file.Write([]byte(city.Response))
				if nil != err {
					logger.Println(city.Id + ".json 写入失败！")
				}
			}
			file.Close()
		}
	}
	logger.Println("文件写入go程结束")
}

//发送http请求
func startRequest(ch chan City, result chan City, quit chan int) {
	for {
		var city City
		quitCount := 0
		select {
		case city = <-ch:
			if len(city.Id) == 0 || len(city.AccuKey) == 0 {
				continue
			}
			request, _ := NewRequest("GET", "http://apidev.accuweather.com/forecasts/v1/hourly/24hour/"+city.AccuKey+".json?apiKey="+apikey+"&language=en&details=true", nil)
			//resp, err := http.Get(accuweather.GetForecastHourly(params.ENVRIOMENT_DEVELOPMENT, params.HOURLY_HOUR_24, params.VERSION_1, params.FORMAT_JSON, city.AccuKey, api_key, params.LANGUAGE_ENGLISH, params.DETAILS_FALSE))
			//resp, err := http.Get("http://apidev.accuweather.com/forecasts/v1/hourly/24hour/" + city.AccuKey + ".json?apiKey=" + api_key + "&language=en&details=true")
			resp, err := DefaultClient.Do(request)
			if nil != err {
				logger.Println("城市：" + city.Id + "请求失败：" + city.AccuKey)
				ch <- city
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if nil != err {
				logger.Println("获取响应内容失败！")
				ch <- city
				continue
			}
			resp.Body.Close()
			city.Response = string(body)
			result <- city
		case quitCount = <-quit:
			if quitCount == taskCount {
				end <- 1
			} else {
				quit <- (quitCount + 1)
			}
		}
	}
}

//设置日志保存路径和文件文件名
func setLoggerSaveFile(filePath string, fileName string) (loger *log.Logger, err error) {
	dirErr := os.MkdirAll(filePath, 0700)
	if dirErr != nil {
		fmt.Println("日志文件目录创建失败！")
		return nil, dirErr
	} else {
		logfile, fileErr := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if fileErr != nil {
			fmt.Println("打开日志保存文件失败！")
			return nil, fileErr
		}
		var logger *log.Logger
		logger = log.New(logfile, "", log.Ldate|log.Ltime)
		return logger, nil
	}
}

//读取配置文件方法
func readSetting(fileName string) (setting map[string]string, err error) {
	//#开头的正则表达式
	reg := regexp.MustCompile(`^#.*`)
	settings := make(map[string]string)
	settingFile, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	settingReader := bufio.NewReader(settingFile)
	for {
		str, _, err := settingReader.ReadLine()
		if err != nil {
			if io.EOF == err {
				break
			} else {
				fmt.Println("读取配置文件错误")
				break
			}
		}
		content := string(str[:])
		if 0 == len(content) || "\r\n" == content || reg.FindAllString(content, -1) != nil {
			continue
		}
		items := strings.Split(strings.TrimSpace(content), "=")
		settings[items[0]] = items[1]
	}
	return settings, nil
}

//读入城市请求key值
func readFileArray(fileName string) (result []City, err error) {
	var cities = make([]City, 0)
	srcFile, err := os.OpenFile(fileName, os.O_RDONLY, 0440)
	if nil != err {
		logger.Println("打开城市信息文件失败")
		return nil, err
	}
	defer srcFile.Close()
	srcReader := bufio.NewReader(srcFile)
	for {
		str, _, err := srcReader.ReadLine()
		if nil != err {
			if io.EOF == err {
				break
			} else {
				logger.Println("读取城市信息文件发生错误")
			}
		}

		content := string(str[:])
		if 0 == len(content) || "\r\n" == content {
			continue
		}
		var city City
		items := strings.Split(strings.TrimSpace(content), ",")
		if len(items) == 2 {
			city.Id = items[0]
			city.AccuKey = items[1]
			cities = append(cities, city)
		}
	}
	return cities, nil
}
