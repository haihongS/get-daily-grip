package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", httpIndex)
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("./data"))))
	http.HandleFunc("/login", login)
	go delFile()
	http.ListenAndServe(":6688", nil)
}

func httpIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`hell`))
	r.ParseForm()
	fmt.Println(r)
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, err := template.ParseFiles("html/login.html")
		checkErr(err)
		t.Execute(w, nil)
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username != "aha" || password != "aha" {
			w.Write([]byte("Login fail!!!!!!!"))
		} else {
			w.Write([]byte("Login success, waiting a little time to check file"))
			getData()
		}
	}
}

func delFile() {
	for {
		time.Sleep(5 * 60 * time.Second)
		fileLoca := "./data/water.xlsx"
		err := os.Remove(fileLoca)
		if err != nil {
			fmt.Println("file remove Error!")
			fmt.Printf("%s\n", err)
		}
	}
}

func getData() {
	db, err := sql.Open("mysql", `user:password@tcp(ip:port)/dbname?charset=utf8`)

	checkErr(err)
	defer db.Close()

	fmt.Println("start ping")
	err = db.Ping()
	checkErr(err)
	fmt.Println("access succeed")

	var file *xlsx.File
	file = xlsx.NewFile()

	waterBill(db, file)
	getBalance(db, file)
	fruitWaterBill(db, file)
	spWithdrawal(db, file)

	err = file.Save("./data/water.xlsx")
	checkErr(err)
}

func waterBill(db *sql.DB, file *xlsx.File) {
	sentence := `
		select
		  tb_order.order_id as id, tb_order.ads_user_name as merchant, seller.user_name as seller, buyer.user_name as buyer, tb_order.ads_type as type,
		  tb_order.final_price as price, tb_order.fiat_amount as fiat, tb_order.crypto_received_amount as received,
		  tb_order.crypto_amount as paid, tb_order.create_time as timestamp
		from tb_order
		inner join tb_user seller on seller.id = tb_order.sell_id
		inner join tb_user buyer on buyer.id = tb_order.buy_id
		where
		  tb_order.order_status = 5
		  and tb_order.order_id not in ('1538421415409253', '1538422805852253', '1538423188108253', '1538424399466253', '1538428554132150', '1538428638041820')
		order by timestamp
		limit 100000;`
	rows, err := db.Query(sentence)
	checkErr(err)

	var sheet *xlsx.Sheet
	sheet, err = file.AddSheet("交易流水")
	checkErr(err)

	err = sheet.SetColWidth(0, 10, 18.0)
	checkErr(err)

	/* Add describe row */
	func () {
		var row *xlsx.Row
		row = sheet.AddRow()
		var cell *xlsx.Cell
		cell = row.AddCell();cell.Value = "id"
		cell = row.AddCell();cell.Value = "merchant"
		cell = row.AddCell();cell.Value = "seller"
		cell = row.AddCell();cell.Value = "buyer"
		cell = row.AddCell();cell.Value = "type"
		cell = row.AddCell();cell.Value = "price"
		cell = row.AddCell();cell.Value = "fiat"
		cell = row.AddCell();cell.Value = "received"
		cell = row.AddCell();cell.Value = "paid"
		cell = row.AddCell();cell.Value = "timestamp"
	}()

	for rows.Next() {
		var id uint64
		var merchant, seller, buyer, typee string
		var price, fiat, received, paid float64
		var timestamp string
		err = rows.Scan(&id, &merchant, &seller, &buyer, &typee, &price, &fiat, &received, &paid, &timestamp)
		checkErr(err)

		var row *xlsx.Row
		row = sheet.AddRow()
		var cell *xlsx.Cell

		cell = row.AddCell();cell.Value = strconv.FormatUint(id, 10)
		cell = row.AddCell();cell.Value = merchant
		cell = row.AddCell();cell.Value = seller
		cell = row.AddCell();cell.Value = buyer
		cell = row.AddCell();cell.Value = typee
		cell = row.AddCell();cell.Value = strconv.FormatFloat(price, 'f', 8, 64)
		cell = row.AddCell();cell.Value = strconv.FormatFloat(fiat, 'f', 8, 64)
		cell = row.AddCell();cell.Value = strconv.FormatFloat(received, 'f', 8, 64)
		cell = row.AddCell();cell.Value = strconv.FormatFloat(paid, 'f', 8, 64)
		cell = row.AddCell();cell.Value = timestamp
	}
}

