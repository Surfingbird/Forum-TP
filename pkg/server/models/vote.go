package models

import (
	"DB_Project_TP/api"
	"log"
	"net/http"

	"DB_Project_TP/config"
)

func VoteBranch(vote api.Vote, id uint64) (status int, diff int) {
	newv := vote.Voice
	diff = newv

	// tx, err := config.DB.Begin()
	// if err != nil {
	// 	config.Logger.Fatal("VoteBranch", err.Error())
	// }

	// ok, old := CheckUserVoteInThread(vote.Nickname, id, tx)
	ok, old := CheckUserVoteInThread(vote.Nickname, id)
	if !ok {
		// SaveUserVote(vote, id, tx)
		SaveUserVote(vote, id)
	} else {
		// UpdateUserVote(vote, id, tx)
		UpdateUserVote(vote, id)
		diff = newv - old
		if diff == 0 {
			status = http.StatusOK

			return
		}
	}

	// res, err := tx.Exec(sqlVoteForThread, diff, id)
	res, err := config.DB.Exec(sqlVoteForThread, diff, id)
	if err != nil {
		// tx.Rollback()
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	rows := res.RowsAffected()
	if rows != 1 {
		// tx.Rollback()
		config.Logger.Fatalf("VoteBranch update: expected: %v have %v", 1, rows)
	}

	status = http.StatusOK

	// err = tx.Commit()

	return
}

// func SaveUserVote(vote api.Vote, id uint64, tx *pgx.Tx) {
func SaveUserVote(vote api.Vote, id uint64) {
	res, err := config.DB.Exec(sqlSaveUserVote, vote.Nickname, id, vote.Voice)
	if err != nil {
		log.Fatalln("SaveUserVote:", err.Error())
	}

	rows := res.RowsAffected()
	if rows != 1 {

		log.Fatalln("SaveUserVote: expected %v, have %v", 1, rows)
	}
}

// func CheckUserVoteInThread(nickname string, threadID uint64, tx *pgx.Tx) (ok bool, old int) {
func CheckUserVoteInThread(nickname string, threadID uint64) (ok bool, old int) {
	err := config.DB.QueryRow(sqlCheckUserVote,
		nickname,
		threadID).Scan(&old)

	if err != nil {
		return
	}

	ok = true

	return
}

// func UpdateUserVote(vote api.Vote, id uint64, tx *pgx.Tx) {
func UpdateUserVote(vote api.Vote, id uint64) {
	commandTag, err := config.DB.Exec(SqlUpdateUserVote, vote.Voice, vote.Nickname, id)
	if err != nil {
		config.Logger.Fatal("UpdateUserVote: ", err.Error())
	}

	if commandTag.RowsAffected() != 1 {
		config.Logger.Fatal("UpdateUserVote: ", err.Error())
	}
}

var sqlVoteForThread = `update threads set votes = votes + $1 where id = $2`

var sqlSaveUserVote = `insert into votes (v_user, thread, u_vote) values ($1, $2, $3)`

var sqlCheckUserVote = `select u_vote from votes  where v_user = $1 and thread = $2`

var SqlUpdateUserVote = `update votes set u_vote = $1 where v_user = $2 and thread = $3`
