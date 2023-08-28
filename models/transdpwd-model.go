package models

import (
	"context"
	"database/sql"
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

func Fetch_transdpwdHome(idmasteragen string) (helpers.Response, error) {
	var obj entities.Model_transdpwd
	var arraobj []entities.Model_transdpwd
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	_, _, tbl_trx_dpwd, _ := Get_mappingdatabase(idmasteragen)
	sql_select := `SELECT 
			iddpwd , to_char(COALESCE(date_dpwd,now()), 'YYYY-MM-DD'), idcurr,  
			tipedocuser_dpwd, tipedoc_dpwd , tipeakun_dpwd, idagenmember,  ipaddress_dpwd, timezone_dpwd,  
			bank_in, bank_in_info , bank_out, bank_out_info, 
			round(amount_dpwd*multiplier_dpwd) as amount_dpwd , round(before_dpwd*multiplier_dpwd) as before_dpwd , round(after_dpwd*multiplier_dpwd) as after_dpwd ,  status_dpwd,
			create_dpwd, to_char(COALESCE(createdate_dpwd,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			update_dpwd, to_char(COALESCE(updatedate_dpwd,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + tbl_trx_dpwd + `  
			ORDER BY createdate_dpwd DESC   `

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			iddpwd_db, date_dpwd_db, idcurr_db                                                                                                                              string
			tipedocuser_dpwd_db, tipedoc_dpwd_db, tipeakun_dpwd_db, idagenmember_db, ipaddress_dpwd_db, timezone_dpwd_db, bank_in_info_db, bank_out_info_db, status_dpwd_db string
			bank_in_db, bank_out_db                                                                                                                                         int
			amount_dpwd_db, before_dpwd_db, after_dpwd_db                                                                                                                   float64
			create_dpwd_db, createdate_dpwd_db, update_dpwd_db, updatedate_dpwd_db                                                                                          string
		)

		err = row.Scan(&iddpwd_db, &date_dpwd_db, &idcurr_db,
			&tipedocuser_dpwd_db, &tipedoc_dpwd_db, &tipeakun_dpwd_db, &idagenmember_db, &ipaddress_dpwd_db, &timezone_dpwd_db,
			&bank_in_db, &bank_in_info_db, &bank_out_db, &bank_out_info_db,
			&amount_dpwd_db, &before_dpwd_db, &after_dpwd_db, &status_dpwd_db,
			&create_dpwd_db, &createdate_dpwd_db, &update_dpwd_db, &updatedate_dpwd_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if create_dpwd_db != "" {
			create = create_dpwd_db + ", " + createdate_dpwd_db
		}
		if update_dpwd_db != "" {
			update = update_dpwd_db + ", " + updatedate_dpwd_db
		}
		switch status_dpwd_db {
		case "PROCESS":
			status_css = configs.STATUS_RUNNING
		case "APPROVED":
			status_css = configs.STATUS_COMPLETE
		case "REJECT":
			status_css = configs.STATUS_CANCEL
		}

		obj.Transdpwd_id = iddpwd_db
		obj.Transdpwd_date = date_dpwd_db
		obj.Transdpwd_idcurr = idcurr_db
		obj.Transdpwd_tipedoc = tipedoc_dpwd_db
		obj.Transdpwd_tipeuserdoc = tipedocuser_dpwd_db
		obj.Transdpwd_tipeakun = tipeakun_dpwd_db
		obj.Transdpwd_idmember = idagenmember_db
		obj.Transdpwd_nmmember = _GetInfoMember(idmasteragen, idagenmember_db)
		obj.Transdpwd_ipaddress = ipaddress_dpwd_db
		obj.Transdpwd_timezone = timezone_dpwd_db
		obj.Transdpwd_bank_in = bank_in_db
		obj.Transdpwd_bank_in_info = bank_in_info_db
		obj.Transdpwd_bank_out = bank_out_db
		obj.Transdpwd_bank_out_info = bank_out_info_db
		obj.Transdpwd_amount = amount_dpwd_db
		obj.Transdpwd_before = before_dpwd_db
		obj.Transdpwd_after = after_dpwd_db
		obj.Transdpwd_status = status_dpwd_db
		obj.Transdpwd_status_css = status_css
		obj.Transdpwd_create = create
		obj.Transdpwd_update = update
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
func Save_transdpwd(admin, idrecord, idmasteragen, idmaster, tipedoc, idmember, note_dpwd, status, sData string, bank_in, bank_out int, amount float32) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	_, _, tbl_trx_dpwd, _ := Get_mappingdatabase(idmasteragen)

	idcurr := _GetDefaultCurr(idmasteragen)
	multiplier := _GetMultiplier(idcurr)
	before := 0
	after := 0
	if sData == "New" {
		sql_insert := `
				insert into
				` + tbl_trx_dpwd + ` (
					iddpwd , idmasteragen, idmaster, 
					yearmonth_dpwd , date_dpwd, idcurr, tipedocuser_dpwd, tipedoc_dpwd, tipeakun_dpwd, idagenmember, 
					bank_in, bank_in_info , bank_out, bank_out_info, 
					multiplier_dpwd, amountdefault_dpwd, amount_dpwd, before_dpwd, after_dpwd, status_dpwd, note_dpwd, 
					create_dpwd, createdate_dpwd  
				) values (
					$1, $2, $3,   
					$4, $5, $6, $7, $8, $9, $10,    
					$11, $12, $13, $14,     
					$15, $16, $17, $18, $19, $20, $21,      
					$22, $23
				)
			`
		tipeakun_dpwd := ""
		temp_bank_out := ""
		temp_bank_in := ""
		switch tipedoc {
		case "DEPOSIT":
			tipeakun_dpwd = "IN"
			temp_bank_out = _GetInfoBank(idmasteragen, idmember, "MEMBER", bank_out)
			temp_bank_in = _GetInfoBank(idmasteragen, idmember, "AGEN", bank_in)
		case "WITHDRAW":
			tipeakun_dpwd = "OUT"
			temp_bank_out = _GetInfoBank(idmasteragen, idmember, "AGEN", bank_out)
			temp_bank_in = _GetInfoBank(idmasteragen, idmember, "MEMBER", bank_in)
		case "BONUS":
			tipeakun_dpwd = "OUT"
			temp_bank_out = _GetInfoBank(idmasteragen, idmember, "AGEN", bank_out)
			temp_bank_in = _GetInfoBank(idmasteragen, idmember, "MEMBER", bank_in)
			tipeakun_dpwd = "OUT"
		}
		amount_db := amount / multiplier

		field_column := tbl_trx_dpwd + tglnow.Format("YYYY-MM")
		idrecord_counter := Get_counter(field_column)
		iddpwd := idmasteragen + "DPWD" + tglnow.Format("YY") + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idrecord_counter)

		flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_trx_dpwd, "INSERT",
			iddpwd, idmasteragen, idmaster,
			tglnow.Format("YYYY-MM"), tglnow.Format("YYYY-MM-DD"), idcurr, "A", tipedoc, tipeakun_dpwd, idmember,
			bank_in, temp_bank_in, bank_out, temp_bank_out,
			multiplier, amount, amount_db, before, after, status, note_dpwd,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		if tipedoc == "WITDRAW" {
			sql_update := `
				UPDATE 
				` + tbl_trx_dpwd + `  
				SET bank_out=$1, bank_out_info=$2,    
				update_dpwd=$3, updatedate_dpwd=$4      
				WHERE iddpwd=$5  AND idmasteragen=$6   
			`
			temp_bank_out := _GetInfoBank(idmasteragen, idmember, "AGEN", bank_out)

			flag_update, msg_update := Exec_SQL(sql_update, tbl_trx_dpwd, "UPDATE",
				bank_out, temp_bank_out,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord, idmasteragen)

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
func Update_statustransdpwd(admin, idrecord, idmasteragen, idmaster, idmember, note, status string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	tbl_mst_member, _, tbl_trx_dpwd, _ := Get_mappingdatabase(idmasteragen)
	idcurr := _GetDefaultCurr(idmasteragen)
	if status == "APPROVED" {
		tipe, status_db, _, amount := _GetDepoWd(idmasteragen, idrecord)
		if status_db == "PROCESS" {
			flag := true
			cash_in, cash_out := _GetMemberCredit(idmasteragen, idmember, tbl_mst_member)

			var credit_member float64 = cash_in - cash_out
			var before float64 = 0
			var after float64 = 0
			fmt.Println("CREDIT :", credit_member)
			tipeakun := ""
			switch tipe {
			case "DEPOSIT":
				before = credit_member
				after = before + amount
				tipeakun = "IN"
			case "WITHDRAW":
				before = credit_member
				after = before - amount
				tipeakun = "OUT"
				if amount > credit_member {
					flag = false
					fmt.Println("The Amount exceed Credit Member")
				}
			}
			if flag {
				sql_update := `
					UPDATE 
					` + tbl_trx_dpwd + `  
					SET status_dpwd=$1, before_dpwd=$2, after_dpwd=$3,  
					update_dpwd=$4, updatedate_dpwd=$5 
					WHERE iddpwd=$6  AND idmasteragen=$7   
				`
				flag_update, msg_update := Exec_SQL(sql_update, tbl_trx_dpwd, "UPDATE",
					status, before, after,
					admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord, idmasteragen)

				if flag_update {
					msg = "Succes"
					flag := _Save_transaksi(admin, idrecord, idmasteragen, idmaster, idcurr, "TRANSAKSI", tipeakun, idmember, amount, before, after)
					if flag {
						_Save_creditmember(admin, idmember, idmasteragen, tipeakun, tbl_mst_member, amount)
						//UPDATE CREDIT AGEN + MASTER
						_Save_creditagenandmaster(admin, idmasteragen, tipeakun, amount)
					}
				} else {
					fmt.Println(msg_update)
				}
			}
		}
	} else if status == "REJECTED" {
		_, status_db, _, _ := _GetDepoWd(idmasteragen, idrecord)
		if status_db == "PROCESS" {
			sql_update := `
				UPDATE 
				` + tbl_trx_dpwd + `  
				SET status_dpwd=$1, note_dpwd=$2,  
				update_dpwd=$3, updatedate_dpwd=$4  
				WHERE iddpwd=$5  AND idmasteragen=$6   
			`

			flag_update, msg_update := Exec_SQL(sql_update, tbl_trx_dpwd, "UPDATE",
				status, note,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord, idmasteragen)

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

func _Save_transaksi(admin, iddpwd, idmasteragen, idmaster, idcurr, tipedoc, tipeakun, idmember string, amount, before, after float64) bool {
	tglnow, _ := goment.New()
	_, _, tbl_trx_dpwd, tbl_trx_transaksi := Get_mappingdatabase(idmasteragen)
	flag := false
	sql_insert := `
				insert into
				` + tbl_trx_transaksi + ` (
					idtransaksi , idother, idmasteragen, idmaster, 
					yearmonth_transaksi , date_transaksi, idcurr, tipedoc_transaksi, tipeakun_transaksi, idagenmember, 
					amount_transaksi, before_transaksi, after_transaksi,  
					create_transaksi, createdate_transaksi  
				) values (
					$1, $2, $3, $4,  
					$5, $6, $7, $8, $9, $10, 
					$11, $12, $13,          
					$14, $15 
				)
			`

	field_column := tbl_trx_transaksi + tglnow.Format("YYYY-MM")
	idrecord_counter := Get_counter(field_column)
	idtransaksi := idmasteragen + "TRANS" + tglnow.Format("YY") + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idrecord_counter)

	flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_trx_transaksi, "INSERT",
		idtransaksi, iddpwd, idmasteragen, idmaster,
		tglnow.Format("YYYY-MM"), tglnow.Format("YYYY-MM-DD"), idcurr, tipedoc, tipeakun, idmember,
		amount, before, after,
		admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

	if flag_insert {
		//update document tbl_trx_dpwd
		sql_update := `
				UPDATE 
				` + tbl_trx_dpwd + `  
				SET connect_idtransaksi=$1, connect_createtransaksi=$2,  connect_createdatetransaksi=$3, 
				update_dpwd=$4, updatedate_dpwd=$5   
				WHERE iddpwd=$6  AND idmasteragen=$7    
			`

		flag_update, msg_update := Exec_SQL(sql_update, tbl_trx_dpwd, "UPDATE",
			idtransaksi, "SYSTEM", tglnow.Format("YYYY-MM-DD HH:mm:ss"),
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), iddpwd, idmasteragen)

		if flag_update {
			flag = true
		} else {
			fmt.Println(msg_update)
		}
	} else {
		fmt.Println(msg_insert)
	}
	return flag
}
func _Save_creditmember(admin, idmember, idmasteragen, tipe, table string, amount float64) {
	tglnow, _ := goment.New()

	c_in_db, c_out_db := _GetMemberCredit(idmasteragen, idmember, table)
	var c_in float64 = 0
	var c_out float64 = 0
	if tipe == "IN" {
		c_in = c_in_db + amount
		c_out = c_out_db
	} else {
		c_in = c_in_db
		c_out = c_out_db + amount
	}
	sql_update := `
		UPDATE 
		` + table + `  
		SET cashin_member=$1, cashout_member=$2, 
		update_member=$3, updatedate_member=$4 
		WHERE idmember=$5  AND idmasteragen=$6   
	`

	flag_update, msg_update := Exec_SQL(sql_update, table, "UPDATE",
		c_in, c_out,
		admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmember, idmasteragen)
	if !flag_update {
		fmt.Println(msg_update)
	}
}
func _Save_creditagenandmaster(admin, idmasteragen, tipe string, amount float64) {
	tglnow, _ := goment.New()

	c_in_db, c_out_db, idmaster_db := _GetAgenCredit(idmasteragen)
	var c_in float64 = 0
	var c_out float64 = 0
	if tipe == "IN" {
		c_in = c_in_db + amount
		c_out = c_out_db
	} else {
		c_in = c_in_db
		c_out = c_out_db + amount
	}
	sql_update := `
		UPDATE 
		` + configs.DB_tbl_mst_master_agen + `  
		SET agencredit_in=$1, agencredit_out=$2, 
		updatemasteragen=$3, updatedatemasteragen=$4 
		WHERE idmasteragen=$5   
	`

	flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_master_agen, "UPDATE",
		c_in, c_out,
		admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmasteragen)
	if flag_update {
		master_in_db, master_out_db := _GetMasterCredit(idmaster_db)
		var master_in float64 = 0
		var master_out float64 = 0
		if tipe == "IN" {
			master_in = master_in_db + amount
			master_out = master_out_db
		} else {
			master_in = master_in_db
			master_out = master_out_db + amount
		}
		sql_updatemaster := `
			UPDATE 
			` + configs.DB_tbl_mst_master + `  
			SET mastercredit_in=$1, mastercredit_out=$2, 
			updatemaster=$3, updatedatemaster=$4 
			WHERE idmaster=$5   
		`
		flag_updatemaster, msg_updatemaster := Exec_SQL(sql_updatemaster, configs.DB_tbl_mst_master, "UPDATE",
			master_in, master_out,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmaster_db)
		if flag_updatemaster {
		} else {
			fmt.Println(msg_updatemaster)
		}
	} else {
		fmt.Println(msg_update)
	}
}
func _GetDefaultCurr(idrecord string) string {
	con := db.CreateCon()
	ctx := context.Background()
	idcurr := ""

	sql_select := `SELECT
		idcurr   
		FROM ` + configs.DB_tbl_mst_master_agen + `  
		WHERE idmasteragen = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idrecord)
	switch e := row.Scan(&idcurr); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return idcurr
}
func _GetMultiplier(idrecord string) float32 {
	con := db.CreateCon()
	ctx := context.Background()
	multipliercurr := 0

	sql_select := `SELECT
		multipliercurr   
		FROM ` + configs.DB_tbl_mst_curr + `  
		WHERE idcurr = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idrecord)
	switch e := row.Scan(&multipliercurr); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return float32(multipliercurr)
}
func _GetInfoBank(idmasteragen, idagenmember, tipe string, idrecord int) string {
	con := db.CreateCon()
	ctx := context.Background()
	info := ""
	bank_id := ""
	bank_norek := ""
	bank_nmrek := ""
	sql_select := ""
	if tipe == "AGEN" {
		sql_select = `SELECT
			idbanktype, norekbank, nmownerbank   
			FROM ` + configs.DB_tbl_mst_master_agen_bank + `  
			WHERE idagenbank=` + strconv.Itoa(idrecord) + ` AND idmasteragen='` + idmasteragen + `'    
		`
	} else {
		sql_select = `SELECT
			idbanktype, norekbank_agenmemberbank, nmownerbank_agenmemberbank   
			FROM ` + configs.DB_tbl_mst_master_agen_member_bank + `  
			WHERE idagenmemberbank=` + strconv.Itoa(idrecord) + ` AND idagenmember='` + idagenmember + `'    
		`
	}

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&bank_id, &bank_norek, &bank_nmrek); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	info = bank_id + "-" + bank_norek + "-" + bank_nmrek
	return info
}
func _GetInfoMember(idmasteragen, idagenmember string) string {
	con := db.CreateCon()
	ctx := context.Background()
	username_agenmember_db := ""
	name_agenmember_db := ""

	sql_select := `SELECT
		username_agenmember, name_agenmember   
		FROM ` + configs.DB_tbl_mst_master_agen_member + `  
		WHERE idagenmember=$1 AND idmasteragen=$2
	`
	row := con.QueryRowContext(ctx, sql_select, idagenmember, idmasteragen)
	switch e := row.Scan(&username_agenmember_db, &name_agenmember_db); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return username_agenmember_db + "-" + name_agenmember_db
}
func _GetDepoWd(idmasteragen, iddpwd string) (string, string, float64, float64) {
	con := db.CreateCon()
	ctx := context.Background()
	tipedoc_db := ""
	status_db := ""
	multiplier_db := 0
	amount_db := 0

	_, _, tbl_trx_dpwd, _ := Get_mappingdatabase(idmasteragen)

	sql_select := `SELECT
		tipedoc_dpwd,status_dpwd, multiplier_dpwd, amount_dpwd   
		FROM ` + tbl_trx_dpwd + `  
		WHERE iddpwd=$1 AND idmasteragen=$2
	`
	row := con.QueryRowContext(ctx, sql_select, iddpwd, idmasteragen)
	switch e := row.Scan(&tipedoc_db, &status_db, &multiplier_db, &amount_db); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return tipedoc_db, status_db, float64(multiplier_db), float64(amount_db)
}
func _GetMemberCredit(idmasteragen, idmember, table string) (float64, float64) {
	con := db.CreateCon()
	ctx := context.Background()
	cash_in_db := 0
	cash_out_db := 0

	sql_select := `SELECT
		cashin_member, cashout_member    
		FROM ` + table + `  
		WHERE idmember=$1 AND idmasteragen=$2 
	`
	row := con.QueryRowContext(ctx, sql_select, idmember, idmasteragen)
	switch e := row.Scan(&cash_in_db, &cash_out_db); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return float64(cash_in_db), float64(cash_out_db)
}
func _GetAgenCredit(idmasteragen string) (float64, float64, string) {
	con := db.CreateCon()
	ctx := context.Background()
	cash_in_db := 0
	cash_out_db := 0
	idmaster_db := ""

	sql_select := `SELECT
		agencredit_in, agencredit_out, idmaster    
		FROM ` + configs.DB_tbl_mst_master_agen + `  
		WHERE idmasteragen=$1 
	`
	row := con.QueryRowContext(ctx, sql_select, idmasteragen)
	switch e := row.Scan(&cash_in_db, &cash_out_db, &idmaster_db); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return float64(cash_in_db), float64(cash_out_db), idmaster_db
}
func _GetMasterCredit(idmaster string) (float64, float64) {
	con := db.CreateCon()
	ctx := context.Background()
	cash_in_db := 0
	cash_out_db := 0

	sql_select := `SELECT
		mastercredit_in, mastercredit_out   
		FROM ` + configs.DB_tbl_mst_master + `  
		WHERE idmaster=$1 
	`
	row := con.QueryRowContext(ctx, sql_select, idmaster)
	switch e := row.Scan(&cash_in_db, &cash_out_db); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return float64(cash_in_db), float64(cash_out_db)
}
