package executor

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// SerialExecutor allows one remote program to run at a time
type SerialExecutor struct {
	lock *sync.Mutex
}

// Response stores stats about the code run
type Response struct {
	Hash       string `json:"hash"`
	Body       string `json:"code,omitempty"`
	RunOutput  string `json:"output"`
	ExecTimeUS int    `json:"buildandruntime_us"`
}

// NewSerialExecutor returns a serial executor
func NewSerialExecutor() *SerialExecutor {
	var lock = &sync.Mutex{}
	s := SerialExecutor{}
	s.lock = lock
	return &s
}

func (s *SerialExecutor) makeTempFile(body []byte) (string, string, error) {
	hash := md5.New()
	hash.Write(body)
	hashed := hash.Sum(nil)
	nameHash := fmt.Sprintf("%x", hashed)
	logrus.Info(nameHash)
	codeDir := strings.TrimRight(viper.GetString("codedir"), "/")
	f, e := os.Create(fmt.Sprintf("%s/%s.go", codeDir, nameHash))

	if e != nil {
		logrus.Error(e)
		return "", "", e
	}

	_, e = f.Write(body)
	if e != nil {
		return "", "", e
	}

	return f.Name(), nameHash, nil

}

// Load a previouly run code file
func (s *SerialExecutor) Load(hash string) (*Response, error) {
	start := time.Now()
	codeDir := strings.TrimRight(viper.GetString("codedir"), "/")
	source, e := ioutil.ReadFile(fmt.Sprintf("%s/%s.go", codeDir, hash))
	end := time.Now()
	if e != nil {
		return nil, e
	}

	response := Response{
		Hash: hash,
		Body: string(source),
		ExecTimeUS: end.Nanosecond() - start.Nanosecond()
	}

	return &response, nil
}

// Format runs the body through Go Fmt
func (s *SerialExecutor) Format(body []byte) (*Response, error) {
	f, h, e := s.makeTempFile(body)
	if e != nil {
		return nil, e
	}
	logrus.Info(f)

	cmd := viper.GetString("gobinlocation")
	usr := viper.GetString("runuser")
	start := time.Now()
	s.lock.Lock()
	out, _ := exec.Command("sudo", "-u", usr, cmd, "fmt", f).CombinedOutput()
	s.lock.Unlock()
	end := time.Now()
	formatted, e := ioutil.ReadFile(f)
	if e != nil {
		return nil, e
	}

	response := Response{
		Hash:       h,
		Body:       string(formatted),
		RunOutput:  string(out),
		ExecTimeUS: end.Nanosecond() - start.Nanosecond(),
	}

	return &response, nil
}

// Run the code
func (s *SerialExecutor) Run(body []byte) (*Response, error) {
	f, h, e := s.makeTempFile(body)
	if e != nil {
		return nil, e
	}
	logrus.Info(f)

	cmd := viper.GetString("gobinlocation")
	usr := viper.GetString("runuser")
	start := time.Now()
	s.lock.Lock()
	out, _ := exec.Command("sudo", "-u", usr, cmd, "run", f).CombinedOutput()
	s.lock.Unlock()
	end := time.Now()
	response := Response{
		Hash:       h,
		RunOutput:  string(out),
		ExecTimeUS: end.Nanosecond() - start.Nanosecond(),
	}

	return &response, nil
}
