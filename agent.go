package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"syscall"
	"time"

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
	Threshold   uint64
	SensePeriod string
	Email       Email
}

func (c *Conf) sensePeriod() time.Duration {
	period, err := time.ParseDuration(c.SensePeriod)
	if err != nil {
		log.Fatal(err)
	}
	return period
}

func main() {
	conf, err := conf()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Sensing started.")
	sensePeriodicallyAndReport(conf)
	select {}
	fmt.Println("Sensing finished.")
}

func sensePeriodicallyAndReport(conf Conf) {
	period := conf.sensePeriod()
	fmt.Printf("I'll sense disc space every %s.\n	", period)
	ticker := time.NewTicker(period)
	go func() {
		for _ = range ticker.C {
			fmt.Println("Sensing...")
			senseAndReport(conf)
		}
	}()
}

func senseAndReport(conf Conf) {
	var total, free = discSpaceInfo()
	var message = message(hostname(), conf.Threshold, total, free)
	fmt.Println(message)
	if free < conf.Threshold {
		notify(conf.Email, message)
	}
}

func message(hostname string, threshold uint64, totalSpace uint64, freeSpace uint64) string {
	return fmt.Sprintf(
		"[%s] There're %s available of a total %s. The minimum is %s.\n",
		hostname,
		humanize.Bytes(freeSpace),
		humanize.Bytes(totalSpace),
		humanize.Bytes(threshold),
	)
}

func discSpaceInfo() (uint64, uint64) {
	var stat syscall.Statfs_t
	wd, _ := os.Getwd()

	syscall.Statfs(wd, &stat)
	// Available blocks * size per block = available space in bytes
	return stat.Blocks * uint64(stat.Bsize), stat.Bavail * uint64(stat.Bsize)
}

func hostname() string {
	var n, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func notify(conf Email, message string) {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Eye notification\n\n%s", conf.From, conf.To, message)
	err := sendMail(
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

/* This is a copy of smtp.SendMail, modified to add  InsecureSkipVerify: true
to the configuration. It's only required if something is wrong with the
server's certificate and you want to skip the verification. */
func sendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: addr, InsecureSkipVerify: true}

		if err = c.StartTLS(config); err != nil {
			return err
		}
	}
	if err = c.Auth(a); err != nil {
		return err
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
