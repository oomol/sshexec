package main

import (
	"github.com/sirupsen/logrus"
	"sshd/pkg/sshd"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	//
	//executable, err := os.Executable()
	//if err != nil {
	//	logrus.Fatalf("Failed to get executable path: %v", err)
	//}
	//
	//err = os.Chdir(filepath.Dir(executable))
	//if err != nil {
	//	logrus.Fatalf("Failed to change workdir: %v", err)
	//}
	//file, err := os.OpenFile("exec_log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	logrus.Fatalf("Failed to open log file: %v", err)
	//}
	//logrus.SetOutput(file)
}
func main() {
	if err := sshd.SSHExec(); err != nil {
		logrus.Fatalf("Error: %s", err)
	}
}
