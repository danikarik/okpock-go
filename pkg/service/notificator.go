package service

import (
	"errors"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/apns"
)

func (s *Service) getNotificator(passType api.PassType) (apns.Notificator, error) {
	switch passType {
	case api.BoardingPass:
		break
	case api.Coupon:
		return s.env.CouponNotificator, nil
	case api.EventTicket:
		break
	case api.Generic:
		break
	case api.StoreCard:
		break
	}
	return nil, errors.New("not supported yet")
}
