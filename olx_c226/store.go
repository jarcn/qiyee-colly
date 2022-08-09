package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"myclooy/db"
	"os"
	"strconv"
	"time"
)

var fn = fmt.Sprintf("olx_%s.csv", time.Now().Format("20060102"))

// var fn = "olx_20220801_phone.csv"

func insert() {
	// 初始化数据库
	mydb := db.InitMySQL()
	query := `INSERT INTO jobs_olx(ext_id,url,title,name,employment_type,salary_currency,salary_min_value,salary_max_value,
		image,description,address_region,address_locality,address,date_posted,valid_through,images,user_id,user_name,user_created_at,
		highlight,phone) 
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := mydb.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	f, err := os.Open(fn)
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
		row2 := make([]interface{}, 0, len(row))
		for _, v := range row {
			var v2 interface{}
			v2 = v
			row2 = append(row2, v2)
		}
		_, err = stmt.Exec(row2...)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(row[0])
	}
	fmt.Println("insert success")
}

func unixToTime(e string) (datatime time.Time, err error) {
	data, err := strconv.ParseInt(e, 10, 64)
	datatime = time.Unix(data/1000, 0)
	return
}
