package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type VoteHelpfulRatingTxArgs struct {
	RatingID       uuid.UUID `json:"rating_id" binding:"arguired,uuid"`
	Helpful        bool      `json:"helpful" binding:"arguired"`
	UserID         uuid.UUID `json:"user_id" binding:"arguired,uuid"`
	HelpfulVotes   int32
	UnhelpfulVotes int32
}

func (s *pgRepo) VoteHelpfulRatingTx(ctx context.Context, arg VoteHelpfulRatingTxArgs) (id uuid.UUID, err error) {
	err = s.execTx(ctx, func(q *Queries) (err error) {
		existedVote, err := q.GetRatingVote(ctx, GetRatingVoteParams{
			RatingID: arg.RatingID,
			UserID:   arg.UserID,
		})
		isUpdate := true
		if err != nil {
			if errors.Is(err, ErrRecordNotFound) {
				isUpdate = false
				// add rating helpful record
				existedVote, err = q.InsertRatingVotes(ctx, InsertRatingVotesParams{
					RatingID:  arg.RatingID,
					UserID:    arg.UserID,
					IsHelpful: arg.Helpful,
				})
				if err != nil {
					return
				}
			} else {
				return
			}
		}

		if existedVote.IsHelpful != arg.Helpful {
			// update rating helpful record
			_, err = q.UpdateRatingVote(ctx, UpdateRatingVoteParams{
				ID:        existedVote.ID,
				IsHelpful: &arg.Helpful,
			})
			if err != nil {
				return
			}
		}

		if err != nil {
			return
		}

		updateRatingParams := UpdateProductRatingParams{
			ID: arg.RatingID,
		}
		if arg.Helpful {
			updateRatingParams.HelpfulVotes = utils.Int32Ptr(arg.HelpfulVotes + 1)
			if isUpdate {
				updateRatingParams.UnhelpfulVotes = utils.Int32Ptr(arg.UnhelpfulVotes - 1)
			}
		} else {
			updateRatingParams.UnhelpfulVotes = utils.Int32Ptr(arg.UnhelpfulVotes + 1)
			if isUpdate {
				updateRatingParams.HelpfulVotes = utils.Int32Ptr(arg.HelpfulVotes - 1)
			}
		}
		_, err = q.UpdateProductRating(ctx, updateRatingParams)
		if err != nil {
			return
		}
		id = existedVote.ID
		return
	})

	return
}
