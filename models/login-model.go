package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikitamirzani323/wl_agen_backend_api/configs"
	"github.com/nikitamirzani323/wl_agen_backend_api/db"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nleeper/goment"
)

func Login_Model(username, password, ipaddress string) (bool, string, string, string, string, error) {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	tglnow, _ := goment.New()
	var idagenadmin, idmasteragen, passwordDB, ruleDB, tipeDB string
	sql_select := `
			SELECT
			idagenadmin, idmasteragen, passwordagen_admin, idagenadminrule, tipeagen_admin    
			FROM ` + configs.DB_tbl_mst_master_agen_admin + ` 
			WHERE usernameagen_admin  = $1
			AND statusagenadmin = 'Y' 
		`

	row := con.QueryRowContext(ctx, sql_select, username)
	switch e := row.Scan(&idagenadmin, &idmasteragen, &passwordDB, &ruleDB, &tipeDB); e {
	case sql.ErrNoRows:
		return false, "", "", "", "", errors.New("Username and Password Not Found")
	case nil:
		flag = true
	default:
		return false, "", "", "", "", errors.New("Username and Password Not Found")
	}

	hashpass := helpers.HashPasswordMD5(password)

	if hashpass != passwordDB {
		return false, "", "", "", "", nil
	}

	if flag {

		sql_update := `
			UPDATE ` + configs.DB_tbl_mst_master_agen_admin + ` 
			SET lastloginagen_admin=$1, ipaddress_admin=$2,  
			updateagenadmin=$3,  updatedateagenadmin=$4   
			WHERE idagenadmin  = $5 
			AND usernameagen_admin  = $6  
			AND statusagenadmin = 'Y' 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_master_agen_admin, "UPDATE",
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), ipaddress, idagenadmin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), idagenadmin, username)

		if flag_update {
			flag = true
		} else {
			fmt.Println(msg_update)
		}
	}

	return true, idmasteragen, idagenadmin, ruleDB, tipeDB, nil
}
