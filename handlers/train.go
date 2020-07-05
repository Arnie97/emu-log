package handlers

import (
	"net/http"
	"strings"

	"github.com/arnie97/emu-log/common"
	"github.com/go-chi/chi"
)

// singleTrainNoHandler returns the used vehicle and the corresponding date
// for the 30 most recent log items that matches the given train number.
func singleTrainNoHandler(w http.ResponseWriter, r *http.Request) {
	trainNo := chi.URLParam(r, "trainNo")
	rows, err := common.DB().Query(`
		SELECT *
		FROM emu_log
		WHERE train_no = ?
			OR train_no LIKE ?
			OR train_no LIKE ?
			OR train_no LIKE ?
		ORDER BY date DESC
		LIMIT 30;
	`, trainNo, trainNo+"/%", "%/"+trainNo+"/%", "%/"+trainNo)
	common.Must(err)
	defer rows.Close()
	serializeLogEntries(rows, w)
}

// multiTrainNoHandler returns the last used vehicle for the first train
// numbers in lexicographical order that matches the given fuzzy pattern.
func multiTrainNoHandler(w http.ResponseWriter, r *http.Request) {
	trainNoList := strings.Split(chi.URLParam(r, "trainNo"), ",")
	trainNoArgs := make([]interface{}, len(trainNoList))
	trainNoArgsPlaceHolder := strings.Repeat(", ?", len(trainNoList))[2:]
	for i := range trainNoList {
		trainNoArgs[i] = trainNoList[i]
	}
	rows, err := common.DB().Query(`
		SELECT date, emu_no, train_no
		FROM emu_latest
		WHERE train_no IN (`+trainNoArgsPlaceHolder+`)
	`, trainNoArgs...)
	common.Must(err)
	defer rows.Close()
	serializeLogEntries(rows, w)
}
