package main

import (
	config "SDK/config_managers"
	"SDK/lib/logger"
	"SDK/lib/processcontrol"
	"flag"
	"fmt"
	"time"

	"SDK/lib/grpc"

	"github.com/sirupsen/logrus"
)

func init() {
	config.LoadConfig("./config/conf.ini")
	logPath := config.GetValue("log", "Path")
	if logPath == "" {
		logPath = "./logs/"
	}
	var logLevel logrus.Level
	logLevelString := config.GetValue("log", "level")
	if logLevelString == "" {
		logLevel = logrus.InfoLevel
	} else {
		switch logLevelString {
		case "Trace":
			logLevel = logrus.TraceLevel
		case "Debug":
			logLevel = logrus.DebugLevel
		case "Info":
			logLevel = logrus.InfoLevel
		case "Warn":
			logLevel = logrus.WarnLevel
		case "Error":
			logLevel = logrus.ErrorLevel
		case "Fatal":
			logLevel = logrus.FatalLevel
		case "Panic":
			logLevel = logrus.PanicLevel
		default:
			logLevel = logrus.InfoLevel
		}
	}
	logger.ConfigLocalFilesystemLogger(logPath, "gcpSdkService", time.Duration(24)*time.Hour, logLevel)
}

var log = logrus.WithFields(logrus.Fields{"package": "main"})

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.String("port", "9999", "The server port")
)

func main() {
	flag.Parse()

	log.Info("sdk start")

	//从配置文件中获取相关参数
	myIp := config.GetValue("myHost", "ip")
	myExamplePort := config.GetValue("myHost", "examplePort")
	mySdkPort := config.GetValue("myHost", "sdkPort")
	certFile := config.GetValue("ssl", "certFile")
	keyFile := config.GetValue("ssl", "keyFile")

	//启动各种处理线程
	go grpc.ExampleStart(myIp, myExamplePort, true, certFile, keyFile)
	go grpc.SdkStart(myIp, mySdkPort, true, certFile, keyFile)

	//阻塞,否则主Go退出， listenner的go将会退出
	exitFlag := <-processcontrol.ProcessExit
	log.Info("sdk exit code: " + fmt.Sprint(exitFlag))
}
