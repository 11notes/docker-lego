package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"io"
	"bufio"
	"regexp"
	"strings"
	"crypto/x509"
	"path/filepath"
	"time"
	"encoding/pem"

	yaml "gopkg.in/yaml.v3"
	"github.com/go-co-op/gocron/v2"
  "github.com/11notes/go-eleven"
)

const APP_BIN string = "/usr/local/bin/lego"
const APP_CONFIG_ENV string = "LEGO_CONFIG"
const APP_CONFIG string = "/lego/etc/default.yml"
const APP_CRT_ROOT string = "/lego/var/certificates"
const APP_TRAEFIK_ENV string = "TRAEFIK_ROOT"

var (
	extractLog = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+(?:Z|[-+]\d{2}:\d{2})\s+level=(\S+) msg="(.+)`)
)

type TraefikDynamicCertificates struct {
	CertFile string
	KeyFile string
}

type TraefikDynamicTLS struct {
	Certificates []TraefikDynamicCertificates
}

type TraefikDynamic struct {
	TLS TraefikDynamicTLS
}

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

	// run at init
	cron()

	// set schedule
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		eleven.LogFatal("cron error: %s", err)
	}
	_, err = scheduler.NewJob(gocron.CronJob("0 9 * * *", false), gocron.NewTask(cron))
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
				if len(match) >= 3 {
					eleven.Log(match[1], sanitizeLog(match[2]))
				}else{
					eleven.Log("INF", line)
				}
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

func cron(){
	eleven.Log("INF", "starting certificate creation/renewal process ...")
	run(APP_BIN, []string{"--config", APP_CONFIG})
	eleven.Log("INF", "certificate creation/renewal process complete.")
	cleanup()
	if _, ok := os.LookupEnv(APP_TRAEFIK_ENV); ok {
		createTraefikConfiguration()
	}
}

func cleanup(){
	// cleanup expired certificates
	eleven.Log("INF", "running certificate cleanup process ...")
	err := filepath.Walk(APP_CRT_ROOT,
		func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matched, _ := regexp.MatchString("(?i)json", path); matched {
			base := strings.ReplaceAll(path, ".json", "")
			paths := strings.Split(path,`/`)
			files, err := filepath.Glob(base + "*")
			if _, err := os.Stat("/lego/var/traefik"); !os.IsNotExist(err) {
				symbolicFiles, err := filepath.Glob("/lego/var/traefik/" + strings.ReplaceAll(paths[len(paths) - 1], ".json", "") + "*")
				if err == nil {
					files = append(files, symbolicFiles...)
				}
			}
			if err != nil {
				eleven.Log("ERR", "filepath.Glob error: %s", err)
			}else{
				if checkCertificateExpired(base + ".crt"){
					for _, file := range files {
						eleven.Log("WRN", "delete expired certificate file %s", file)
						os.Remove(file)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		eleven.Log("ERR", "filepath.Walk error: %s", err)
	}else{
		eleven.Log("INF", "certificate cleanup process complete.")
	}
}

func checkCertificateExpired(crt string) bool {
	file, err := eleven.Util.ReadFile(crt)
	if err != nil {
		eleven.Log("ERR", "could not open certificate %s: %s", crt, err)
	}
	b, _ := pem.Decode([]byte(file))
	certificates, err := x509.ParseCertificates(b.Bytes)
	if err != nil {
		eleven.Log("ERR", "could not parse certificate %s: %s", crt, err)
	}
	for _, certificate := range certificates {
		if time.Now().Unix() > certificate.NotAfter.Unix(){
			return true
		}
	}
	return false
}

func createTraefikConfiguration(){
	eleven.Log("INF", "starting traefik configuration process ...")
	err := os.MkdirAll("/lego/var/traefik", 0700)
	traefikRoot := os.Getenv(APP_TRAEFIK_ENV)
	if err != nil {
		eleven.Log("ERR", "could os.MkdirAll error: %s", err)
	}else{
		var etcTraefik TraefikDynamic
		err := filepath.Walk(APP_CRT_ROOT,
			func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if matched, _ := regexp.MatchString("(?i)json", path); matched {
				paths := strings.Split(path,`/`)
				base := strings.ReplaceAll(paths[len(paths) - 1], ".json", "")
				eleven.Log("INF", "added %s to Traefik ssl config", base)
				etcTraefik.TLS.Certificates = append(etcTraefik.TLS.Certificates, TraefikDynamicCertificates{
					CertFile:traefikRoot + "/" + base + ".crt",
					KeyFile:traefikRoot + "/" + base + ".key",
				})
				os.Symlink(strings.ReplaceAll(path, ".json", ".crt"), "/lego/var/traefik/" + base + ".crt")
				os.Symlink(strings.ReplaceAll(path, ".json", ".key"), "/lego/var/traefik/" + base + ".key")
			}
			return nil
		})
		if err != nil {
			eleven.Log("ERR", "filepath.Walk error: %s", err)
		}

		traefikConfig, err := yaml.Marshal(&etcTraefik)
		if err != nil {
			eleven.Log("ERR", "yaml.Marshal error: %s", err)
		}else{
			eleven.Util.WriteFile("/lego/var/traefik/lego.yml", string(traefikConfig))
			eleven.Log("INF", "traefik configuration process complete.")
		}
	}
}

func sanitizeLog(log string) string{
	log = strings.ReplaceAll(log, `"`, ``)
	log = strings.ReplaceAll(log, `\n\n`, ` `)
	log = strings.ReplaceAll(log, `\n`, ` `)
	return log
}