package reposistory

import "github.com/go-sql-driver/mysql"

func IsDuplicateErr(err error) bool {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number == 1062 {
			return true
		}
	}
	return false
}
