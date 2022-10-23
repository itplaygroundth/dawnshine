package utils

import (
	//"github.com/amalfra/etag"
	// "github.com/go-redis/redis/v8"
	//"github.com/tidwall/gjson"
	"crypto/md5"
	"github.com/gofiber/fiber/v2"
	"fmt"
	"time"
 
	"log"
	"github.com/valyala/fasthttp"
	//"github.com/labstack/echo/v4"
)

// var redis_master_host = os.Getenv("REDIS_MASTER_HOST")
// var redis_master_port = os.Getenv("REDIS_MASTER_PORT")
// var redis_master_password = os.Getenv("REDIS_MASTER_PASSWORD")
// var redis_slave_host = os.Getenv("REDIS_SLAVE_HOST")
// var redis_slave_port = os.Getenv("REDIS_SLAVE_PORT")
// var redis_slave_password = os.Getenv("REDIS_SLAVE_PASSWORD")
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
 
func FastGet(url string, c *fiber.Ctx) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetUserAgent(fiber.HeaderUserAgent)
	req.Header.SetUserAgent(string(c.Context().UserAgent()))
	req.Header.SetMethod("GET")

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		log.Println("==fastget error==")
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}

func MD5(text string) string {
    data := []byte(text)
    return fmt.Sprintf("%x", md5.Sum(data))
}


func FastPut(url string, data []byte) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("PUT")
	req.SetBody(data)

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}

func FastPost(url string, referrer string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	req.SetRequestURI(url)
	req.Header.Add("Referer", referrer)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}


 

func FastDelete(url string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(url)
	req.Header.SetMethod("DELETE")

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}