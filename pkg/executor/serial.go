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

type SerialExecutor struct {
	lock *sync.Mutex
}

type ExecutorResponse struct {
	Body       string `json:"code,omitempty"`
	RunOutput  string `json:"output"`
	ExecTimeUS int    `json:"runtime_us"`
}

func NewSerialExecutor() *SerialExecutor {

	s := SerialExecutor{}
	return &s
}

func (s *SerialExecutor) makeTempFile(body []byte) (string, error) {
	hash := md5.New()
	hash.Write(body)
	hashed := hash.Sum(nil)
	name := fmt.Sprintf("%x", hashed)
	logrus.Info(name)
	codeDir := strings.TrimRight(viper.GetString("codedir"), "/")
	f, e := os.Create(fmt.Sprintf("%s/%s.go", codeDir, name))

	if e != nil {
		logrus.Error(e)
		return "", e
	}

	_, e = f.Write(body)
	if e != nil {
		return "", e
	}

	return f.Name(), nil

}

func (s *SerialExecutor) Format(body []byte) (*ExecutorResponse, error) {
	f, e := s.makeTempFile(body)
	if e != nil {
		return nil, e
	}
	logrus.Info(f)

	cmd := viper.GetString("gobinlocation")
	usr := viper.GetString("runuser")
	s.lock.Lock()
	start := time.Now()
	out, e := exec.Command("sudo", "-u", usr, cmd, "fmt", f).CombinedOutput()
	end := time.Now()
	s.lock.Unlock()
	if e != nil {
		logrus.Info(e.Error)
		logrus.Info(string(out))
		return nil, fmt.Errorf("%s - %s", e.Error(), string(out))
	}
	logrus.Info(string(out))

	formatted, e := ioutil.ReadFile(f)
	if e != nil {
		return nil, e
	}

	response := ExecutorResponse{
		Body:       string(formatted),
		RunOutput:  string(out),
		ExecTimeUS: end.Nanosecond() - start.Nanosecond(),
	}

	return &response, nil
}

func (s *SerialExecutor) Run(body []byte) (*ExecutorResponse, error) {
	f, e := s.makeTempFile(body)
	if e != nil {
		return nil, e
	}
	logrus.Info(f)

	cmd := viper.GetString("gobinlocation")
	usr := viper.GetString("runuser")
	s.lock.Lock()
	start := time.Now()
	out, e := exec.Command("sudo", "-u", usr, cmd, "run", f).CombinedOutput()
	end := time.Now()
	s.lock.Unlock()
	if e != nil {
		logrus.Info(e.Error)
		logrus.Info(string(out))
		return nil, fmt.Errorf("%s - %s", e.Error(), string(out))
	}
	response := ExecutorResponse{
		RunOutput:  string(out),
		ExecTimeUS: end.Nanosecond() - start.Nanosecond(),
	}

	return &response, nil
}
