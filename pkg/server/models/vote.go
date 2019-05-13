package models

import (
	"DB_Project_TP/api"
	"database/sql"
	"log"
	"net/http"

	"DB_Project_TP/config"
)

//toDo доделать логику обновления голосований
func VoteBranch(vote api.Vote, id uint64) (status int, diff int64) {
	if ok := CheckUser(vote.Nickname); !ok {
		status = http.StatusNotFound

		return
	}

	diff = int64(vote.Voice)

	tx, _ := config.DB.Begin()

	if ok := CheckUserVoteInThread(vote.Nickname, id); !ok {
		SaveUserVote(vote, id)
	} else {
		prevDiff := UpdateUserVote(vote, id)
		diff -= prevDiff
		if prevDiff == 0 {
			status = http.StatusOK

			return
		}
	}

	res, err := config.DB.Exec(sqlVoteForThread, diff, id)
	if err != nil {
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	rows, _ := res.RowsAffected()
	if rows != 1 {
		tx.Rollback()
		config.Logger.Fatalf("VoteBranch update: expected: %v have %v", 1, rows)
	}

	tx.Commit()
	status = http.StatusOK

	return
}

func SaveUserVote(vote api.Vote, id uint64) {
	res, err := config.DB.Exec(sqlSaveUserVote, vote.Nickname, id, vote.Voice)
	if err != nil {
		log.Fatalln("SaveUserVote:", err.Error())
	}

	rows, _ := res.RowsAffected()
	if rows != 1 {
		log.Fatalln("SaveUserVote: expected %v, have %v", 1, rows)
	}
}

func CheckUserVoteInThread(nickname string, threadID uint64) bool {
	row := config.DB.QueryRow(sqlCheckUserVote, nickname, threadID)
	err := row.Scan()
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func UpdateUserVote(vote api.Vote, id uint64) (prevDiff int64) {
	row := config.DB.QueryRow(sqlCheckUserVote, vote.Nickname, id)
	err := row.Scan(&prevDiff)
	if err != nil {
		log.Fatalln("UpdateUserVote: ", err.Error())
	}

	_, err = config.DB.Exec(sqlUpdateUserVote, vote.Voice, id, vote.Nickname)
	if err != nil {
		log.Fatalln("UpdateUserVote: ", err.Error())
	}

	return prevDiff
}

var sqlVoteForThread = `update threads set votes = votes + $1 where id = $2`

var sqlSaveUserVote = `insert into votes (v_user, thread, u_vote) values ($1, $2, $3)`

var sqlCheckUserVote = `select u_vote from votes  where v_user = $1 and thread = $2`

var sqlUpdateUserVote = `update votes set u_vote = $1 where thread = $2 and v_user = $3`
