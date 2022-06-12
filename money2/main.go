/*


 */

package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func add_one_day(date string) string {
	//　うるう年、月終り、完全無視
	arr := strings.Split(date, "/")

	year, _ := strconv.Atoi(arr[0])
	month, _ := strconv.Atoi(arr[1])
	day, _ := strconv.Atoi(arr[2])

	day += 1

	if day >= 31 {
		day = 1
		month += 1
	}
	if month > 12 {
		month = 1
		year += 1
	}

	day_string := strconv.Itoa(year) + "/" + strconv.Itoa(month) + "/" + strconv.Itoa(day)

	return day_string

	/*
		for _, s := range arr1 {
			fmt.Printf("%s\n", s)
			year, _ := strconv.Atoi(s)
		}
	*/

}

func main() {

	toshiba_db, err := sql.Open("postgres", "user=postgres password=password host=localhost port=15432 dbname=toshiba sslmode=disable")
	if err != nil {
		log.Fatal("OpenError: ", err)
	}
	defer toshiba_db.Close()

	add_day := "2021/6/2"

	cash := 1000000.00     // 100万円からスタート
	potential_asset := 0.0 //潜在的資金
	has_toshiba := 0       // 持っていない:0 持っている:1
	old_close := 0.0       //

	for i := 0; i < 365; i++ {
		fmt.Println("count=", i)

		add_day = add_one_day(add_day) // 日付が一日加算される
		fmt.Println(add_day)

		select_string := "SELECT date, close from stock where date = '" + add_day + "'"

		rows, err := toshiba_db.Query(select_string)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var date string
		var close float64

		for rows.Next() {
			if err := rows.Scan(&date, &close); err != nil {
				fmt.Println(err)
			}
		}
		//fmt.Println(date, close)

		if close == 0 { // 市場がクローズしている
			continue
		}

		// 資産調整(これがキモか)
		if has_toshiba == 1 {
			potential_asset += close - old_close
		}

		// 本日の江端の行動(昨日より値上がりしていれば買うし、値下がりしていれば売る、という単純な行動)
		if old_close < close && has_toshiba == 0 {
			has_toshiba = 1 // 買い
			cash -= close   // 現金は減って
			//potential_asset += close // 資産は
		} else if old_close > close && has_toshiba == 1 {
			has_toshiba = 0 // 売り
			cash += close
			//potential_asset -= close
		}

		old_close = close

		fmt.Println("potential_asset:", potential_asset, "cash=", cash)
	}

}
