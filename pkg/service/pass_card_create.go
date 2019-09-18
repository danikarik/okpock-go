package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/filestore"
	"github.com/danikarik/okpock/pkg/pkpass"
	"github.com/danikarik/okpock/pkg/secure"
	"github.com/danikarik/okpock/pkg/store"
	uuid "github.com/satori/go.uuid"
)

// CreatePassCardRequest holds pass card info to be saved.
type CreatePassCardRequest struct {
	// Associated App Keys
	AppLaunchURL       string  `json:"appLaunchURL,omitempty"`
	AssociatedStoreIDs []int64 `json:"associatedStoreIdentifiers,omitempty"`

	// Companion App Keys
	UserInfo api.JSONMap `json:"userInfo,omitempty"`

	// Expiration Keys
	ExpirationDate string `json:"expirationDate,omitempty"`
	Voided         bool   `json:"voided,omitempty"`

	// Relevance Keys
	Beacons      []*api.Beacon   `json:"beacons,omitempty"`
	Locations    []*api.Location `json:"locations,omitempty"`
	MaxDistance  int64           `json:"maxDistance,omitempty"`
	RelevantDate string          `json:"relevantDate,omitempty"`

	// Style Keys
	Structure *api.PassStructure `json:"structure,omitempty"`

	// Visual Appearance Keys
	Barcodes           []*api.Barcode `json:"barcodes,omitempty"`
	BackgroundColor    string         `json:"backgroundColor,omitempty"`
	ForegroundColor    string         `json:"foregroundColor,omitempty"`
	GroupingIdentifier string         `json:"groupingIdentifier,omitempty"`
	LabelColor         string         `json:"labelColor,omitempty"`
	LogoText           string         `json:"logoText,omitempty"`

	// NFC-Enabled Pass Keys
	NFC *api.NFC `json:"nfc,omitempty"`
}

func (s *Service) createPassCardHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	id, err := s.idFromRequest(r, "id")
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "IDFromRequest", err)
	}

	project, err := s.env.Logic.LoadProject(ctx, user, id)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadProject", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadProject", err)
	}

	var req CreatePassCardRequest
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	passcard, err := s.newProjectPassCard(&req, project)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "NewProjectPassCard", err)
	}

	err = passcard.IsValid()
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	err = s.env.Logic.SaveNewPassCard(ctx, project, passcard)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SaveNewPassCard", err)
	}

	err = s.env.PassKit.InsertPass(
		ctx,
		passcard.Data.SerialNumber,
		passcard.Data.AuthenticationToken,
		passcard.Data.PassTypeID)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "InsertPass", err)
	}

	upload, err := s.newPassUpload(ctx, project, passcard)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "NewPassUpload", err)
	}

	err = s.env.Storage.UploadFile(ctx, s.env.Config.PassesBucket, upload)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UploadFile", err)
	}

	return sendJSON(w, http.StatusCreated, M{
		"id":           passcard.ID,
		"serialNumber": passcard.Data.SerialNumber,
		"url": fmt.Sprintf(
			"%s/downloads/%s%s",
			s.hostURL(),
			passcard.Data.SerialNumber,
			pkpass.Extension,
		),
	})
}

func (s *Service) newProjectPassCard(req *CreatePassCardRequest, project *api.Project) (*api.PassCardInfo, error) {
	data := &api.PassCard{
		Description:         project.Description,
		FormatVersion:       1,
		OrganizationName:    project.OrganizationName,
		PassTypeID:          s.passTypeToString(project.PassType),
		SerialNumber:        uuid.NewV4().String(),
		TeamID:              s.env.Config.Certificates.Team,
		AppLaunchURL:        req.AppLaunchURL,
		AssociatedStoreIDs:  req.AssociatedStoreIDs,
		UserInfo:            req.UserInfo,
		ExpirationDate:      req.ExpirationDate,
		Voided:              req.Voided,
		Beacons:             req.Beacons,
		Locations:           req.Locations,
		MaxDistance:         req.MaxDistance,
		RelevantDate:        req.RelevantDate,
		Barcodes:            req.Barcodes,
		BackgroundColor:     req.BackgroundColor,
		ForegroundColor:     req.ForegroundColor,
		GroupingIdentifier:  req.GroupingIdentifier,
		LabelColor:          req.LabelColor,
		LogoText:            req.LogoText,
		WebServiceURL:       s.hostURL(),
		AuthenticationToken: secure.Token(),
		NFC:                 req.NFC,
	}

	data = s.setPassStructure(req, project.PassType, data)
	err := data.IsValid()
	if err != nil {
		return nil, err
	}

	return api.NewPassCardInfo(data), nil
}

