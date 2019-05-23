package models

import (
	"DB_Project_TP/api"

	"github.com/jackc/pgx"

	"DB_Project_TP/config"
)

func VoteBranch(vote api.Vote, id uint64) (sum int) {
	tx, err := config.DB.Begin()
	if err != nil {
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	_, err = tx.Exec(`DELETE FROM votes
	WHERE v_user = $1 AND thread = $2`,
		vote.Nickname, id)
	if err != nil {
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	_, err = tx.Exec(`INSERT INTO votes (v_user, u_vote, thread)
	VALUES 
		  ($1, $2, $3)`,
		vote.Nickname, vote.Voice, id)
	if err != nil {
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	err = tx.QueryRow(`SELECT votes
	 FROM threads 
	 WHERE id = $1`,
		id).Scan(&sum)
	if err != nil {
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		config.Logger.Fatal("VoteBranch", err.Error())
	}

	return sum
}

func UpdateUserVote(vote api.Vote, id uint64, tx *pgx.Tx) {
	commandTag, err := tx.Exec(SqlUpdateUserVote, vote.Voice, vote.Nickname, id)
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
