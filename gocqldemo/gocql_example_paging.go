package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/gocql/gocql"
)

func iterate(iter *gocql.Iter) {
	fmt.Println("start iterate()")
	scanner := iter.Scanner()
	for scanner.Next() {
		var (
			id string
		)
		err := scanner.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id)
	}
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("end of iterate()")
}

func iterate2(iter *gocql.Iter) {
	fmt.Println("start iterate2()")
	columns := iter.Columns()

	for {
		rd, err := iter.RowData()
		if err != nil {
			panic(err)
		}

		if !iter.Scan(rd.Values...) {
			break
		}

		for index, value := range rd.Values {
			switch columns[index].TypeInfo.Type() {
			case gocql.TypeBigInt:
				fmt.Println(strconv.Itoa(int(*value.(*int64))))
			case gocql.TypeInt:
				fmt.Println(strconv.Itoa(int(*value.(*int))))
			case gocql.TypeFloat:
				fmt.Println(strconv.FormatFloat(*value.(*float64), 'f', 2, 64))
			case gocql.TypeVarchar:
				fmt.Println(*value.(*string))
			case gocql.TypeBlob:
				val := *value.(*[]byte)
				str := make([]byte, len(val))
				for index, num := range val {
					str[index] = byte(num)
				}
				fmt.Println(string(str))
			default:
				// We've encountered a type that we don't know yet.
				t := reflect.TypeOf(value)
				str := "?nil?"
				if t != nil {
					str = "?" + t.String() + "?"
				}
				fmt.Println(str)
			}
		}
	}
	fmt.Println("end of iterate2()")
}

func main() {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "routingjournal"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	loadRows := func(pageState []byte) []byte {
		// We use PageSize(2) for the sake of example, use larger values in production (default is 5000) for performance
		// reasons.
		iter := session.Query(`SELECT persistence_id FROM messages`).PageSize(10).PageState(pageState).Iter()
		pageState = iter.PageState()
		// iterate(iter)
		iterate2(iter)
		fmt.Printf("next page state: %+v\n", pageState)
		return pageState

	}

	var pageState []byte

	for {
		pageState = loadRows(pageState)
		if len(pageState) == 0 {
			break
		}
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
	}

	// 5 five
	// 1 one
	// next page state: [4 0 0 0 1 0 240 127 255 255 253 0]
	// 2 two
	// 4 four
	// next page state: [4 0 0 0 4 0 240 127 255 255 251 0]
	// 6 six
	// 3 three
	// next page state: [4 0 0 0 3 0 240 127 255 255 249 0]
	// next page state: []
}
