package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/dustin/go-humanize"
)

type Email struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       string
}

type Conf struct {
	Threshold uint64
	Email     Email
}

func main() {
	conf, err := conf()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var freeSpace = availableSpace()
	var message = message(conf.Threshold, freeSpace)
	fmt.Println(message)
	if freeSpace < conf.Threshold {
		notify(conf.Email, message)
	}
}

func message(threshold uint64, freeSpace uint64) string {
	return fmt.Sprintf(
		"There're %s available of a minimum %s.\n",
		humanize.Bytes(freeSpace),
		humanize.Bytes(threshold),
	)
}

func availableSpace() uint64 {
	var stat syscall.Statfs_t
	wd, _ := os.Getwd()

	syscall.Statfs(wd, &stat)
	// Available blocks * size per block = available space in bytes
	return stat.Bavail * uint64(stat.Bsize)
}

func notify(conf Email, message string) {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Eye notification\n%s", conf.From, conf.To, message)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		smtp.PlainAuth("", conf.Username, conf.Password, conf.Host),
		conf.From,
		[]string{conf.To},
		[]byte(msg),
	)
	if err != nil {
		log.Printf("SMTP error: %s", err)
	}
}

func conf() (Conf, error) {
	configFile := "/etc/eye/eye.cfg"
	var err error
	_, err = os.Stat(configFile)
	if err != nil {
		log.Fatal("Config file is missing: ", configFile)
	}
	var conf Conf
	_, err = toml.DecodeFile(configFile, &conf)
	return conf, err
}
