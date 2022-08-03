package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

// const fileName = "/home/meng/文档/myshare/colly/olx/record_test.csv"

var fileName = fmt.Sprintf("olx_%s.csv", time.Now().Format("20060102"))

// 获取手机号链接
var _phoneUrl = "https://www.olx.co.id/api/users/%s"

var id = 16

func getPhone(ch chan []string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	i := 0
	for {
		i++
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("read row error:%v\n", err)
		}
		if i < 2 {
			continue
		}
		// if len(row) == 19 {
		row = append(row, reqPhone(row[16]))
		// }
		ch <- row
	}
}

var gCurCookieJar *cookiejar.Jar

func init() {
	//var err error;
	gCurCookieJar, _ = cookiejar.New(nil)

}

var first bool

var cookie = `locationPath=[{"id":4000030,"name":"Jakarta Selatan","type":"CITY","longitude":106.80105,"latitude":-6.2609,"parentId":2000007}]; laquesis=pan-59312@a#pan-59740@a#pan-60601@b#pan-67471@b; lqstatus=1658195958; G_ENABLED_IDPS=google; user=j:{"id":"122507225","name":"wang test"}; lf=1; t=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJncmFudFR5cGUiOiJncGx1cyIsImNsaWVudFR5cGUiOiJ3ZWIiLCJ0b2tlblR5cGUiOiJhY2Nlc3NUb2tlbiIsImlzTmV3VXNlciI6ZmFsc2UsImlhdCI6MTY1OTUxMDM0MSwiZXhwIjoxNjU5NTExMjQxLCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjUwNzIyNSIsImp0aSI6IjRjZjQzMjI4ZjY3OTkwMTQ2YmU3MGVjYWMyMDgzNDdhMzYzNmMyMGYifQ.c28Ls8MqcLfxoQZD4SpfzqM-zczRpWuaXuNVkWRo5QzhuE6lucHGHJwPuZQ2O_hu5cdyxDBDHU8a_4yXWAaoviP4QfIoGZRhRi-x8zmDF8PeU2sezaXLs3rbokFDfYmndNVt8fkuPbcOC4MCSvp1qNrwUvOf-CQD3JjsQ8Z84NXB4UgQUOuKKmVofXbw7qxDwWjAoSvhGubqji36rXdn7Dlv_s1GAW_IgjBnrhsqYTApGTAu8WLJzFI9MNMbTfolTztx_QNN0GEgnEWsQTca6Qtft1BMWmU-K6bdrsApbIPDAfOfmeptkpbCmQfS4zLivb-XAX0lkdkO6umxx9duIw; rt=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJncmFudFR5cGUiOiJncGx1cyIsImNsaWVudFR5cGUiOiJ3ZWIiLCJ0b2tlblR5cGUiOiJyZWZyZXNoVG9rZW4iLCJyZWFjdGl2YXRlZCI6ZmFsc2UsImlhdCI6MTY1OTUxMDM0MSwiZXhwIjoxNjY3NTQ1NTQxLCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjUwNzIyNSIsImp0aSI6IjRjZjQzMjI4ZjY3OTkwMTQ2YmU3MGVjYWMyMDgzNDdhMzYzNmMyMGYifQ.Pgf03cbF5CnmW0WCw6KZ4hnVE2bbwK_GZBoErkR3TNhs7y4SQDT7cdP7Gg-oFa0IKTpKhGxAg61AtmJF96YAehQUlhe_vSH_ffted4vWs2UIt6fCg8BSycG3KhWQ9sqYBuADpBq9alv-oV56hmD6hz-KVRm55RzWLypOg13DHPugJjPJn2hI8ZtVoraj9jv36IWXE51EQJbqRLIAtaymxCS7-9b6eTbRY4CyZm8UclsNWAAuDvX-yweu_3HUzOeZ4feRCMh_cwq1OPeLeAhqTiBI_DWO1Bh6Vjhb7t-hTl9JzhmXTFK7T6-ekc_Dkt7KtiXjWo5LKfC7TxowwvAJvQ; ct=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJ0b2tlblR5cGUiOiJjaGF0VG9rZW4iLCJ2ZXJzaW9uIjoiMSIsImlhdCI6MTY1OTUxMDM0MSwiZXhwIjoxNjU5NTk2NzQxLCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjUwNzIyNSIsImp0aSI6IjRjZjQzMjI4ZjY3OTkwMTQ2YmU3MGVjYWMyMDgzNDdhMzYzNmMyMGYifQ.RTNAUDALg8QRuV9V64nxlugZyPVAI3Ef56BytlakfXK1lSz49qmACYTBFL3IvOgxJpAoaIvCA_o6RBiNxBVOkdCLZ0MafmSnqpa6FdThUjg4cNFb17rbtMTdr4LwZoso-XkN5fDFAG74RQjU8V_R7lqWfyhnOSsU0tC6MJXVifIR8lc-qGAJbxVTPWFwWZeAOHkpS_U_tleO9LGrt1CJEBkbArpl_4s9iHVpb32egjzBr_qGQOr2Xa0x1PGwbXgxHOo9xlArSeTmoC9YJMbPeGPbg193j6QtJ6uD1RXcr5Ro-oa3IP1iVnlMk9iXX5T5PVuNivzoeuVpC2WIfvKBbQ; nt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTk1MTAzNDEsImJyYW5kIjoib2x4IiwiY291bnRyeUNvZGUiOiJpZCIsInVzZXJJZCI6MTIyNTA3MjI1fQ.xq1ACc209E9-udvs-vZbo1hwZgGLkQ2eg6dSiIcnGGo; _abck=EBED42F344620BC816C09B5808A2C2E3~0~YAAQz/xkX57kkVuCAQAAi0iGYghWiwdi/QKClwE82Sw69lNjP5pS/K+GwIhltfuSdPAF6Q+LdDaezVDcgvZZ5vAB/JBOID/cqkAce5H03NPFMHpRAfRUp65n2Kloo4Dw7A8vbaLFq5pkirH+KymcNRGmJhsaOpb1+y2ztXLNITNvTY7o4YiaKQBOY7rv7AbA3fdo13oP2MwW3ZhofHe/l5Cin/QoyGYvlfu6R+ls23MhdSWJH0nopzsUUqMtNtrNFBaEY7k3xcSyoHjjr5Dl9rtec/05Kbw+a7zciXSCT6F2BIYm7OMTR1mLkpwZ0VhvVEBCFx+fp/9oq1AlL6QpgY+e1evGo2q3F0kjtPmWsC2zFWrOjYRT/DawD63tHKV1i+g02OZi2i72hox+1WFKc1dI9dXi5A==~-1~-1~-1; ak_bmsc=FC8FC81A57B85E8333409AA2D01A3C39~000000000000000000000000000000~YAAQz/xkX5/kkVuCAQAAi0iGYhAP+u34QWa2e3p+rxFd7J5wiOkVv/2g57e3jdhOwef1lAL/zxTffK06eM+qFKyeojOjY+x/VYkJ2Gv+OzPRN81lGOc7Rm0u67IDwNQ+dSu48Sm3ickPvKK4UnmCdHEWueU9+JxbBnjHKVTpGwCJsf+ytRSBvy82GCq/IJZXHOLqoXB1Gqp8LCFbyWeRiGMiIzNLa1/AQqdUp75PKXtscYEB3OGGkUbiW/rFxovuk9G9zSmhtFpyrZjD+ZVjEr608Djg6TSX7bJXzh3pxWIkT3FG3/c+6ffWOyHgEq7dIGJX+8Os2IKAoruMMijkNyqg9gqD58m5zP/zE0C6Peskjk8k0xKfq6ImR6XTkvnLbd3y63n4dH6ENu0=; bm_sz=F13EDCB9E5A1C82DC8700376DDC738F8~YAAQz/xkX6HkkVuCAQAAi0iGYhDJLLKSXFGC8Dp2f1N3V0/IczqnWFBiDqaatDasSrcxyNBdaw6ndIrL8CjFlS1h8J1/ZrP7Kk7dKoHNYDpBMbIS/VAMrRY6/742igPFLfTi1zWyZuMDPXjRAmvJsbuDZoXtN9/NqLTSd32xhreZi6oyKqk5vZS5rbrq4athmSg2pVGT9MuCUgLkjDrWltdwHvfQR0gkzLetmp3ix8FJSECeguvee1EjKeVEPYK9cjqOOSq63v5J9MSwrKWL76ehwYoxVnPIEIYXFQCew2IGZA==~3555640~4404547; reply_restriction=false; kyc_reply_count=0; bm_mi=B6E0FFA90C90E9810757676FCC9BE801~YAAQz/xkX3TokVuCAQAA+syGYhAHlcRTxvhL2bjNT+TwQ1tQgOBSAAoZmSu+60+0qb4SRYa7yd2yey8GYgfPMZbPDUa5iiKh1sWNNiBqSy1No13EIQa10AjAiy9XnMUysmQy7aSIu3DN0/L+NBkFFpKj3/GPFDRQ8YuWM2OZORMQ7t/WTJ2i+N7HliQ8W1n5nLVcnmZRptGuWxOMWFoEQ04M2ZSOF8Csw8bdZxns+lMCfW/rw818kRtR+UfpgT06xq7SX1IuKai8PVCKzUXcO4L3IKwlZf9Y7vw54YtsYc4lmPVJuzHDSAnp4r9aQgBBwb4IgkFl8nEHTWRwOza9HGRNGS4uGlUBv0rVmMcrPkCM2EeILshy9e5gIl9M~1; bm_sv=583598E7CFEBB368A5BB2D6A9AED3E7B~YAAQz/xkX5LskVuCAQAA4vyGYhC3ntEQIR7ex6LV/yyvE9HPOuMAZm8CdFFgjM0VpshUdGUhTqkjUI/dE5LVmNw07lsNxBKRpdx+rtlnj+etN8YB/Y1jhJWoL8YlAREu0q7t9MJLoyDMBfd8M52TXJztUxIF1/JsjtoDDSjGxPgNtWA1XEW14nsP76D6XSroKIWcYx30FFHTNDzzFMDqhUHV4s8Cq5cj/fC04gTC4xFmj8HTfISa952Z9iSEZje2~1`

func reqPhone(userID string) string {
	client := http.Client{Timeout: time.Second * 20}
	req, err := http.NewRequest("GET", fmt.Sprintf(_phoneUrl, userID), nil) //GET大写

	req.Header.Set("cookie", cookie)
	// req.Header.Set("cookie", os.Getenv("OLX_COOKIE"))

	if err != nil {
		log.Printf("req header err:%+v\n", err)
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("req header err:%+v\n", err)
		return ""
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("resp read err:%+v\n", err)
		return ""
	}

	//全局保存
	fmt.Printf("data:%s\n", data)
	if resp.StatusCode == http.StatusUnauthorized {
		time.Sleep(time.Second * 3)
		panic(err)
	}
	// 提取手机号
	return gjson.Get(string(data), "data.phone").String()
}
