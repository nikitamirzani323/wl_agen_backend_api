package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_agen_backend_api/configs"
	"github.com/nikitamirzani323/wl_agen_backend_api/db"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nleeper/goment"
)

func Fetch_memberHome(idmasteragen string) (helpers.Responsemember, error) {
	var obj entities.Model_member
	var arraobj []entities.Model_member
	var objbanktype entities.Model_bankTypeshare
	var arraobjbanktype []entities.Model_bankTypeshare
	var res helpers.Responsemember
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	tbl_mst_member, tbl_mst_member_bank, _, _ := Get_mappingdatabase(idmasteragen)

	sql_select := `SELECT 
			idmember , username_member, timezone_member,  ipaddress_member, 
			to_char(COALESCE(lastlogin_member,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			name_member , phone_member, email_member,  status_member, 
			(cashin_member-cashout_member) as credit , 
			create_member, to_char(COALESCE(createdate_member,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			update_member, to_char(COALESCE(updatedate_member,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + tbl_mst_member + `  
			WHERE idmasteragen=$1 
			ORDER BY lastlogin_member DESC   `

	row, err := con.QueryContext(ctx, sql_select, idmasteragen)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idmember_db, username_member_db, timezone_member_db, ipaddress_member_db, lastlogin_member_db string
			name_member_db, phone_member_db, email_member_db, status_member_db                            string
			credit_db                                                                                     float64
			create_member_db, createdate_member_db, update_member_db, updatedate_member_db                string
		)

		err = row.Scan(&idmember_db, &username_member_db, &timezone_member_db, &ipaddress_member_db, &lastlogin_member_db,
			&name_member_db, &phone_member_db, &email_member_db, &status_member_db,
			&credit_db,
			&create_member_db, &createdate_member_db, &update_member_db, &updatedate_member_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if create_member_db != "" {
			create = create_member_db + ", " + createdate_member_db
		}
		if update_member_db != "" {
			update = update_member_db + ", " + updatedate_member_db
		}
		if status_member_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}

		//BANK
		var objbank entities.Model_memberbank
		var arraobjbank []entities.Model_memberbank
		sql_selectbank := `SELECT 
			idmemberbank,idbanktype, norekbank_memberbank, nmownerbank_memberbank 
			FROM ` + tbl_mst_member_bank + ` 
			WHERE idmember = $1   
		`
		row_bank, err_bank := con.QueryContext(ctx, sql_selectbank, idmember_db)
		helpers.ErrorCheck(err_bank)
		for row_bank.Next() {
			var (
				idmemberbank_db                                                   int
				idbanktype_db, norekbank_memberbank_db, nmownerbank_memberbank_db string
			)
			err_bank = row_bank.Scan(&idmemberbank_db, &idbanktype_db, &norekbank_memberbank_db, &nmownerbank_memberbank_db)

			objbank.Memberbank_id = idmemberbank_db
			objbank.Memberbank_idbanktype = idbanktype_db
			objbank.Memberbank_nmownerbank = nmownerbank_memberbank_db
			objbank.Memberbank_norek = norekbank_memberbank_db
			arraobjbank = append(arraobjbank, objbank)
		}
		defer row_bank.Close()

		idcurr := _GetDefaultCurr(idmasteragen)
		multiplier := _GetMultiplier(idcurr)
		var credit float64 = credit_db * float64(multiplier)

		obj.Member_id = idmember_db
		obj.Member_username = username_member_db
		obj.Member_timezone = timezone_member_db
		obj.Member_ipaddress = ipaddress_member_db
		obj.Member_lastlogin = lastlogin_member_db
		obj.Member_name = name_member_db
		obj.Member_phone = phone_member_db
		obj.Member_email = email_member_db
		obj.Member_credit = credit
		obj.Member_listbank = arraobjbank
		obj.Member_status = status_member_db
		obj.Member_status_css = status_css
		obj.Member_create = create
		obj.Member_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	sql_selectbanktype := `SELECT 
			B.nmcatebank, A.idbanktype  
			FROM ` + configs.DB_tbl_mst_banktype + ` as A 
			JOIN ` + configs.DB_tbl_mst_cate_bank + ` as B ON B.idcatebank = A.idcatebank 
			ORDER BY B.nmcatebank,A.idbanktype ASC    
	`
	rowbanktype, errbanktype := con.QueryContext(ctx, sql_selectbanktype)
	helpers.ErrorCheck(errbanktype)
	for rowbanktype.Next() {
		var (
			nmcatebank_db, idbanktype_db string
		)

		errbanktype = rowbanktype.Scan(&nmcatebank_db, &idbanktype_db)

		helpers.ErrorCheck(errbanktype)

		objbanktype.Catebank_name = nmcatebank_db
		objbanktype.Banktype_id = idbanktype_db
		arraobjbanktype = append(arraobjbanktype, objbanktype)
		msg = "Success"
	}
	defer rowbanktype.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listbank = arraobjbanktype
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_memberSearch(idmasteragen, search string) (helpers.Response, error) {
	var obj entities.Model_membershare
	var arraobj []entities.Model_membershare
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	tbl_mst_member, tbl_mst_member_bank, _, _ := Get_mappingdatabase(idmasteragen)
	perpage := 50

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "idmember , username_member, name_member, "
	sql_select += "(cashin_member-cashout_member) as credit  "
	sql_select += "FROM " + tbl_mst_member + "  "
	if search == "" {
		sql_select += "WHERE idmasteragen = '" + idmasteragen + "' "
		sql_select += "ORDER BY name_member DESC   LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE idmasteragen = '" + idmasteragen + "' "
		sql_select += "AND LOWER(username_member) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY name_member DESC LIMIT " + strconv.Itoa(perpage)
	}

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idmember_db, username_member_db, name_member_db string
			credit_db                                       float64
		)

		err = row.Scan(&idmember_db, &username_member_db, &name_member_db, &credit_db)
		helpers.ErrorCheck(err)

		//BANK
		var objbank entities.Model_memberbankshare
		var arraobjbank []entities.Model_memberbankshare
		sql_selectbank := `SELECT 
			idmemberbank,idbanktype, norekbank_memberbank, nmownerbank_memberbank  
			FROM ` + tbl_mst_member_bank + ` 
			WHERE idmember = $1   
		`
		row_bank, err_bank := con.QueryContext(ctx, sql_selectbank, idmember_db)
		helpers.ErrorCheck(err_bank)
		for row_bank.Next() {
			var (
				idmemberbank_db                                                   int
				idbanktype_db, norekbank_memberbank_db, nmownerbank_memberbank_db string
			)
			err_bank = row_bank.Scan(&idmemberbank_db, &idbanktype_db, &norekbank_memberbank_db, &nmownerbank_memberbank_db)

			objbank.Memberbank_id = idmemberbank_db
			objbank.Memberbank_info = idbanktype_db + "-" + norekbank_memberbank_db + "-" + nmownerbank_memberbank_db
			arraobjbank = append(arraobjbank, objbank)
		}
		defer row_bank.Close()
		idcurr := _GetDefaultCurr(idmasteragen)
		multiplier := _GetMultiplier(idcurr)
		var credit float64 = credit_db * float64(multiplier)

		obj.Member_id = idmember_db
		obj.Member_name = username_member_db + "-" + name_member_db
		obj.Member_credit = credit
		obj.Member_listbank = arraobjbank
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

	tbl_mst_member, _, _, _ := Get_mappingdatabase(idmasteragen)
	if sData == "New" {
		flag = CheckDB(tbl_mst_member, "username_member", username)
		if !flag {
			sql_insert := `
				insert into
				` + tbl_mst_member + ` (
					idmember , idmaster, idmasteragen,
					username_member, password_member, lastlogin_member, 
					name_member, phone_member,email_member,status_member,
					create_member, createdate_member   
				) values (
					$1, $2, $3,   
					$4, $5, $6,   
					$7, $8, $9, $10, 
					$11, $12
				)
			`
			field_column := tbl_mst_member + tglnow.Format("YYYY")
			idrecord_counter := Get_counter(field_column)
			hashpass := helpers.HashPasswordMD5(password)
			create_date := tglnow.Format("YYYY-MM-DD HH:mm:ss")
			flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_mst_member, "INSERT",
				idmasteragen+"MBR"+tglnow.Format("YY")+tglnow.Format("MM")+tglnow.Format("DD")+tglnow.Format("HH")+strconv.Itoa(idrecord_counter), idmaster, idmasteragen,
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
		// idmember , idmaster, idmasteragen,
		// username_member, password_member, lastlogin_member,
		// name_member, phone_member,email_member,status_member,
		if password == "" {
			sql_update := `
				UPDATE 
				` + tbl_mst_member + `  
				SET name_member=$1, phone_member=$2, email_member=$3,
				status_member=$4,    
				update_member=$5, updatedate_member=$6      
				WHERE idmember=$7    
			`

			flag_update, msg_update := Exec_SQL(sql_update, tbl_mst_member, "UPDATE",
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
				` + tbl_mst_member + `  
				SET password_member=$1, name_member=$2, phone_member=$3, email_member=$4, 
				status_member=$5,     
				update_member=$6, updatedate_member=$7       
				WHERE idmember=$8     
			`

			flag_update, msg_update := Exec_SQL(sql_update, tbl_mst_member, "UPDATE",
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
func Save_memberbank(admin, idmasteragen, idmember, idbanktype, norek, name, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	_, tbl_mst_member_bank, _, _ := Get_mappingdatabase(idmasteragen)
	if sData == "New" {
		sql_insert := `
			insert into
			` + tbl_mst_member_bank + ` (
				idmemberbank , idmember, 
				idbanktype, norekbank_memberbank, nmownerbank_memberbank, 
				create_memberbank, createdate_memberbank     
			) values (
				$1, $2,    
				$3, $4, $5,    
				$6, $7
			)
		`
		field_column := tbl_mst_member_bank + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_mst_member_bank, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), idmember, idbanktype, norek, name,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Delete_memberbank(idmember, idmasteragen string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	render_page := time.Now()

	_, tbl_mst_member_bank, _, _ := Get_mappingdatabase(idmasteragen)
	sql_delete := `
				DELETE FROM
				` + tbl_mst_member_bank + ` 
				WHERE idmemberbank=$1 AND idmember=$2  
			`
	flag_delete, msg_delete := Exec_SQL(sql_delete, tbl_mst_member_bank, "DELETE", idrecord, idmember)

	if !flag_delete {
		fmt.Println(msg_delete)
	} else {
		msg = "Succes"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
