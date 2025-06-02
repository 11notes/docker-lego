package main

import (
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"strings"
	"io"
	"bufio"
	"encoding/json"
	"crypto/x509"
	"path/filepath"
	"regexp"
	"encoding/pem"
	"time"
)

type Paths struct {
	Base     string
	BaseVar string
	BaseEtc string
	Lego string
}

type Etc struct {
	Paths Paths
	MaxNameLength int
	Email string
}

var etc = Etc{
	Paths:Paths{
		Base:os.Getenv("APP_ROOT"),
		BaseVar:fmt.Sprintf("%s%s", os.Getenv("APP_ROOT"), "/var"),
		BaseEtc:fmt.Sprintf("%s%s", os.Getenv("APP_ROOT"), "/etc"),
		Lego:"/usr/local/bin/lego",
	},
	MaxNameLength:0,
}

type Config struct {
	Domains []struct {
		Name string `yaml:"name"`
		FQDNs []string `yaml:"fqdns"`
		Commands []string `yaml:"commands"`
		Environment map[string]interface{} `yaml:"environment,omitempty"`
	} `yaml:"domains"`
	Global map[string]interface{} `yaml:"global"`
}

type ACMEAccount struct {
	Email        string `json:"email"`
	Registration struct {
		Body struct {
			Status  string   `json:"status"`
			Contact []string `json:"contact"`
		} `json:"body"`
		URI string `json:"uri"`
	} `json:"registration"`
}

func log(caller string, msg string){
	var spaces string
	for i := 1; i <= etc.MaxNameLength - len(caller); i++ {
		spaces += " "
	}
	fmt.Fprintf(os.Stdout, "%s%s | %s\n", caller, spaces, msg)
}

func config() (*Config, error){
	cfg := &Config{}
	file, err := ioutil.ReadFile(os.Getenv("LEGO_CONFIG"))
	if err != nil{
		if err := yaml.Unmarshal([]byte(os.Getenv("LEGO_CONFIG")), cfg); err != nil {
			return cfg, err
		}
	}else{
		if err := yaml.Unmarshal(file, cfg); err != nil {
			return cfg, err
		}
	}
	for _, domain := range cfg.Domains {
		if len(domain.Name) > etc.MaxNameLength{
			etc.MaxNameLength = len(domain.Name)
		}
	}
	return cfg, nil
}

func path(name string) (error){
	// set paths
	path := fmt.Sprintf("%s/%s", etc.Paths.BaseVar, name)
	symbolicPath := fmt.Sprintf("%s/accounts", path)

	// check if base path for certificates exists, if not create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	// check if symolic link for accounts exists, if not create it
	if _, err := os.Stat(symbolicPath); os.IsNotExist(err) {
		if err := os.Symlink(etc.Paths.BaseEtc, symbolicPath); err != nil {
			return err
		}
	}

	return nil
}

func accountValid(name string, file string) (bool){
	accountJson, err := os.Open(file)
	if err != nil{
		log(name, fmt.Sprintf("lego account error: %s", err))
		return false
	}
	byteValue, err := ioutil.ReadAll(accountJson)
	if err != nil{
		log(name, fmt.Sprintf("lego account error: %s", err))
		return false
	}
	var account ACMEAccount
	err = json.Unmarshal(byteValue, &account);
	if err != nil{
		log(name, fmt.Sprintf("lego account error: %s", err))
		return false
	}
	if account.Registration.Body.Status != "valid"{
		log(name, fmt.Sprintf("lego account error: Account status is not valid, it is %s", account.Registration.Body.Status))
		return false
	}
	return true
}

