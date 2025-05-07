package repository

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// VerifyEmailTxArgs contains the parameters needed for the verify email transaction
type VerifyEmailTxArgs struct {
	VerifyEmail VerifyEmail
	VerifyCode  string
}

// VerifyEmailTx updates both the verify email record and user's verification status
// within a single database transaction to ensure consistency
func (s *pgRepo) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxArgs) error {
	err := s.execTx(ctx, func(q *Queries) error {
		// 1. Update the verify email record to mark it as used
		_, err := q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.VerifyEmail.ID,
			VerifyCode: arg.VerifyCode,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to update verify email")
			return err
		}

		// 2. Update the user's verified_email status to true
		trueVal := true
		_, err = q.UpdateUser(ctx, UpdateUserParams{
			ID:            arg.VerifyEmail.UserID,
			VerifiedEmail: &trueVal,
			UpdatedAt:     utils.GetPgTypeTimestamp(time.Now()),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to update user verified email status")
			return err
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("verify email transaction failed")
	}

	return err
}
