package models

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_agen_backend_api/configs"
	"github.com/nikitamirzani323/wl_agen_backend_api/db"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nleeper/goment"
)

const database_master_local = configs.DB_tbl_mst_master
const database_masteradmin_local = configs.DB_tbl_mst_master_admin
const database_masteragen_local = configs.DB_tbl_mst_master_agen

func Fetch_masterHome() (helpers.Responsemaster, error) {
	var obj entities.Model_master
	var arraobj []entities.Model_master
	var objcurr entities.Model_currshare
	var arraobjcurr []entities.Model_currshare
	var objbank entities.Model_bankTypeshare
	var arraobjbank []entities.Model_bankTypeshare
	var res helpers.Responsemaster
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idmaster ,  
			to_char(COALESCE(startjoinmaster,now()), 'YYYY-MM-DD HH24:MI:SS'),
			to_char(COALESCE(endjoinmaster,now()), 'YYYY-MM-DD HH24:MI:SS'),
			idcurr , nmmaster, nmowner, phone1owner, phone2owner, emailowner, notemaster, statusmaster, 
			idbanktype , norekbank, nmownerbank, 
			createmaster, to_char(COALESCE(createdatemaster,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatemaster, to_char(COALESCE(updatedatemaster,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + database_master_local + `  
			ORDER BY createdatemaster DESC   `

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idmaster_db, idcurr_db, nmmaster_db, nmowner_db, phone1owner_db, phone2owner_db, emailowner_db, notemaster_db, statusmaster_db string
			idbanktype_db, norekbank_db, nmownerbank_db                                                                                    string
			startjoinmaster_db, endjoinmaster_db                                                                                           string
			createmaster_db, createdatemaster_db, updatemaster_db, updatedatemaster_db                                                     string
		)

		err = row.Scan(&idmaster_db, &startjoinmaster_db, &endjoinmaster_db,
			&idcurr_db, &nmmaster_db, &nmowner_db, &phone1owner_db, &phone2owner_db, &emailowner_db, &notemaster_db, &statusmaster_db,
			&idbanktype_db, &norekbank_db, &nmownerbank_db,
			&createmaster_db, &createdatemaster_db, &updatemaster_db, &updatedatemaster_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if createmaster_db != "" {
			create = createmaster_db + ", " + createdatemaster_db
		}
		if updatemaster_db != "" {
			update = updatemaster_db + ", " + updatedatemaster_db
		}
		if statusmaster_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}

		//MASTER ADMIN
		var objmasteradmin entities.Model_masteradmin
		var arraobjmasteradmin []entities.Model_masteradmin
		sql_selectmasteradmin := `SELECT 
			idmasteradmin,tipe_masteradmin, username_masteradmin, name_masteradmin, phone1_masteradmin,phone2_masteradmin, statusmasteradmin,  
			createmasteradmin, to_char(COALESCE(createdatemasteradmin,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatemasteradmin, to_char(COALESCE(updatedatemasteradmin,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + configs.DB_tbl_mst_master_admin + ` 
			WHERE idmaster = $1   
		`
		row_banktype, err := con.QueryContext(ctx, sql_selectmasteradmin, idmaster_db)
		helpers.ErrorCheck(err)
		for row_banktype.Next() {
			var (
				idmasteradmin_db                                                                                                                      int
				tipe_masteradmin_db, username_masteradmin_db, name_masteradmin_db, phone1_masteradmin_db, phone2_masteradmin_db, statusmasteradmin_db string
				createmasteradmin_db, createdatemasteradmin_db, updatemasteradmin_db, updatedatemasteradmin_db                                        string
			)
			err = row_banktype.Scan(&idmasteradmin_db, &tipe_masteradmin_db, &username_masteradmin_db, &name_masteradmin_db, &phone1_masteradmin_db, &phone2_masteradmin_db, &statusmasteradmin_db,
				&createmasteradmin_db, &createdatemasteradmin_db, &updatemasteradmin_db, &updatedatemasteradmin_db)

			create_admin := ""
			update_admin := ""
			status_css_admin := configs.STATUS_CANCEL
			if createmasteradmin_db != "" {
				create_admin = createmasteradmin_db + ", " + createdatemasteradmin_db
			}
			if updatemasteradmin_db != "" {
				update_admin = updatemasteradmin_db + ", " + updatedatemasteradmin_db
			}
			if statusmasteradmin_db == "Y" {
				status_css_admin = configs.STATUS_COMPLETE
			}

			objmasteradmin.Masteradmin_id = idmasteradmin_db
			objmasteradmin.Masteradmin_tipe = tipe_masteradmin_db
			objmasteradmin.Masteradmin_username = username_masteradmin_db
			objmasteradmin.Masteradmin_name = name_masteradmin_db
			objmasteradmin.Masteradmin_phone1 = phone1_masteradmin_db
			objmasteradmin.Masteradmin_phone2 = phone2_masteradmin_db
			objmasteradmin.Masteradmin_status = statusmasteradmin_db
			objmasteradmin.Masteradmin_status_css = status_css_admin
			objmasteradmin.Masteradmin_create = create_admin
			objmasteradmin.Masteradmin_update = update_admin
			arraobjmasteradmin = append(arraobjmasteradmin, objmasteradmin)
		}
		defer row_banktype.Close()

		//MASTER AGEN
		var objmasteragen entities.Model_masteragen
		var arraobjmasteraagen []entities.Model_masteragen
		sql_selectmasteraagen := `SELECT 
			idmasteragen,idcurr, nmagen, nmowneragen, phone1agen,phone2agen, emailagen,noteagen,  
			idbanktype,norekbank, nmownerbank,statusmasteragen,  
			createmasteragen, to_char(COALESCE(createdatemasteragen,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatemasteragen, to_char(COALESCE(updatedatemasteragen,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + database_masteragen_local + ` 
			WHERE idmaster = $1   
		`
		row_masteragen, err := con.QueryContext(ctx, sql_selectmasteraagen, idmaster_db)
		helpers.ErrorCheck(err)
		for row_masteragen.Next() {
			var (
				idmasteragen_db, idcurr_db, nmagen_db, nmowneragen_db, phone1agen_db, phone2agen_db, emailagen_db, noteagen_db string
				idbanktype_db, norekbank_db, nmownerbank_db, statusmasteragen_db                                               string
				createmasteragen_db, createdatemasteragen_db, updatemasteragen_db, updatedatemasteragen_db                     string
			)
			err = row_masteragen.Scan(&idmasteragen_db, &idcurr_db, &nmagen_db, &nmowneragen_db, &phone1agen_db, &phone2agen_db, &emailagen_db, &noteagen_db,
				&idbanktype_db, &norekbank_db, &nmownerbank_db, &statusmasteragen_db,
				&createmasteragen_db, &createdatemasteragen_db, &updatemasteragen_db, &updatedatemasteragen_db)

			create_agen := ""
			update_agen := ""
			status_css_agen := configs.STATUS_CANCEL
			if createmasteragen_db != "" {
				create_agen = createmasteragen_db + ", " + createdatemasteragen_db
			}
			if updatemasteragen_db != "" {
				update_agen = updatemasteragen_db + ", " + updatedatemasteragen_db
			}
			if statusmasteragen_db == "Y" {
				status_css_agen = configs.STATUS_COMPLETE
			}

			objmasteragen.Masteragen_id = idmasteragen_db
			objmasteragen.Masteragen_idcurr = idcurr_db
			objmasteragen.Masteragen_nmagen = nmagen_db
			objmasteragen.Masteragen_owner = nmowneragen_db
			objmasteragen.Masteragen_phone1 = phone1agen_db
			objmasteragen.Masteragen_phone2 = phone2agen_db
			objmasteragen.Masteragen_email = emailagen_db
			objmasteragen.Masteragen_note = noteagen_db
			objmasteragen.Masteragen_bank_id = idbanktype_db
			objmasteragen.Masteragen_bank_name = nmownerbank_db
			objmasteragen.Masteragen_bank_norek = norekbank_db
			objmasteragen.Masteragen_status = statusmasteragen_db
			objmasteragen.Masteragen_status_css = status_css_agen
			objmasteragen.Masteragen_create = create_agen
			objmasteragen.Masteragen_update = update_agen
			arraobjmasteraagen = append(arraobjmasteraagen, objmasteragen)
		}
		defer row_masteragen.Close()

		obj.Master_id = idmaster_db
		obj.Master_start = startjoinmaster_db
		obj.Master_end = endjoinmaster_db
		obj.Master_idcurr = idcurr_db
		obj.Master_name = nmmaster_db
		obj.Master_owner = nmowner_db
		obj.Master_phone1 = phone1owner_db
		obj.Master_phone2 = phone2owner_db
		obj.Master_email = emailowner_db
		obj.Master_note = notemaster_db
		obj.Master_bank_id = idbanktype_db
		obj.Master_bank_name = nmownerbank_db
		obj.Master_bank_norek = norekbank_db
		obj.Master_status = statusmaster_db
		obj.Master_status_css = status_css
		obj.Master_listadmin = arraobjmasteradmin
		obj.Master_listagen = arraobjmasteraagen
		obj.Master_create = create
		obj.Master_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	sql_selectcurr := `SELECT 
			idcurr  
			FROM ` + configs.DB_tbl_mst_curr + ` 
			ORDER BY idcurr ASC    
	`
	rowcurr, errcurr := con.QueryContext(ctx, sql_selectcurr)
	helpers.ErrorCheck(errcurr)
	for rowcurr.Next() {
		var (
			idcurr_db string
		)

		errcurr = rowcurr.Scan(&idcurr_db)

		helpers.ErrorCheck(errcurr)

		objcurr.Curr_id = idcurr_db
		arraobjcurr = append(arraobjcurr, objcurr)
		msg = "Success"
	}
	defer rowcurr.Close()

	sql_selectbank := `SELECT 
			B.nmcatebank, A.idbanktype  
			FROM ` + configs.DB_tbl_mst_banktype + ` as A 
			JOIN ` + configs.DB_tbl_mst_cate_bank + ` as B ON B.idcatebank = A.idcatebank 
			ORDER BY B.nmcatebank,A.idbanktype ASC    
	`
	rowbank, errbank := con.QueryContext(ctx, sql_selectbank)
	helpers.ErrorCheck(errbank)
	for rowbank.Next() {
		var (
			nmcatebank_db, idbanktype_db string
		)

		errbank = rowbank.Scan(&nmcatebank_db, &idbanktype_db)

		helpers.ErrorCheck(errbank)

		objbank.Catebank_name = nmcatebank_db
		objbank.Banktype_id = idbanktype_db
		arraobjbank = append(arraobjbank, objbank)
		msg = "Success"
	}
	defer rowbank.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listcurr = arraobjcurr
	res.Listbank = arraobjbank
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_masteragenAdmin(idmasteragen string) (helpers.Response, error) {
	var obj entities.Model_masteragenadmin
	var arraobj []entities.Model_masteragenadmin
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	tbl_mst_admin := configs.DB_tbl_mst_master_agen_admin
	sql_select := `SELECT 
			idagenadmin, tipeagen_admin, usernameagen_admin,     
			nameagen_admin , phone1agen_admin, phone2agen_admin, statusagenadmin, 
			to_char(COALESCE(lastloginagen_admin,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			createagenadmin, to_char(COALESCE(createdateagenadmin,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updateagenadmin, to_char(COALESCE(updatedateagenadmin,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + tbl_mst_admin + `  
			WHERE idmasteragen= '` + idmasteragen + `'   
			ORDER BY lastloginagen_admin DESC   `
	log.Println(sql_select)
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idadmin_db, tipe_admin_db, username_admin_db                                        string
			name_admin_db, phone1_admin_db, phone2_admin_db, statusadmin_db, lastlogin_admin_db string
			createadmin_db, createdateadmin_db, updateadmin_db, updatedateadmin_db              string
		)

		err = row.Scan(&idadmin_db, &tipe_admin_db, &username_admin_db,
			&name_admin_db, &phone1_admin_db, &phone2_admin_db, &statusadmin_db, &lastlogin_admin_db,
			&createadmin_db, &createdateadmin_db, &updateadmin_db, &updatedateadmin_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		lastlogin := ""
		if createadmin_db != "" {
			create = createadmin_db + ", " + createdateadmin_db
		}
		if updateadmin_db != "" {
			update = updateadmin_db + ", " + updatedateadmin_db
		}
		if statusadmin_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}
		if lastlogin_admin_db != createdateadmin_db {
			lastlogin = lastlogin_admin_db
		}

		obj.Masteragenadmin_id = idadmin_db
		obj.Masteragenadmin_tipe = tipe_admin_db
		obj.Masteragenadmin_username = username_admin_db
		obj.Masteragenadmin_lastlogin = lastlogin
		obj.Masteragenadmin_name = name_admin_db
		obj.Masteragenadmin_phone1 = phone1_admin_db
		obj.Masteragenadmin_phone2 = phone2_admin_db
		obj.Masteragenadmin_status = statusadmin_db
		obj.Masteragenadmin_status_css = status_css
		obj.Masteragenadmin_create = create
		obj.Masteragenadmin_update = update
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
func Save_master(admin, idrecord, idcurr, name, owner, phone1, phone2, email, note, status, idbanktype, norekbank, nmownerbank, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(database_master_local, "idmaster", idrecord)
		if !flag {
			sql_insert := `
				insert into
				` + database_master_local + ` (
					idmaster , startjoinmaster, endjoinmaster, 
					idcurr , nmmaster, nmowner , phone1owner, phone2owner, emailowner, notemaster, statusmaster, 
					idbanktype , norekbank, nmownerbank, 
					createmaster, createdatemaster   
				) values (
					$1, $2, $3,    
					$4, $5, $6, $7, $8, $9, $10, $11,    
					$12, $13, $14,  
					$15, $16  
				)
			`
			start := tglnow.Format("YYYY-MM-DD HH:mm:ss")
			flag_insert, msg_insert := Exec_SQL(sql_insert, database_master_local, "INSERT",
				idrecord, start, start, idcurr, name, owner, phone1, phone2, email, note, status,
				idbanktype, norekbank, nmownerbank,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update := `
				UPDATE 
				` + database_master_local + `  
				SET idcurr=$1, nmmaster=$2, nmowner=$3,  phone1owner=$4,  phone2owner=$5, emailowner=$6, notemaster=$7, statusmaster=$8,  
				idbanktype=$9, norekbank=$10, nmownerbank=$11,      
				updatemaster=$12, updatedatemaster=$13      
				WHERE idmaster=$14     
			`

		flag_update, msg_update := Exec_SQL(sql_update, database_master_local, "UPDATE",
			idcurr, name, owner, phone1, phone2, email, note, status,
			idbanktype, norekbank, nmownerbank,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
		} else {
			fmt.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Save_masteradmin(admin, idmaster, tipe, username, password, name, phone1, phone2, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(database_masteradmin_local, "username_masteradmin", username)
		if !flag {
			sql_insert := `
				insert into
				` + database_masteradmin_local + ` (
					idmasteradmin, idmaster , tipe_masteradmin, username_masteradmin, password_masteradmin, 
					name_masteradmin , phone1_masteradmin, phone2_masteradmin, statusmasteradmin, 
					createmasteradmin, createdatemasteradmin   
				) values (
					$1, $2, $3, $4, $5,   
					$6, $7, $8, $9, 
					$10, $11
				)
			`
			field_column := database_masteradmin_local + tglnow.Format("YYYY")
			idrecord_counter := Get_counter(field_column)
			hashpass := helpers.HashPasswordMD5(password)
			flag_insert, msg_insert := Exec_SQL(sql_insert, database_masteradmin_local, "INSERT",
				tglnow.Format("YY")+strconv.Itoa(idrecord_counter), idmaster, tipe, username, hashpass, name, phone1, phone2, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		if password == "" {
			sql_update := `
				UPDATE 
				` + database_masteradmin_local + `  
				SET tipe_masteradmin=$1, name_masteradmin=$2, phone1_masteradmin=$3, phone2_masteradmin=$4, statusmasteradmin=$5,  
				updatemasteradmin=$6, updatedatemasteradmin=$7        
				WHERE idmaster=$8 AND idmasteradmin=$9        
			`

			flag_update, msg_update := Exec_SQL(sql_update, database_masteradmin_local, "UPDATE",
				tipe, name, phone1, phone2, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmaster, idrecord)

			if flag_update {
				msg = "Succes"
			} else {
				fmt.Println(msg_update)
			}
		} else {
			hashpass := helpers.HashPasswordMD5(password)
			sql_update := `
				UPDATE 
				` + database_masteradmin_local + `  
				SET tipe_masteradmin=$1, password_masteradmin=$2, name_masteradmin=$3, phone1_masteradmin=$4, phone2_masteradmin=$5, statusmasteradmin=$6,  
				updatemasteradmin=$7, updatedatemasteradmin=$8         
				WHERE idmaster=$9 AND idmasteradmin=$10        
			`

			flag_update, msg_update := Exec_SQL(sql_update, database_masteradmin_local, "UPDATE",
				tipe, hashpass, name, phone1, phone2, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmaster, idrecord)

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
func Save_masteragen(admin, idrecord, idmaster, idcurr, name, owner, phone1, phone2, email, note, status, idbanktype, norekbank, nmownerbank, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + database_masteragen_local + ` (
					idmasteragen, idmaster , startjoinagen, endjoinagen, 
					idcurr, nmagen, nmowneragen , phone1agen, phone2agen, emailagen, noteagen, statusmasteragen, 
					idbanktype , norekbank, nmownerbank, 
					createmasteragen, createdatemasteragen    
				) values (
					$1, $2, $3, $4,   
					$5, $6, $7, $8, $9, $10, $11, $12,     
					$13, $14, $15,  
					$16, $17   
				)
			`
		field_column := database_masteragen_local + tglnow.Format("YY") + tglnow.Format("MM")
		idrecord_counter := Get_counter(field_column)
		start := tglnow.Format("YYYY-MM-DD HH:mm:ss")
		flag_insert, msg_insert := Exec_SQL(sql_insert, database_masteragen_local, "INSERT",
			idmaster+tglnow.Format("YY")+tglnow.Format("MM")+strconv.Itoa(idrecord_counter), idmaster, start, start,
			idcurr, name, owner, phone1, phone2, email, note, status,
			idbanktype, norekbank, nmownerbank,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + database_masteragen_local + `  
				SET nmagen=$1, nmowneragen=$2, phone1agen=$3,  phone2agen=$4,  emailagen=$5, noteagen=$6, statusmasteragen=$7, 
				idbanktype=$8, norekbank=$9, nmownerbank=$10,      
				updatemasteragen=$11, updatedatemasteragen=$12       
				WHERE idmasteragen=$13      
			`

		flag_update, msg_update := Exec_SQL(sql_update, database_masteragen_local, "UPDATE",
			name, owner, phone1, phone2, email, note, status,
			idbanktype, norekbank, nmownerbank,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			msg = "Succes"
		} else {
			fmt.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Save_masteragenadmin(admin, idrecord, idmasteragen, tipe, username, password, name, phone1, phone2, status, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false
	tbl_mst_admin := configs.DB_tbl_mst_master_agen_admin
	if sData == "New" {
		flag = CheckDB(tbl_mst_admin, "usernameagen_admin", username)
		if !flag {
			sql_insert := `
					insert into
					` + tbl_mst_admin + ` (
						idagenadmin, idmasteragen , tipeagen_admin, usernameagen_admin, passwordagen_admin, lastloginagen_admin,   
						nameagen_admin, phone1agen_admin, phone2agen_admin , statusagenadmin, 
						createagenadmin, createdateagenadmin    
					) values (
						$1, $2, $3, $4, $5, $6,   
						$7, $8, $9, $10,     
						$11, $12  
					)
				`
			field_column := idmasteragen + tbl_mst_admin + tglnow.Format("YY")
			idrecord_counter := Get_counter(field_column)
			hashpass := helpers.HashPasswordMD5(password)
			create_date := tglnow.Format("YYYY-MM-DD HH:mm:ss")
			flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_mst_admin, "INSERT",
				idmasteragen+"-"+tglnow.Format("YY")+strconv.Itoa(idrecord_counter), idmasteragen, tipe, username, hashpass, create_date,
				name, phone1, phone2, status,
				admin, create_date)

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		if password == "" {
			sql_update := `
				UPDATE 
				` + tbl_mst_admin + `  
				SET tipeagen_admin=$1, nameagen_admin=$2, phone1agen_admin=$3, phone2agen_admin=$4, statusagenadmin=$5,  
				updateagenadmin=$6, updatedateagenadmin=$7        
				WHERE idmasteragen=$8 AND idagenadmin=$9        
			`

			flag_update, msg_update := Exec_SQL(sql_update, tbl_mst_admin, "UPDATE",
				tipe, name, phone1, phone2, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmasteragen, idrecord)

			if flag_update {
				msg = "Succes"
			} else {
				fmt.Println(msg_update)
			}
		} else {
			hashpass := helpers.HashPasswordMD5(password)
			sql_update := `
				UPDATE 
				` + tbl_mst_admin + `  
				SET tipeagen_admin=$1, passwordagen_admin=$2, nameagen_admin=$3, phone1agen_admin=$4, phone2agen_admin=$5, statusagenadmin=$6,  
				updateagenadmin=$7, updatedateagenadmin=$8         
				WHERE idmasteragen=$8 AND idagenadmin=$9        
			`

			flag_update, msg_update := Exec_SQL(sql_update, tbl_mst_admin, "UPDATE",
				tipe, hashpass, name, phone1, phone2, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmasteragen, idrecord)

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
