package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/akamajoris/ql"

	"github.com/akamajoris/ngorm"
	"github.com/akamajoris/ngorm/errmsg"

	"github.com/vikingo-project/vsat/shared"
)

var (
	connection *ngorm.DB //*gorm.DB
)

func Init() {
	var err error
	// connection, err = gorm.Open(mysql.Open(shared.Config.DB), &gorm.Config{}) // sql.Open("sqlite3_extended", shared.Config.DB)
	// connection, err = sql.Open("sqlite3_extended", shared.Config.DB)
	connection, err = ngorm.Open("ql", shared.Config.DB)
	if err != nil {
		log.Fatal(err)
	}
	migrate()

}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func SQLQuery(query string, args ...interface{}) (result []map[string]string, err error) {
	result = make([]map[string]string, 0)

	rows, err := connection.Begin().SQLCommon().Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	columnNames, err := rows.Columns()
	if err != nil {
		return
	}

	vals := make([]interface{}, len(columnNames))
	for rows.Next() {
		for i := range columnNames {
			vals[i] = &vals[i]
		}
		err = rows.Scan(vals...)
		if err != nil {
			return
		}
		var row = make(map[string]string)
		for i := range columnNames {
			switch vals[i].(type) {
			default:
				row[columnNames[i]] = fmt.Sprintf("%s", vals[i])
			case int, int64:
				row[columnNames[i]] = fmt.Sprintf("%d", vals[i])
			case float32, float64:
				row[columnNames[i]] = fmt.Sprintf("%f", vals[i])
			case nil:
				row[columnNames[i]] = ""
			}
		}
		result = append(result, row)
	}
	return
}

func ErrRecordNotFound(err error) bool {
	return errors.Is(err, errmsg.ErrRecordNotFound)
}

func GetConnection() *ngorm.DB {
	return connection
}

func Close() {
	connection.Close()
}