func fruitWaterBill(db *sql.DB, file *xlsx.File) {
	sentence := `
		select
		  tb_order.order_id as id, (case when ads_type = 'buy' then buy_id else sell_id end) as merchant, tb_order.sell_id as seller, tb_order.buy_id as buyer, tb_order.ads_type as type,
		  tb_order.fiat_amount as fiat, tb_order.crypto_received_amount as received,
		  tb_order.crypto_amount as paid, tb_order.create_time as timestamp
		from tb_order
		where tb_order.order_status = 5 and (buy_id in (162, 170, 423, 812) or sell_id in (162, 170, 423, 812))
		order by timestamp
		limit 10000;`
	rows, err := db.Query(sentence)
	checkErr(err)

	var sheet *xlsx.Sheet
	sheet, err = file.AddSheet("水果店交易流水")
	checkErr(err)

	err = sheet.SetColWidth(0, 10, 18.0)
	checkErr(err)

	/* Add describe row */
	func () {
		var row *xlsx.Row
		row = sheet.AddRow()
		var cell *xlsx.Cell
		cell = row.AddCell();cell.Value = "id"
		cell = row.AddCell();cell.Value = "merchant"
		cell = row.AddCell();cell.Value = "seller"
		cell = row.AddCell();cell.Value = "buyer"
		cell = row.AddCell();cell.Value = "type"
		cell = row.AddCell();cell.Value = "fiat"
		cell = row.AddCell();cell.Value = "received"
		cell = row.AddCell();cell.Value = "paid"
		cell = row.AddCell();cell.Value = "timestamp"
	}()

	for rows.Next() {
		var id, merchant, seller, buyer uint64
		var typee string
		var fiat, received, paid float64
		var timestamp string
		err = rows.Scan(&id, &merchant, &seller, &buyer, &typee, &fiat, &received, &paid, &timestamp)
		checkErr(err)

		var row *xlsx.Row
		row = sheet.AddRow()
		var cell *xlsx.Cell

		cell = row.AddCell();cell.Value = strconv.FormatUint(id, 10)
		cell = row.AddCell();cell.Value = strconv.FormatUint(merchant, 10)
		cell = row.AddCell();cell.Value = strconv.FormatUint(seller, 10)
		cell = row.AddCell();cell.Value = strconv.FormatUint(buyer, 10)
		cell = row.AddCell();cell.Value = typee
		cell = row.AddCell();cell.Value = strconv.FormatFloat(fiat, 'f', 8, 64)
		cell = row.AddCell();cell.Value = strconv.FormatFloat(received, 'f', 8, 64)
		cell = row.AddCell();cell.Value = strconv.FormatFloat(paid, 'f', 8, 64)
		cell = row.AddCell();cell.Value = timestamp
	}
}

func spWithdrawal(db *sql.DB, file *xlsx.File) {
	sentence := `
		select id, crypto, amount, create_time, sender_id from tb_withdraw_info
		where
  			sender_id in (162, 170, 423, 812) and withdraw_status = 4
  			and (recipient_addr = '0x5bca22bc562ee8e7cbc81a678509aef1f230c258' or recipient_addr = '0x5BCA22BC562Ee8E7CBc81a678509aEF1F230C258')
		order by id;`
	rows, err := db.Query(sentence)
	checkErr(err)

	var sheet *xlsx.Sheet

	sheet, err = file.AddSheet("特殊扣款流水")
	checkErr(err)

	err = sheet.SetColWidth(0, 10, 18.0)
	checkErr(err)

	/* Add describe row */
	func () {
		var row *xlsx.Row
		row = sheet.AddRow()
		var cell *xlsx.Cell
		cell = row.AddCell();cell.Value = "id"
		cell = row.AddCell();cell.Value = "crypto"
		cell = row.AddCell();cell.Value = "amount"
		cell = row.AddCell();cell.Value = "createTime"
		cell = row.AddCell();cell.Value = "senderId"
	}()

	for rows.Next() {
		var id uint
		var crypto string
		var amount float64
		var createTime string
		var senderId int
		err = rows.Scan(&id, &crypto, &amount, &createTime, &senderId)
		checkErr(err)

		var row *xlsx.Row
		row = sheet.AddRow()
		var cell *xlsx.Cell

		cell = row.AddCell();cell.Value = strconv.FormatUint(uint64(id), 10)
		cell = row.AddCell();cell.Value = crypto
		cell = row.AddCell();cell.Value = strconv.FormatFloat(amount, 'f', 8, 64)
		cell = row.AddCell();cell.Value = createTime
		cell = row.AddCell();cell.Value = strconv.FormatInt(int64(senderId), 10)
	}
}

func getBalance(db *sql.DB, file *xlsx.File) {
	sentence := `
		select tb_user_wallet.user_id, tb_user.user_name as nickname, tb_user_wallet.crypto as crypto,
			tb_user_wallet.available, tb_user_wallet.frozen as frozen
		from tb_user_wallet
		left outer join tb_user on tb_user_wallet.user_id = tb_user.id
		where tb_user_wallet.user_id in (162, 170, 423, 812) and crypto = 'CFC';`

	rows, err := db.Query(sentence)
	checkErr(err)

	var sheet *xlsx.Sheet
	sheet, err = file.AddSheet("水果店余额")
	checkErr(err)

	err = sheet.SetColWidth(0, 10, 18.0)
	checkErr(err)

	/* Add describe row */
	func () {
		var row *xlsx.Row
		row = sheet.AddRow()
		var cell *xlsx.Cell
		cell = row.AddCell();cell.Value = "userId"
		cell = row.AddCell();cell.Value = "nickname"
		cell = row.AddCell();cell.Value = "crypto"
		cell = row.AddCell();cell.Value = "available"
		cell = row.AddCell();cell.Value = "frozen"
	}()

	for rows.Next() {
		var userId uint
		var nickname string
		var crypto string
		var available float64
		var frozen float64
		err = rows.Scan(&userId, &nickname, &crypto, &available, &frozen)
		checkErr(err)

		var row *xlsx.Row
		row = sheet.AddRow()
		var cell *xlsx.Cell

		cell = row.AddCell();cell.Value = strconv.FormatUint(uint64(userId), 10)
		cell = row.AddCell();cell.Value = nickname
		cell = row.AddCell();cell.Value = crypto
		cell = row.AddCell();cell.Value = strconv.FormatFloat(float64(available), 'f', 8, 64)
		cell = row.AddCell();cell.Value = strconv.FormatFloat(frozen, 'f', 8, 64)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}