func (s *Service) passTypeToString(passType api.PassType) string {
	domainLayout := "pass.com.okpock.%s"
	if s.env.Config.IsDevelopment() {
		domainLayout += "-dev"
	}

	switch passType {
	case api.BoardingPass:
		return fmt.Sprintf(domainLayout, "boardingpass")
	case api.Coupon:
		return fmt.Sprintf(domainLayout, "coupon")
	case api.EventTicket:
		return fmt.Sprintf(domainLayout, "eventticket")
	case api.Generic:
		return fmt.Sprintf(domainLayout, "generic")
	case api.StoreCard:
		return fmt.Sprintf(domainLayout, "storecard")
	}

	return ""
}

func (s *Service) setPassStructure(req *CreatePassCardRequest, passType api.PassType, passCard *api.PassCard) *api.PassCard {
	switch passType {
	case api.BoardingPass:
		passCard.BoardingPass = req.Structure
	case api.Coupon:
		passCard.Coupon = req.Structure
	case api.EventTicket:
		passCard.EventTicket = req.Structure
	case api.Generic:
		passCard.Generic = req.Structure
	case api.StoreCard:
		passCard.StoreCard = req.Structure
	}
	return passCard
}

func (s *Service) newPassUpload(ctx context.Context, project *api.Project, passCard *api.PassCardInfo) (*filestore.Object, error) {
	files := []pkpass.File{}

	pass, err := json.Marshal(passCard.Data)
	if err != nil {
		return nil, err
	}
	files = append(files, pkpass.NewFile(pkpass.PassFilename, pass))

	if project.BackgroundImage != "" {
		background, err := s.env.Storage.GetFile(ctx, s.env.Config.UploadBucket, project.BackgroundImage)
		if err != nil {
			return nil, err
		}
		files = append(files, pkpass.NewFile(pkpass.BackgroundFilename, background.Body))
	}

	if project.FooterImage != "" {
		footer, err := s.env.Storage.GetFile(ctx, s.env.Config.UploadBucket, project.FooterImage)
		if err != nil {
			return nil, err
		}
		files = append(files, pkpass.NewFile(pkpass.FooterFilename, footer.Body))
	}

	if project.IconImage != "" {
		icon, err := s.env.Storage.GetFile(ctx, s.env.Config.UploadBucket, project.IconImage)
		if err != nil {
			return nil, err
		}
		files = append(files, pkpass.NewFile(pkpass.IconFilename, icon.Body))
	}

	if project.LogoImage != "" {
		logo, err := s.env.Storage.GetFile(ctx, s.env.Config.UploadBucket, project.LogoImage)
		if err != nil {
			return nil, err
		}
		files = append(files, pkpass.NewFile(pkpass.LogoFilename, logo.Body))
	}

	if project.StripImage != "" {
		strip, err := s.env.Storage.GetFile(ctx, s.env.Config.UploadBucket, project.StripImage)
		if err != nil {
			return nil, err
		}
		files = append(files, pkpass.NewFile(pkpass.StripFilename, strip.Body))
	}

	manifest, err := pkpass.CreateManifest(files...)
	if err != nil {
		return nil, err
	}
	files = append(files, *manifest)

	if project.PassType != api.Coupon {
		return nil, errors.New("pkpass: unsupported signer")
	}

	signature, err := s.env.CouponSigner.Sign(manifest.Data)
	if err != nil {
		return nil, err
	}
	files = append(files, *signature)

	zip, err := pkpass.Zip(files...)
	if err != nil {
		return nil, err
	}

	return &filestore.Object{
		Key:         passCard.Data.SerialNumber,
		Body:        zip,
		ContentType: filestore.ApplePkpass,
	}, nil
}
