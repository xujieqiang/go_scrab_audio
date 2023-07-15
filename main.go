package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/gocolly/colly"
)

var (
	//Bible_mulu []string
	ta        []string = make([]string, 66)
	zhang     []string
	Bible_map map[string]string
	list      []string
)

func c2get(url string) (int, string) {
	surl := strings.Split(url, "/")
	last := surl[len(surl)-1] //4.html  5.html
	nums := strings.Split(last, ".")
	n := nums[0]
	if n == "" {
		return -1, ""
	}

	num, err := strconv.Atoi(n)
	if len(n) == 1 {
		n = "0" + n
	}
	if err != nil {
		fmt.Println("error in string converting to number,", err.Error())
	}

	return num - 1, n
}
func co3(bookname string, bookno string, stonum string, url string) {
	//colly.AllowedDomains("www.godcom.net")
	fmt.Println(bookname, "--------", zhang)
	c3 := colly.NewCollector()
	//audio[id=jquery_jplayer_1]
	c3.OnHTML("div[id=jquery_jplayer_1]", func(e *colly.HTMLElement) {
		fname := bookno + "-" + stonum + ".mp3"
		url := "http://file.kuanye.net:2153/55bible/shouji/" + bookno + "/" + fname
		fmt.Println(url)

		//f := bookno + bookname + stonum + ".mp3"
		c3.Visit(url)

	})
	// c3.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("src", r.URL)
	// })
	c3.OnResponse(func(res *colly.Response) {
		b := res.Request.Body
		fmt.Println(b)
	})
	c3.Visit(url)
}

//http://file.kuanye.net:2153/55bible/shouji/05/05-20.mp3

// http://file.kuanye.net:2153/55bible/shouji/06/06-16.mp3

func co2(bookname string, bookno string, surl string) {
	//colly.AllowedDomains("www.godcom.net")
	// err := os.Mkdir(bookname, os.ModePerm)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	var stonum string
	c2 := colly.NewCollector()
	audiodownload := c2.Clone()
	c2.OnHTML("div[class]", func(e *colly.HTMLElement) {
		cl_z := e.Attr("class")
		list = []string{}
		str := ""
		if cl_z == "zhang1" {
			list = e.ChildAttrs("a", "href")
			str = e.Text
			decode := mahonia.NewDecoder("GBK")
			atext := decode.ConvertString(str)
			tt := strings.TrimSpace(atext)
			text_slice := strings.Split(tt, "\n")
			zhang = text_slice

		}
		//fmt.Println(list)
		for _, val := range list {
			c2.Visit(e.Request.AbsoluteURL(val))
		}
	})
	c2.OnRequest(func(r *colly.Request) {
		//fmt.Println("second::", r.URL)
		s := r.URL.String()[21:]
		stonum = ""
		tag := 0
		//fmt.Println(list)
		for i, v := range list {
			//fmt.Println(v)
			if v == s {
				tag = i
				break
				//fmt.Println(tag)
			}
		}
		//fmt.Println(tag)
		tag += 1
		stonum = strconv.Itoa(tag)
		//fmt.Println(bookname, stonum)
		fname := bookno + "-" + stonum + ".mp3"
		url := "http://file.kuanye.net:2153/55bible/shouji/" + bookno + "/" + fname
		fmt.Println("audio url:", url)
		audiodownload.Visit(url)

		//co3(bookname, bookno, stonum, r.URL.String())
	})

	audiodownload.OnResponse(func(res *colly.Response) {

		_, err := os.Stat("./bible")
		if os.IsNotExist(err) {
			err = os.Mkdir("bible", os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		dirname := "./bible/" + bookno + bookname
		os.Mkdir(dirname, os.ModePerm)
		f := dirname + "/" + bookno + bookname + stonum + ".mp3"
		res.Save(f)
	})

	c2.Visit(surl)
}
func main() {
	fmt.Println("hello")
	//colly.AllowedDomains("www.godcom.net")
	c := colly.NewCollector()
	c.OnHTML("div[class]", func(e *colly.HTMLElement) {

		cc := e.Attr("class")
		// ss := e.ChildAttrs("div[class]", "mulu")
		// fmt.Println(ss)
		//fmt.Println("ccc", cc)
		if cc == "mulu" {
			fmt.Println("======")
			xx := e.ChildAttrs("a", "href")
			str := e.Text
			decode := mahonia.NewDecoder("GBK")
			atext := decode.ConvertString(str)
			tt := strings.TrimSpace(atext)
			text_slice := strings.Split(tt, "\n")
			text_slice = text_slice[1:]
			Bible_mulu := text_slice
			//tag := 0
			k := 0
			for _, val := range Bible_mulu {
				if strings.TrimSpace(val) == "新约" {
					//tag = j
					continue
				}
				//Bible_mulu[k] = val
				ta[k] = val
				k += 1
			}

			//fmt.Println(text_slice[0])
			for _, v := range xx {

				// t := Bible_mulu[i]
				// Bible_map[t] = e.Request.AbsoluteURL(v)

				c.Visit(e.Request.AbsoluteURL(v))
			}
		}
		//fmt.Printf("Link :: %q --> %s", e.Text, cc)
		//c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnRequest(func(r *colly.Request) {
		p, bookno := c2get(r.URL.String())
		mm := ""
		if p != -1 {
			mm = ta[p]

		}
		fmt.Println("#####  visiting url:", mm, "-", bookno, r.URL)
		if mm != "" {
			co2(mm, bookno, r.URL.String())
		}

	})
	c.Visit("http://www.godcom.net/")
}
