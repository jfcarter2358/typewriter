// routes.go

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"typewriter/config"
	Middleware "typewriter/middleware"

	"github.com/gomarkdown/markdown"

	"github.com/gin-gonic/gin"

	"path/filepath"
	// "log"
)

var contentsRoute string

func initializeRoutes() {
	router.Use(Middleware.CORSMiddleware())

	contentsRoute = "contents"
	err := filepath.Walk("./"+contentsRoute,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasPrefix(path, fmt.Sprintf("%s/.theme", contentsRoute)) || strings.HasPrefix(path, fmt.Sprintf("%s/.git", contentsRoute)) {
				return nil
			}
			if !info.IsDir() {
				extension := filepath.Ext(path)
				htmlPath := strings.TrimPrefix(path[:len(path)-len(extension)], contentsRoute)

				switch extension {
				case ".html":
					log.Println("Adding HTML route to " + htmlPath)
					router.GET("/"+htmlPath, htmlHandler)
				case ".md":
					log.Println("Adding Markdown route to " + htmlPath)
					router.GET("/"+htmlPath, markdownHandler)
				default:
					log.Println("Adding static route to " + htmlPath + extension)
					router.StaticFile("/"+htmlPath+extension, path)
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	for _, proxyItem := range config.Config.ReverseProxies {
		router.Any(fmt.Sprintf("/%s/*path", proxyItem.Prefix), ReverseProxy(proxyItem.Server, "/"+proxyItem.Prefix))
	}

	router.GET("/", RedirectLandingPage)
}

func htmlHandler(c *gin.Context) {
	data, err := ioutil.ReadFile(fmt.Sprintf("./%s%s.html", contentsRoute, c.FullPath()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}
	c.Data(http.StatusOK, "text/html", data)
}

func markdownHandler(c *gin.Context) {
	metadataString := ""
	contents := ""
	metadata := make(map[string]interface{})
	isHeader := true

	// Open file
	f, err := os.Open(fmt.Sprintf("./%s%s.md", contentsRoute, c.FullPath()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	// Read file line by line
	r := bufio.NewReader(f)
	s, e := Readln(r)
	for e == nil {
		if s == "---" {
			isHeader = false
		} else {
			if isHeader {
				metadataString += s + "\n"
			} else {
				contents += s + "\n"
			}
		}
		s, e = Readln(r)
	}

	err = yaml.Unmarshal([]byte(metadataString), metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	markdownHTML := markdown.ToHTML([]byte(contents), nil, nil)

	template_file := "template.html"
	if val, ok := metadata["template"]; ok {
		template_file = val.(string)
	}

	// Load markdown template
	template_data, err := ioutil.ReadFile(fmt.Sprintf("./%s/.theme/%s", contentsRoute, template_file))
	if err != nil {
		panic(err)
	}
	template := string(template_data)

	data := strings.Replace(template, "{{MARKDOWN CONTENTS}}", string(markdownHTML), -1)
	data = strings.Replace(data, "{{TITLE}}", metadata["title"].(string), -1)

	c.Data(http.StatusOK, "text/html", []byte(data))
}

func ReverseProxy(target, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
			req.Header["X-Forwarded-Prefix"] = []string{prefix}
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func RedirectLandingPage(c *gin.Context) {
	c.Redirect(301, config.Config.Redirect)
}
