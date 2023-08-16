package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_agen_backend_api/configs"
	"github.com/nikitamirzani323/wl_agen_backend_api/db"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nleeper/goment"
)

const database_member_local = configs.DB_tbl_mst_master_agen_member

func Fetch_memberHome(idmasteragen string) (helpers.Response, error) {
	var obj entities.Model_member
	var arraobj []entities.Model_member
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idagenmember , username_agenmember, timezone_agenmember,  ipaddress_agenmember, 
			to_char(COALESCE(lastlogin_agenmember,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			name_agenmember , phone_agenmember, email_agenmember,  status_agenmember,
			create_agenmember, to_char(COALESCE(createdate_agenmember,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			update_agenmember, to_char(COALESCE(updatedate_agenmember,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + database_member_local + `  
			WHERE idmasteragen=$1 
			ORDER BY lastlogin_agenmember DESC   `

	row, err := con.QueryContext(ctx, sql_select, idmasteragen)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idagenmember_db, username_agenmember_db, timezone_agenmember_db, ipaddress_agenmember_db, lastlogin_agenmember_db string
			name_agenmember_db, phone_agenmember_db, email_agenmember_db, status_agenmember_db                                string
			create_agenmember_db, createdate_agenmember_db, update_agenmember_db, updatedate_agenmember_db                    string
		)

		err = row.Scan(&idagenmember_db, &username_agenmember_db, &timezone_agenmember_db, &ipaddress_agenmember_db, &lastlogin_agenmember_db,
			&name_agenmember_db, &phone_agenmember_db, &email_agenmember_db, &status_agenmember_db,
			&create_agenmember_db, &createdate_agenmember_db, &update_agenmember_db, &updatedate_agenmember_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if create_agenmember_db != "" {
			create = create_agenmember_db + ", " + createdate_agenmember_db
		}
		if update_agenmember_db != "" {
			update = update_agenmember_db + ", " + updatedate_agenmember_db
		}
		if status_agenmember_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}

		obj.Member_id = idagenmember_db
		obj.Member_username = username_agenmember_db
		obj.Member_timezone = timezone_agenmember_db
		obj.Member_ipaddress = ipaddress_agenmember_db
		obj.Member_lastlogin = lastlogin_agenmember_db
		obj.Member_name = name_agenmember_db
		obj.Member_phone = phone_agenmember_db
		obj.Member_email = email_agenmember_db
		obj.Member_status = status_agenmember_db
		obj.Member_status_css = status_css
		obj.Member_create = create
		obj.Member_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_member(admin, idmaster, idmasteragen, username, password, name, phone, email, status, sData, idrecord string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(database_member_local, "username_agenmember", username)
		if !flag {
			sql_insert := `
				insert into
				` + database_member_local + ` (
					idagenmember , idmaster, idmasteragen,
					username_agenmember, password_agenmember, lastlogin_agenmember, 
					name_agenmember, phone_agenmember,email_agenmember,status_agenmember,
					create_agenmember, createdate_agenmember   
				) values (
					$1, $2, $3,   
					$4, $5, $6,   
					$7, $8, $9, $10, 
					$11, $12
				)
			`
			field_column := database_member_local + tglnow.Format("YYYY")
			idrecord_counter := Get_counter(field_column)
			hashpass := helpers.HashPasswordMD5(password)
			create_date := tglnow.Format("YYYY-MM-DD HH:mm:ss")
			flag_insert, msg_insert := Exec_SQL(sql_insert, database_member_local, "INSERT",
				idmasteragen+tglnow.Format("YY")+tglnow.Format("MM")+tglnow.Format("DD")+tglnow.Format("HH")+strconv.Itoa(idrecord_counter), idmaster, idmasteragen,
				username, hashpass, create_date,
				name, phone, email, status,
				admin, create_date)

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Username"
		}
	} else {
		if password == "" {
			sql_update := `
				UPDATE 
				` + database_member_local + `  
				SET name_agenmember=$1, phone_agenmember=$2, email_agenmember=$3,
				status_agenmember=$4,    
				update_agenmember=$5, updatedate_agenmember=$6      
				WHERE idagenmember=$7    
			`

			flag_update, msg_update := Exec_SQL(sql_update, database_member_local, "UPDATE",
				name, phone, email, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

			if flag_update {
				msg = "Succes"
			} else {
				fmt.Println(msg_update)
			}
		} else {
			hashpass := helpers.HashPasswordMD5(password)
			sql_update := `
				UPDATE 
				` + database_member_local + `  
				SET password_agenmember=$1, name_agenmember=$2, phone_agenmember=$3, email_agenmember=$4, 
				status_agenmember=$5,     
				update_agenmember=$6, updatedate_agenmember=$7       
				WHERE idagenmember=$8     
			`

			flag_update, msg_update := Exec_SQL(sql_update, database_member_local, "UPDATE",
				hashpass, name, phone, email, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

			if flag_update {
				msg = "Succes"
			} else {
				fmt.Println(msg_update)
			}
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
