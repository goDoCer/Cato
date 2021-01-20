package cate

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/kardianos/osext"
	"golang.org/x/crypto/ssh/terminal"
)

var path string

func init() {
	p, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
	path = strings.ReplaceAll(p, "\\", "/")
}

const (
	cateURL      = "https://cate.doc.ic.ac.uk"
	timeTableURL = cateURL + "/timetable.cgi?keyt=%d:%d:%s:%s"
)

var auth string

//GetLoginDetails asks the user for their login details
func GetLoginDetails() {
	if err := getAuth(); err != nil {
		fmt.Println("Enter your shortcode")
		reader := bufio.NewReader(os.Stdin)
		shortcode, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}

		fmt.Println("Enter your password")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		str := base64.StdEncoding.EncodeToString([]byte(string(shortcode) + ":" + string(bytePassword)))
		auth = "Basic" + str
	}
}

func login(url string) (*http.Response, error) {
	GetLoginDetails()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", auth)
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func getAuth() error {
	var s map[string]string
	file, err := ioutil.ReadFile(path + "/" + "secrets.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		return err
	}
	auth = "Basic " + s["Auth"]
	return nil
}
