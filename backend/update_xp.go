package opengamifylms

import (
	"context"

	"github.com/puskunalis/opengamifylms/store/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func updateUserXPAndCheckThreshold(ctx context.Context, pool *pgxpool.Pool, courseStore *db.Queries, userID int64, xpToAdd int32) error {
	// Start a transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	// Ensure the transaction is rolled back if an error occurs
	defer tx.Rollback(ctx) //nolint:errcheck

	// Get current user XP
	user, err := courseStore.WithTx(tx).GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Calculate new XP
	newXP := user.Xp + xpToAdd

	// Update user XP
	err = courseStore.WithTx(tx).UpdateUserXp(ctx, db.UpdateUserXpParams{
		ID: userID,
		Xp: newXP,
	})
	if err != nil {
		return err
	}

	if newXP >= 100 { // TODO
		err = courseStore.WithTx(tx).AddUserBadge(ctx, db.AddUserBadgeParams{
			UserID:  userID,
			BadgeID: 2, // TODO
		})
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit(ctx)
}
