package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const dbSchema = `
	CREATE TABLE IF NOT EXISTS emu_latest (
		date        VARCHAR NOT NULL,
		emu_no      VARCHAR NOT NULL,
		train_no    VARCHAR NOT NULL,
		log_id      INTEGER NOT NULL,
		UNIQUE(train_no)
	);
	CREATE TABLE IF NOT EXISTS emu_log (
		date        VARCHAR NOT NULL,
		emu_no      VARCHAR NOT NULL,
		train_no    VARCHAR NOT NULL,
		UNIQUE(date, emu_no, train_no)
	);
	CREATE TABLE IF NOT EXISTS emu_qrcode (
		emu_no      VARCHAR NOT NULL,
		emu_bureau  CHAR(1) NOT NULL,
		emu_qrcode  VARCHAR NOT NULL,
		UNIQUE(emu_bureau, emu_qrcode)
	);
	CREATE INDEX IF NOT EXISTS idx_emu_no ON emu_log(emu_no);
`

type (
	LogEntry struct {
		Date      string `json:"date"`
		VehicleNo string `json:"emu_no"`
		TrainNo   string `json:"train_no"`
	}
	Bureau struct {
		Code       string
		Name       string
		BruteForce func(chan<- string)
		TrainNo    func(this *Bureau, qrCode string) (trainNo, date string, err error)
		VehicleNo  func(this *Bureau, qrCode string) (vehicleNo string, err error)
		Info       func(qrCode string) (info jsonObject, err error)
	}
	jsonObject map[string]interface{}
)

var bureaus = []Bureau{
	Bureau{
		Code: "H",
		Name: "中国铁路上海局集团有限公司",
		BruteForce: func(pqCodes chan<- string) {
			for i := 2000; i < 11000; i += 200 {
				pqCodes <- fmt.Sprintf("PQ%07d", i)
			}
			for i := 11000; i < 2000000; i += 500 {
				pqCodes <- fmt.Sprintf("PQ%07d", i)
			}
		},
		TrainNo: func(this *Bureau, pqCode string) (trainNo, date string, err error) {
			var info jsonObject
			info, err = this.Info(pqCode)
			if err == nil {
				defer catch(&err)
				trainNo = info["trainName"].(string)
				date = time.Now().Format("2006-01-02")
			}
			return
		},
		VehicleNo: func(this *Bureau, pqCode string) (vehicleNo string, err error) {
			var info jsonObject
			info, err = this.Info(pqCode)
			if err == nil {
				defer catch(&err)
				vehicleNo = normalizeVehicleNo(info["cdh"].(string))
			}
			return
		},
		Info: func(pqCode string) (info jsonObject, err error) {
			const api = "https://g.xiuxiu365.cn/railway_api/web/index/train"
			query := url.Values{"pqCode": {pqCode}}.Encode()
			resp, err := httpClient.Get(api + "?" + query)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			var result struct {
				Status int `json:"code"`
				Msg    string
				Data   jsonObject
			}
			err = parseResult(resp, &result)
			info = result.Data
			return
		},
	},
	Bureau{
		Code: "P",
		Name: "中国铁路北京局集团有限公司",
		BruteForce: func(qrCodes chan<- string) {
			for x := 1; x <= 5; x += 2 {
				for y := 0; y < 50; y++ {
					for z := 0; z < 10000; z += 500 {
						qrCodes <- fmt.Sprintf("%d%03d%04d", x, y, z)
					}
				}
			}
		},
		TrainNo: func(this *Bureau, qrCode string) (trainNo, date string, err error) {
			var info jsonObject
			info, err = this.Info(qrCode)
			if err == nil {
				defer catch(&err)
				trainNo = info["TrainnoId"].(string)
				date = info["TrainnoDate"].(string)
			}
			return
		},
		VehicleNo: func(this *Bureau, qrCode string) (vehicleNo string, err error) {
			var info jsonObject
			info, err = this.Info(qrCode)
			if err == nil {
				defer catch(&err)
				vehicleNo = normalizeVehicleNo(info["TrainId"].(string))
			}
			return
		},
		Info: func(qrCode string) (info jsonObject, err error) {
			const api = "https://aymaoto.jtlf.cn/webapi/otoshopping/ewh_getqrcodetrainnoinfo"
			const key = "qrcode=%s&key=ltRsjkiM8IRbC80Ni1jzU5jiO6pJvbKd"
			sign := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf(key, qrCode))))
			form := url.Values{"qrCode": {qrCode}, "sign": {sign}}
			resp, err := httpClient.PostForm(api, form)
			if err != nil {
				return
			}
			var result struct {
				Status int `json:"state"`
				Msg    string
				Data   struct {
					TrainInfo jsonObject
					URLStr    string
				}
			}
			err = parseResult(resp, &result)
			info = result.Data.TrainInfo
			return
		},
	},
}

