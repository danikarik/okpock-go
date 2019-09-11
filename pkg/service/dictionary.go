package service

import (
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
)

func (s *Service) detectorTypesHandler(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, M{
		"data": []string{
			api.PKDataDetectorTypePhoneNumber,
			api.PKDataDetectorTypeLink,
			api.PKDataDetectorTypeAddress,
			api.PKDataDetectorTypeCalendarEvent,
		},
	})
}

func (s *Service) textAlignmentHandler(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, M{
		"data": []string{
			api.PKTextAlignmentLeft,
			api.PKTextAlignmentCenter,
			api.PKTextAlignmentRight,
			api.PKTextAlignmentNatural,
		},
	})
}

func (s *Service) dateStyleHandler(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, M{
		"data": []string{
			api.PKDateStyleNone,
			api.PKDateStyleShort,
			api.PKDateStyleMedium,
			api.PKDateStyleLong,
			api.PKDateStyleFull,
		},
	})
}

func (s *Service) numberStyleHandler(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, M{
		"data": []string{
			api.PKNumberStyleDecimal,
			api.PKNumberStylePercent,
			api.PKNumberStyleScientific,
			api.PKNumberStyleSpellOut,
		},
	})
}

func (s *Service) transitTypeHandler(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, M{
		"data": []string{
			api.PKTransitTypeAir,
			api.PKTransitTypeBoat,
			api.PKTransitTypeBus,
			api.PKTransitTypeGeneric,
			api.PKTransitTypeTrain,
		},
	})
}

func (s *Service) barcodeFormatHandler(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, M{
		"data": []string{
			api.PKBarcodeFormatQR,
			api.PKBarcodeFormatPDF417,
			api.PKBarcodeFormatAztec,
			api.PKBarcodeFormatCode128,
		},
	})
}
