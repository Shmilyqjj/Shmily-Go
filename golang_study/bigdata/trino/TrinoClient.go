package main

import (
	"database/sql"
	"fmt"
)
import _ "github.com/trinodb/trino-go-client/trino"

type result struct {
	custkey   int
	name      string
	address   string
	nationkey int
}

/*
*
golang 连接 trino
*/
func main() {
	//args := os.Args
	//cert_path := args[1]
	//dsn := fmt.Sprintf("https://user:password@trino_server:7778?SSLCertPath=%s&catalog=hive&schema=default", cert_path)
	dsn := "http://user:password@localhost:8080?catalog=tpch&schema=sf1"
	fmt.Println("conn url: " + dsn)
	db, err := sql.Open("trino", dsn)

	var r result
	if err == nil {
		query, err := db.Query("select custkey, name, address, nationkey FROM tpch.sf1.customer limit 9")
		if err != nil {
			fmt.Println(err)
			return
		}
		for query.Next() {
			err := query.Scan(&r.custkey, &r.name, &r.address, &r.nationkey)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("custkey:%d name:%v address:%v nationkey:%d\n", r.custkey, r.name, r.address, r.nationkey)
		}
	}
}
