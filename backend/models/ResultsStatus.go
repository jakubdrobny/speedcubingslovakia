package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ResultsStatus struct {
	Id int `json:"id"`
	ApprovalFinished bool `json:"approvalFinished"`
	Approved bool `json:"approved"`
	Visible bool `json:"visible"`
	Displayname string `json:"displayname"`
}

func GetResultsStatus(db *pgxpool.Pool, statusId int) (ResultsStatus, error) {
	rows, err := db.Query(context.Background(), `SELECT rs.results_status_id, rs.approvalfinished, rs.approved, rs.visible, rs.displayname FROM results_status rs WHERE rs.results_status_id = $1;`, statusId);
	if err != nil {
		return ResultsStatus{}, err
	}

	var resultsStatus ResultsStatus
	found := false
	for rows.Next() {
		err = rows.Scan(&resultsStatus.Id, &resultsStatus.ApprovalFinished, &resultsStatus.Approved, &resultsStatus.Visible, &resultsStatus.Displayname)
		if err != nil {
			return ResultsStatus{}, err
		}
		found = true
	}

	if !found {
		return ResultsStatus{}, err
	}

	return resultsStatus, nil
}