var (
	httpClient = &http.Client{
		Timeout: requestTimeout,
	}
	wg sync.WaitGroup
	db *sql.DB
)

const (
	day            = 24 * time.Hour
	repeatInterval = time.Hour
	requestDelay   = 4 * time.Second
	requestTimeout = 5 * time.Second
	startTime      = 5 * time.Hour
	endTime        = 24 * time.Hour
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if len(os.Args) < 3 {
		log.Error().Msg("required command-line arguments: task, bureau codes")
		return
	}

	checkLocalTimezone()
	checkInternetConnection()
	checkDatabase()

	switch os.Args[1] {
	case "serve":
		go http.ListenAndServe("localhost:8080", newRouter())
		scheduleTask(func() {
			iterateBureaus((*Bureau).scanTrainNo, os.Args[2])
		})
	case "trainNo":
		iterateBureaus((*Bureau).scanTrainNo, os.Args[2])
	case "vehicleNo":
		iterateBureaus((*Bureau).scanVehicleNo, os.Args[2])
	case "info":
		if len(os.Args) < 4 {
			log.Error().Msg("missing argument: qr code")
			return
		}
		printInfo(os.Args[2], os.Args[3])
	default:
		log.Error().Msgf("invalid task option: %s", os.Args[1])
	}
}

func scheduleTask(task func()) {
	var nextRun time.Time
	for {
		now := time.Now()
		today := time.Date(
			now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0, time.Local,
		)
		if now.Before(today.Add(startTime)) {
			nextRun = today.Add(startTime)
		} else if now.After(today.Add(endTime - repeatInterval)) {
			nextRun = today.Add(startTime + day)
		} else {
			nextRun = now.Truncate(repeatInterval).Add(repeatInterval)
		}
		log.Info().Msgf("next scheduled run: %v", nextRun)
		time.Sleep(time.Until(nextRun))
		task()
	}
}

func printInfo(bureauCode, qrCode string) {
	for i := range bureaus {
		if bureaus[i].Code == bureauCode {
			info, _ := bureaus[i].Info(qrCode)
			prettyPrint(info)
			return
		}
	}
	log.Error().Msgf("unknown bureau code: %s", bureauCode)
}

func prettyPrint(obj interface{}) {
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	checkFatal(err)
	fmt.Printf("%s\n", jsonBytes)
}

func iterateBureaus(task func(*Bureau, *sql.Tx), bureauCodes string) {
	tx, err := db.Begin()
	checkFatal(err)
	defer tx.Rollback()

	for i := range bureaus {
		if bureauCodes == "" || strings.Contains(bureauCodes, bureaus[i].Code) {
			wg.Add(1)
			go task(&bureaus[i], tx)
		}
	}

	wg.Wait()
	tx.Commit()
}

func (b *Bureau) scanTrainNo(tx *sql.Tx) {
	log.Info().Msgf("[%s] job started: %s", b.Code, b.Name)
	defer wg.Done()

	rows, err := tx.Query(`
		SELECT emu_no, emu_qrcode, MAX(rowid)
		FROM emu_qrcode
		WHERE emu_bureau = ?
		GROUP BY emu_no
		ORDER BY emu_no ASC;
	`, b.Code)
	checkFatal(err)
	defer rows.Close()

	for rows.Next() {
		var e LogEntry
		var qrCode, id string
		checkFatal(rows.Scan(&e.VehicleNo, &qrCode, &id))
		time.Sleep(requestDelay)
		e.TrainNo, e.Date, err = b.TrainNo(b, qrCode)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		log.Debug().Msgf("[%s] %s -> %s", b.Code, e.VehicleNo, e.TrainNo)
		if e.TrainNo != "" {
			res, err := tx.Exec(
				`INSERT OR IGNORE INTO emu_log VALUES (?, ?, ?)`,
				e.Date, e.VehicleNo, e.TrainNo,
			)
			checkFatal(err)
			logID, err := res.LastInsertId()
			checkFatal(err)

			for _, singleTrainNo := range normalizeTrainNo(e.TrainNo) {
				_, err = tx.Exec(
					`REPLACE INTO emu_latest VALUES (?, ?, ?, ?)`,
					e.Date, e.VehicleNo, singleTrainNo, logID,
				)
				checkFatal(err)
			}
		}
	}
	log.Info().Msgf("job done: %s", b.Name)
}

