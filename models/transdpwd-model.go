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

	tbl_trx_dpwd, _ := Get_mappingdatabase(idmasteragen)
	sql_select := `SELECT 
			iddpwd , date_dpwd, idcurr,  
			tipedoc_dpwd , tipeakun_dpwd, idagenmember,  note_bank, ipaddress_dpwd, 
			amount_dpwd , before_dpwd, after_dpwd,  status_dpwd,
			create_dpwd, to_char(COALESCE(createdate_dpwd,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			update_dpwd, to_char(COALESCE(updatedate_dpwd,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + tbl_trx_dpwd + `  
			ORDER BY createdate_dpwd DESC   `

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			iddpwd_db, date_dpwd_db, idcurr_db                                                                  string
			tipedoc_dpwd_db, tipeakun_dpwd_db, idagenmember_db, note_bank_db, ipaddress_dpwd_db, status_dpwd_db string
			amount_dpwd_db, before_dpwd_db, after_dpwd_db                                                       float32
			create_dpwd_db, createdate_dpwd_db, update_dpwd_db, updatedate_dpwd_db                              string
		)

		err = row.Scan(&iddpwd_db, &date_dpwd_db, &idcurr_db,
			&tipedoc_dpwd_db, &tipeakun_dpwd_db, &idagenmember_db, &note_bank_db, &ipaddress_dpwd_db,
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
		obj.Transdpwd_tipeakun = tipeakun_dpwd_db
		obj.Transdpwd_idmember = idagenmember_db
		obj.Transdpwd_notebank = note_bank_db
		obj.Transdpwd_ipaddress = ipaddress_dpwd_db
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
	tbl_trx_dpwd, _ := Get_mappingdatabase(idmasteragen)

	idcurr := _GetDefaultCurr(idmasteragen)
	multiplier := _GetMultiplier(idcurr)
	before := 0
	after := 0
	if sData == "New" {
		sql_insert := `
				insert into
				` + tbl_trx_dpwd + ` (
					iddpwd , idmasteragen, idmaster, 
					yearmonth_dpwd , date_dpwd, idcurr, tipedoc_dpwd, tipeakun_dpwd, idagenmember, 
					bank_int , bank_out, note_bank, multiplier_dpwd, amount_dpwd, before_dpwd, after_dpwd, status_dpwd, note_dpwd, 
					create_dpwd, createdate_dpwd  
				) values (
					$1, $2, $3,   
					$4, $5, $6, $7, $8, $9,   
					$10, $11, $12, $13, $14, $15, $16, $17, $18,     
					$19, $20
				)
			`
		tipeakun_dpwd := ""
		switch tipedoc {
		case "DEPOSIT":
			tipeakun_dpwd = "IN"
		case "WITHDRAW":
			tipeakun_dpwd = "OUT"
		case "BONUS":
			tipeakun_dpwd = "OUT"
		}

		note_bank := "FROM: BANK OUT - TO: BANK IN"
		field_column := tbl_trx_dpwd + tglnow.Format("YYYY-MM")
		idrecord_counter := Get_counter(field_column)
		iddpwd := idmasteragen + "-DPWD-" + tglnow.Format("YY") + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idrecord_counter)

		flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_trx_dpwd, "INSERT",
			iddpwd, idmasteragen, idmaster,
			tglnow.Format("YYYY-MM"), tglnow.Format("YYYY-MM-DD"), idcurr, tipedoc, tipeakun_dpwd, idmember,
			bank_in, bank_out, note_bank, multiplier, amount, before, after, status, note_dpwd,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + tbl_trx_dpwd + `  
				SET tipedoc_dpwd=$1, tipeakun_dpwd=$2,   
				idagenmember=$3, bank_int=$4, bank_out=$5, note_bank=$6,      
				amount_dpwd=$7, before_dpwd=$8, after_dpwd=$9, status_dpwd=$10, note_dpwd=$11,     
				update_dpwd=$12, updatedate_dpwd=$13     
				WHERE iddpwd=$14  AND idmasteragen=$15  
			`

		tipeakun_dpwd := ""
		switch tipedoc {
		case "DEPOSIT":
			tipeakun_dpwd = "IN"
		case "WITHDRAW":
			tipeakun_dpwd = "OUT"
		case "BONUS":
			tipeakun_dpwd = "OUT"
		}
		before := 0
		after := 0
		note_bank := "FROM: BANK OUT - TO: BANK IN"
		flag_update, msg_update := Exec_SQL(sql_update, tbl_trx_dpwd, "UPDATE",
			tipedoc, tipeakun_dpwd,
			idmember, bank_in, bank_out, note_bank,
			amount, before, after, status, note_dpwd,
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
		FROM ` + configs.DB_tbl_mst_cate_bank + `  
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