package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type VoteHelpfulRatingTxArgs struct {
	RatingID uuid.UUID `json:"ratingId" validate:"required,uuid"`
	Helpful  bool      `json:"helpful" validate:"required"`
	UserID   uuid.UUID `json:"userId" validate:"required,uuid"`
}

func (repo *pgRepo) VoteHelpfulRatingTx(ctx context.Context, arg VoteHelpfulRatingTxArgs) (id uuid.UUID, err error) {
	err = repo.execTx(ctx, func(q *Queries) (err error) {
		existedVote, err := q.GetRatingVote(ctx, GetRatingVoteParams{
			RatingID: arg.RatingID,
			UserID:   arg.UserID,
		})
		if err != nil {
			if errors.Is(err, ErrRecordNotFound) {
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

		// The vote record has been updated successfully
		// Helpful vote counts should be calculated from the rating_votes table
		// rather than stored as denormalized data on the rating record
		id = existedVote.ID
		return
	})

	return
}
