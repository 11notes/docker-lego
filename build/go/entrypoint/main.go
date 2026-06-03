package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"io"
	"bufio"
	"regexp"

	"github.com/go-co-op/gocron/v2"
  "github.com/11notes/go-eleven"
)

const APP_BIN = "/usr/local/bin/lego"
const APP_CONFIG_ENV = "LEGO_CONFIG"
const APP_CONFIG string = "/lego/etc/default.yml"

var (
	extractLog = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+(?:Z|[-+]\d{2}:\d{2})\s+level=(\S+) msg=(.+)`)
)

func main(){
	// syscalls
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)
	go func() {
		<- signalChannel
		os.Exit(0)
	}()

	// write env to file if set
	eleven.Container.EnvToFile(APP_CONFIG_ENV, APP_CONFIG)

	// set schedule
	daily()
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		eleven.LogFatal("cron error: %s", err)
	}
	_, err = scheduler.NewJob(gocron.CronJob("0 9 * * *", false), gocron.NewTask(daily))
	if err != nil {
		eleven.LogFatal("cron error: %s", err)
	}
	scheduler.Start()
	select {}
}

func run(bin string, params []string){
	cmd := exec.Command(bin, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid:true}
	cmd.Env = os.Environ()

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	go func() {
		stdoutScanner := bufio.NewScanner(io.MultiReader(stdout,stderr))
		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			if len(line) > 0 {
				match := extractLog.FindStringSubmatch(line)
				eleven.Log(match[1], match[2])
			}
		}
	}()

	err := cmd.Start()
	if err != nil {
		eleven.Log("ERR", err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		eleven.Log("ERR", err.Error())
	}
}

func daily(){
	eleven.Log("INF", "starting certificate renewal process ...")
	run(APP_BIN, []string{"--config", APP_CONFIG})
	eleven.Log("INF", "certificate renewal process complete.")
}