func (b *Bureau) scanVehicleNo(tx *sql.Tx) {
	log.Info().Msgf("[%s] job started: %s", b.Code, b.Name)
	defer wg.Done()

	rows, err := tx.Query(`
		SELECT emu_qrcode
		FROM emu_qrcode
		WHERE emu_bureau = ?
		ORDER BY emu_qrcode ASC;
	`, b.Code)
	checkFatal(err)
	defer rows.Close()

	qrCodes := make(chan string)
	go func() {
		b.BruteForce(qrCodes)
		close(qrCodes)
	}()

	qrCodeFromDB := ""
	for qrCode := range qrCodes {
		// skip existing codes in the database
		for qrCode > qrCodeFromDB && rows.Next() {
			checkFatal(rows.Scan(&qrCodeFromDB))
			log.Debug().Msgf("[%s] loaded: %s", b.Code, qrCodeFromDB)
		}
		if qrCode == qrCodeFromDB {
			log.Debug().Msgf("[%s] skipped: %s", b.Code, qrCodeFromDB)
			continue
		}

		time.Sleep(requestDelay)
		vehicleNo, err := b.VehicleNo(b, qrCode)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		log.Debug().Msgf("[%s] checked: %s -> %s", b.Code, qrCode, vehicleNo)
		if vehicleNo != "" {
			_, err := tx.Exec(
				`INSERT OR IGNORE INTO emu_qrcode VALUES (?, ?, ?)`,
				vehicleNo, b.Code, qrCode,
			)
			checkFatal(err)
		}
	}
	log.Info().Msgf("[%s] job done: %s", b.Code, b.Name)
}

func catch(err *error) {
	if r := recover(); r != nil {
		*err = r.(error)
	}
}

func checkFatal(err error) {
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func checkLocalTimezone() {
	tzName, tzOffset := time.Now().Zone()
	if tzOffset*int(time.Second) != 8*int(time.Hour) {
		log.Warn().Msgf(
			"expected Beijing Timezone (UTC+08), but found %s (UTC%s)",
			tzName, time.Now().Format("-07"),
		)
	}
}

func checkInternetConnection() {
	start := time.Now()
	_, err := bureaus[0].Info("PQ0123456")
	checkFatal(err)
	log.Info().Msgf(
		"internet connection ok, round-trip delay %v",
		time.Since(start),
	)
}

func checkDatabase() {
	dbConn, err := sql.Open("sqlite3", "./emu_log.db")
	checkFatal(err)
	db = dbConn
	// TODO: defer db.Close()

	_, err = db.Exec(dbSchema)
	checkFatal(err)
	log.Info().Msgf(
		"found %d log records in the database",
		countRecords("emu_log"),
	)
	log.Info().Msgf(
		"found %d vehicles and %d qr codes in the database",
		countRecords("emu_qrcode", "DISTINCT emu_no"),
		countRecords("emu_qrcode"),
	)
}

func countRecords(tableName string, fields ...string) (count int) {
	field := "*"
	if len(fields) != 0 {
		field = fields[0]
	}
	row := db.QueryRow(fmt.Sprintf(
		`SELECT COUNT(%s) FROM %s`, field, tableName,
	))
	checkFatal(row.Scan(&count))
	return
}

func parseResult(resp *http.Response, resultPtr interface{}) (err error) {
	err = json.NewDecoder(resp.Body).Decode(resultPtr)
	apiResult := reflect.Indirect(reflect.ValueOf(resultPtr)).FieldByName
	if err == nil && apiResult("Status").Interface().(int) != 200 {
		err = fmt.Errorf(
			"api error %d: %s",
			apiResult("Status"), apiResult("Msg"),
		)
	}
	return
}

func newRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(requestTimeout),
	)
	mux.Get(`/map/{stationName}`, railMapHandler)
	mux.Get(`/train/{trainNo:[GDC]\d{1,4}}`, singleTrainNoHandler)
	mux.Get(`/train/{trainNo:.*,.*}`, multiTrainNoHandler)
	mux.Get(`/emu/{vehicleNo:[A-Z-0-9]*?\d{4}}`, singleVehicleNoHandler)
	mux.Get(`/emu/{vehicleNo:[A-Z-0-9]+}`, multiVehicleNoHandler)
	return mux
}

