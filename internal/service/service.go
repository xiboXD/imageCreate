package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"image-Designer/internal/config"
	"image-Designer/internal/define"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var Config *define.Config

func init() {
	conf, err := ReadConfig()
	if err != nil {
		log.Fatal("Configuration file does not exist, program exits:", err)
		return
	}
	if len(conf.Cookie) == 0 {
		log.Fatal("Cookie does not exist in the configuration file, program exits:", err)
		return
	}

	if conf.ProxyEnable {
		if len(conf.ProxyUrl) == 0 {
			log.Fatal("Proxy enabled, but proxy address is empty, program exits")
			os.Exit(1)
		}
		transport := &http.Transport{
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse(conf.ProxyUrl)
			},
		}
		config.Client.Transport = transport
	}
	Config = conf
}

func ReadConfig() (*define.Config, error) {
	dir, _ := os.Getwd()
	file, err := ioutil.ReadFile(dir + string(os.PathSeparator) + define.ConfigName)
	if err != nil {
		return nil, err
	}
	conf := new(define.Config)
	err = json.Unmarshal(file, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func Submit(message string) (string, error) {

	escape := url.QueryEscape(message)
	// Process according to your code, submit the request
	requestUrl := fmt.Sprintf(config.RequestUrl, escape)
	request, err := http.NewRequest("POST", requestUrl, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Cookie", Config.Cookie)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Transfer-Encoding", "chunked")
	request.Header.Add("Host", "https://www.bing.com")
	response, err := config.Client.Do(request)
	defer response.Body.Close()
	if err != nil {
		return "", err
	}
	location, err := response.Location()
	if err != nil {
		return "", err
	}
	newRequest, err := http.NewRequest("GET", location.String(), nil)
	newRequest.Header = request.Header
	newRequest.Header.Add("Referer", requestUrl)
	r, err := config.Client.Do(newRequest)
	defer r.Body.Close()
	if err != nil {
		return "", err
	}
	id := location.Query().Get("id")
	q := location.Query().Get("q")
	if len(id) == 0 {
		return "", errors.New("Request failed, please check the network or add a proxy")
	}
	config.Cache[id] = q
	return id, nil
}

func Result(id string) ([]string, error) {
	q := config.Cache[id]
	if len(q) == 0 {
		return nil, errors.New("id does not exist, please check")
	}
	requestUrl := config.RequestResultUrl + id + "?q=" + url.QueryEscape(q)
	request, err := http.NewRequest("POST", requestUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Cookie", Config.Cookie)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	response, err := config.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if len(string(bytes)) == 0 {
		return nil, errors.New("Photo is being generated, please try again later")
	}
	doc, err := html.Parse(strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}
	nodes, err := htmlquery.QueryAll(doc, "//div[@class='img_cont hoff']/img/@src")
	if err != nil {
		return nil, err
	}

	srcArrays := make([]string, 0)
	for _, node := range nodes {
		// Get the value of the node (value of src attribute)
		srcValue := htmlquery.InnerText(node)
		srcValue = strings.ReplaceAll(srcValue, "w=270&h=270&c=6&r=0&o=5&dpr=1.5", "")
		srcArrays = append(srcArrays, srcValue)
	}

	if len(srcArrays) == 0 {
		nodes, err = htmlquery.QueryAll(doc, "//*[@id='gir_async']/a/img/@src")
		for _, node := range nodes {
			// Get the value of the node (value of src attribute)
			srcValue := htmlquery.InnerText(node)
			srcValue = strings.ReplaceAll(srcValue, "w=270&h=270&c=6&r=0&o=5&dpr=1.5", "")
			srcArrays = append(srcArrays, srcValue)
		}
	}

	if len(srcArrays) == 0 {
		return srcArrays, nil
	}
	//delete(config.Cache, id)
	return srcArrays, nil
}
