package models

import (
	"DB_Project_TP/api"
	"database/sql"
	"log"
	"net/http"

	"DB_Project_TP/config"
)

//toDo доделать логику обновления голосований
func VoteBranch(vote api.Vote, id int) (status int) {
	if ok := CheckUser(vote.Nickname); !ok {
		return http.StatusNotFound
	}

	diffVote := vote.Voice

	tx, _ := config.DB.Begin()

	if ok := CheckUserVoteInThread(vote.Nickname, id); !ok {
		SaveUserVote(vote, id)
	} else {
		prevDiff := UpdateUserVote(vote, id)
		if prevDiff == diffVote {
			return http.StatusOK
		}

		diffVote -= prevDiff
	}

	res, err := config.DB.Exec(sqlVoteForThread, diffVote, id)
	if err != nil {
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	rows, _ := res.RowsAffected()
	if rows != 1 {
		tx.Rollback()
		config.Logger.Fatalf("VoteBranch update: expected: %v have %v", 1, rows)
	}

	tx.Commit()

	return http.StatusOK
}

func SaveUserVote(vote api.Vote, id int) {
	res, err := config.DB.Exec(sqlSaveUserVote, vote.Nickname, id, vote.Voice)
	if err != nil {
		log.Fatalln("SaveUserVote:", err.Error())
	}

	rows, _ := res.RowsAffected()
	if rows != 1 {
		log.Fatalln("SaveUserVote: expected %v, have %v", 1, rows)
	}
}

func CheckUserVoteInThread(nickname string, threadID int) bool {
	row := config.DB.QueryRow(sqlCheckUserVote, nickname, threadID)
	err := row.Scan()
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func UpdateUserVote(vote api.Vote, id int) (prevDiff int) {
	row := config.DB.QueryRow(sqlCheckUserVote, vote.Nickname, id)
	err := row.Scan(&prevDiff)
	if err != nil {
		log.Fatalln("UpdateUserVote: ", err.Error())
	}

	_, err = config.DB.Exec(sqlUpdateUserVote, vote.Voice, id)
	if err != nil {
		log.Fatalln("UpdateUserVote: ", err.Error())
	}

	return prevDiff
}

var sqlVoteForThread = `update threads set votes = votes + $1 where id = $2`

var sqlSaveUserVote = `insert into votes (v_user, thread, u_vote) values ($1, $2, $3)`

var sqlCheckUserVote = `select u_vote from votes  where v_user = $1 and thread = $2`

var sqlUpdateUserVote = `update votes set u_vote = $1 where thread = $2`
