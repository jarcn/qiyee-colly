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

var fileName = fmt.Sprintf("kitalulus_%s.csv", time.Now().Format("20060102"))

func insert() {
	// 初始化数据库
	mydb := db.InitMySQL()
	query := `INSERT INTO jobs_kitalulus(url,ext_id,position_name,company_id,company_name,posted_date,posted_date_str,
		requirement_str,education_level,gender,max_age,min_experience,province,city,type_str,location_site,salary_lower_bound_str,
		salary_upper_bound_str,description,working_day_str,working_hour_str,company_description,contact_weblink) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := mydb.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
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
		row2 := make([]interface{}, 0, len(row))
		for k, v := range row {
			var v2 interface{}
			v2 = v
			if k == 5 {
				v2, err = unixToTime(v)
				if err != nil {
					fmt.Println(err.Error())
					v2 = time.Now()
				}
			}
			row2 = append(row2, v2)
		}
		_, err = stmt.Exec(row2...)
		if err != nil {
			panic(err)
		}
		fmt.Println(row[1])
	}
	fmt.Println("insert success")
}

func unixToTime(e string) (datatime time.Time, err error) {
	data, err := strconv.ParseInt(e, 10, 64)
	datatime = time.Unix(data/1000, 0)
	return
}