func run(name string, fqdns []string, commands []string, environment []string) (bool){
	// create paths and symbolic link
	err := path(name)
	if err != nil{
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	
	// setup arguments for lego
	var args = []string{"--accept-tos", "--path", fmt.Sprintf("%s/%s", etc.Paths.BaseVar, name), "--pfx", "--pfx.pass", "lego1234", "--pem"}
	for _, command := range commands {
		args = append(args, command)
	}
	for _, fqdn := range fqdns {
		args = append(args, "--domains")
		args = append(args, fqdn)
	}

	// check if an account is already registered for the email address and set command accordingly
	accountFile := fmt.Sprintf("%s/acme-v02.api.letsencrypt.org/%s/account.json", etc.Paths.BaseEtc, etc.Email)
	if _, err := os.Stat(accountFile); os.IsNotExist(err){
		args = append(args, "run")
		log(name, "create certificate")
	}else{
		if accountValid(name, accountFile) {
			// check if a certificate already exists
			if _, err := os.Stat(fmt.Sprintf("%s/%s/certificates/%s.crt", etc.Paths.BaseVar, name, strings.Replace(fqdns[0], "*", "_", -1))); os.IsNotExist(err) {
				args = append(args, "run")
				log(name, "create certificate")
			}else{
				args = append(args, "renew")
				log(name, "renew certificate")
			}
		}else{
			return false
		}
	}

	// prepare command with environment and pipes
	cmd := exec.Command(etc.Paths.Lego, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid:true}
	cmd.Env = environment
	
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	output := []string{}

	go func() {
		stdoutScanner := bufio.NewScanner(io.MultiReader(stdout,stderr))
		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			output = append(output, line)
			log(name, line)
		}
	}()

	err = cmd.Start()
	if err != nil {
		log(name, fmt.Sprintf("lego execution error: %s\n", err))
		return false
	}

	err = cmd.Wait()
	if err != nil {
		log(name, fmt.Sprintf("lego execution error: %s\n", err))
		return false
	}

	log(name, "certificate successfully created")

	return true
}

func cleanup(){
	err := filepath.Walk(etc.Paths.BaseVar,
		func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matched, _ := regexp.MatchString("(?i)json", path); matched {
			crt := strings.ReplaceAll(path, ".json", ".crt")
			paths := strings.Split(path, "/")
			root := strings.Join(paths[:len(paths)-2], "/")
			if checkCertificateExpired(crt){
				os.RemoveAll(root)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "filepath.Walk error: %s\n", err)
	}
}

func checkCertificateExpired(crt string) (bool){
	file, err := ioutil.ReadFile(crt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open certificate %s: %s\n", crt, err)
	}
	b, _ := pem.Decode(file)
	certificates, err := x509.ParseCertificates(b.Bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse certificate %s: %s\n", crt, err)
	}
	for _, certificate := range certificates {
		if time.Now().Unix() > certificate.NotAfter.Unix(){
			return true
		}
	}
	return false
}

func daily(){
	fmt.Println("starting daily job")

	// check if config is valid
	cfg, err := config()
	if err != nil{
		fmt.Fprintf(os.Stderr, "config error: %s\n", err)
	}else{
		fmt.Printf("found %b entires in config file\n", len(cfg.Domains))
		for _, certificate := range cfg.Domains {
			// create env for lego to use
			var env []string
			if(len(certificate.Environment) > 0){
				for key, value := range certificate.Environment {
					env = append(env, fmt.Sprintf("%s=%v", key, value))
				}
			}
			if(len(cfg.Global) > 0){
				for key, value := range cfg.Global {
					env = append(env, fmt.Sprintf("%s=%v", key, value))
					if(key == "LEGO_EMAIL"){
						etc.Email = fmt.Sprintf("%v", value)
					}
				}
			}

			// run
			go run(certificate.Name, certificate.FQDNs, certificate.Commands, env)
		}
	}

	// clean out old certificates
	cleanup()
}

func main(){
	// syscalls
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)

	// event listener
	go func() {
		<- signalChannel
		os.Exit(0)
	}()

	// set schedule
	daily()
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cron error: %s\n", err)
	}
	_, err = scheduler.NewJob(gocron.CronJob("0 9 * * *", false), gocron.NewTask(daily))
	if err != nil {
		fmt.Fprintf(os.Stderr, "cron error: %s\n", err)
	}
	scheduler.Start()
	select {}
}