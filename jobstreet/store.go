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

// var fn = fmt.Sprintf("jobstreet_%s.csv", time.Now().Format("20060102"))
var fn = "jobstreet_20220803.csv"

func insert() {
	// 初始化数据库
	mydb := db.InitMySQL()
	query := `INSERT INTO jobs_jobstreet(ext_id,page_url,logo_url,salary,job_title,company_name,posted_at,company_website,company_size,
		avg_process_time,company_overview,company_telephone_number,company_nearby_locations,job_description,job_career_leveln,
		job_years_of_experience,job_qualification,job_field_of_study,job_skills,job_employment_type,job_languages,job_closing_date,
		job_function_value,category1,category2,job_benefits,location,location_str,source_country) 
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
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
