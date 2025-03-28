package services

import (
	"context"
	"fmt"
	"gopkg.in/guregu/null.v4/zero"
	"med/internal/db/psql/sqlc"
)

func (s *Service) ProcessEmailNotifications() error {
	fmt.Println("-----------------START--------------------")
	ctx := context.Background()
	notifications, err := s.storage.SelectPendingNotifications(ctx)
	if err != nil {
		fmt.Println("-------------------++++++++++++++++", err.Error())
		return nil
	}
	fmt.Println("-----------------------RESOURSE", notifications)

	for _, n := range notifications {
		fmt.Println("-----------------------BEGIN")
		err = s.email.SendConfirmation(n.Email.String, n.Message)
		status := "sent"
		if err != nil {
			status = "failed"
		}

		err = s.storage.UpdateNotificationStatus(ctx, sqlc.UpdateNotificationStatusParams{
			ID:     n.ID,
			Status: zero.StringFrom(status),
		})
		if err != nil {
			return nil
		}
	}
	return nil
}
