package library

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

type db struct {
	DbHost   string
	DbUser   string
	DbPwd    string
	DbPort   int
	DbPrefix string
}

type rd struct {
	Stat        string
	RedisHost   string
	RedisPort   string
	RedisPrefix string
	RedisPwd    string
	RedisDb     int
}

type setting struct {
	ServerName string
	LogFile    string
	PidFile    string
	Daemonize  int
}

type Config struct {
	Db       db
	Redis    rd
	Setting  setting
	ConfPath string
	Server   map[string]string
}

type WsConn struct {
	Fd             int
	Conn           *websocket.Conn
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

//获取当前文件执行路径
func GetConfPath() string {
	execFile, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(execFile)
	index := strings.LastIndex(path, string(os.PathSeparator))
	return path[:index]
}

func GetConf() (Conf Config, err error) {
	path := GetConfPath()
	Conf, err = GetConfInfo(path)
	if err != nil {
		return
	}
	Conf.ConfPath = path
	return
}

func GetConfInfo(path string) (conf Config, err error) {
	if path == "" {
		path = GetConfPath()
	}
	// file, osErr := os.Open(path + "/conf.json")
	// 打开文件
	file, osErr := os.Open("G:/WWW/golang/src/webchat/conf.json")
	if osErr != nil {
		err = errors.New("读取Conf配置文件错误")
		return
	}
	// 关闭文件
	defer file.Close()
	var tmpConf Config
	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(file)
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	errJson := decoder.Decode(&tmpConf)
	if errJson != nil {
		err = errors.New("读取Conf配置错误")
		return
	}
	conf = tmpConf
	return
}

//func main() {
//	conf, err := GetConf()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(conf.Db.DbPort)
//}
