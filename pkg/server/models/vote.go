package models

import (
	"DB_Project_TP/api"
	"log"
	"net/http"

	"github.com/jackc/pgx"

	"DB_Project_TP/config"
)

//toDo доделать логику обновления голосований
func VoteBranch(vote api.Vote, id uint64) (status int, diff int) {
	newv := vote.Voice
	diff = newv

	tx, err := config.DB.Begin()
	if err != nil {
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	ok, old := CheckUserVoteInThread(vote.Nickname, id, tx)
	if !ok {
		SaveUserVote(vote, id, tx)
	} else {
		UpdateUserVote(vote, id, tx)
		diff = newv - old
		if diff == 0 {
			status = http.StatusOK

			return
		}
	}

	res, err := tx.Exec(sqlVoteForThread, diff, id)
	if err != nil {
		tx.Rollback()
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	rows := res.RowsAffected()
	if rows != 1 {
		tx.Rollback()
		config.Logger.Fatalf("VoteBranch update: expected: %v have %v", 1, rows)
	}

	status = http.StatusOK

	err = tx.Commit()

	return
}

func SaveUserVote(vote api.Vote, id uint64, tx *pgx.Tx) {
	res, err := tx.Exec(sqlSaveUserVote, vote.Nickname, id, vote.Voice)
	if err != nil {
		log.Fatalln("SaveUserVote:", err.Error())
	}

	rows := res.RowsAffected()
	if rows != 1 {

		log.Fatalln("SaveUserVote: expected %v, have %v", 1, rows)
	}
}

func CheckUserVoteInThread(nickname string, threadID uint64, tx *pgx.Tx) (ok bool, old int) {
	err := tx.QueryRow(sqlCheckUserVote,
		nickname,
		threadID).Scan(&old)

	if err != nil {

		return
	}

	ok = true

	return
}

func UpdateUserVote(vote api.Vote, id uint64, tx *pgx.Tx) {
	_, err := tx.Exec(sqlUpdateUserVote, vote.Voice, id, vote.Nickname)
	if err != nil {

		log.Fatalln("UpdateUserVote: ", err.Error())
	}
}

var sqlVoteForThread = `update threads set votes = votes + $1 where id = $2`

var sqlSaveUserVote = `insert into votes (v_user, thread, u_vote) values ($1, $2, $3)`

var sqlCheckUserVote = `select u_vote from votes  where v_user = $1 and thread = $2`

var sqlUpdateUserVote = `update votes set u_vote = $1 where v_user = $3 and thread = $2`