func railMapHandler(w http.ResponseWriter, r *http.Request) {
	const site = "http://cnrail.geogv.org"
	stationID := ""
	stationName := chi.URLParam(r, "stationName")
	defer func() {
		http.Redirect(w, r, fmt.Sprintf(
			"%s/zhcn/station/%s?useMapboxGl=true", site, stationID,
		), http.StatusSeeOther)
	}()

	keyword := stationName
	if len(stationName) > 2 {
		keyword = strings.TrimSuffix(stationName, "所")
	}
	resp, err := httpClient.Get(fmt.Sprintf(
		"%s/api/v1/match_feature/%s?locale=zhcn", site, keyword,
	))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	matches := struct {
		Success bool
		Data    [][3]string
	}{}
	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil || !matches.Success {
		return
	}

	for _, m := range matches.Data {
		itemID, itemType, itemName := m[0], m[1], m[2]
		if itemType != "STATION" {
			continue
		} else if strings.Replace(itemName, "线路所", "所", 1) == stationName {
			stationID = itemID
			return
		}
	}
}

func singleTrainNoHandler(w http.ResponseWriter, r *http.Request) {
	trainNo := chi.URLParam(r, "trainNo")
	rows, err := db.Query(`
		SELECT *
		FROM emu_log
		WHERE train_no = ?
			OR train_no LIKE ?
			OR train_no LIKE ?
			OR train_no LIKE ?
		ORDER BY date DESC
		LIMIT 30;
	`, trainNo, trainNo+"/%", "%/"+trainNo+"/%", "%/"+trainNo)
	checkFatal(err)
	defer rows.Close()
	serializeLogEntries(rows, w)
}

func multiTrainNoHandler(w http.ResponseWriter, r *http.Request) {
	trainNoList := strings.Split(chi.URLParam(r, "trainNo"), ",")
	trainNoArgs := make([]interface{}, len(trainNoList))
	trainNoArgsPlaceHolder := strings.Repeat(", ?", len(trainNoList))[2:]
	for i := range trainNoList {
		trainNoArgs[i] = trainNoList[i]
	}
	rows, err := db.Query(`
		SELECT date, emu_no, train_no
		FROM emu_latest
		WHERE train_no IN (`+trainNoArgsPlaceHolder+`)
	`, trainNoArgs...)
	checkFatal(err)
	defer rows.Close()
	serializeLogEntries(rows, w)
}

func singleVehicleNoHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT *
		FROM (
			SELECT *
			FROM emu_log
			WHERE emu_no IN (
				SELECT emu_no
				FROM emu_qrcode
				WHERE emu_no LIKE ?
			)
			ORDER BY date DESC
			LIMIT 30
		)
		ORDER BY emu_no ASC;
	`, "%"+normalizeVehicleNo(chi.URLParam(r, "vehicleNo")))
	checkFatal(err)
	defer rows.Close()
	serializeLogEntries(rows, w)
}

func multiVehicleNoHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT *
		FROM emu_log
		WHERE rowid IN (
			SELECT MAX(rowid)
			FROM emu_log
			WHERE emu_no IN (
				SELECT emu_no
				FROM emu_qrcode
				WHERE emu_no LIKE ?
			)
			GROUP BY emu_no
			LIMIT 30
		)
		ORDER BY emu_no ASC;
	`, "%"+normalizeVehicleNo(chi.URLParam(r, "vehicleNo"))+"%")
	checkFatal(err)
	defer rows.Close()
	serializeLogEntries(rows, w)
}

func normalizeTrainNo(trainNo string) (results []string) {
	trainNoRegExp := regexp.MustCompile(`\b[GDC]?\d{1,4}\b`)
	var initial string
	for i, part := range strings.Split(trainNo, "/") {
		if part = trainNoRegExp.FindString(part); len(part) == 0 {
			return
		} else if i == 0 && part[0] <= '9' {
			return
		} else if i == 0 {
			initial = part
		} else if omitted := len(initial) - len(part); omitted > 0 {
			part = initial[:omitted] + part
		}
		results = append(results, part)
	}
	return
}

func normalizeVehicleNo(vehicleNo string) string {
	return strings.ReplaceAll(vehicleNo, "-", "")
}

func serializeLogEntries(rows *sql.Rows, w http.ResponseWriter) {
	results := make([]LogEntry, 0)
	for rows.Next() {
		var e LogEntry
		checkFatal(rows.Scan(&e.Date, &e.VehicleNo, &e.TrainNo))
		results = append(results, e)
	}
	w.Header().Set("Content-Type", "application/json")
	checkFatal(json.NewEncoder(w).Encode(results))
}
