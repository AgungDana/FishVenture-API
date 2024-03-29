package bannerService

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	bannerconfig "github.com/e-fish/api/banner_http/banner_config"
	"github.com/e-fish/api/pkg/common/helper/logger"
	"github.com/e-fish/api/pkg/common/helper/savefile"
	"github.com/e-fish/api/pkg/common/helper/werror"
	"github.com/e-fish/api/pkg/domain/banner"
	"github.com/e-fish/api/pkg/domain/banner/model"
	"github.com/google/uuid"
)

type Service struct {
	conf bannerconfig.BannerConfig
	repo banner.Repo
}

func NewService(conf bannerconfig.BannerConfig) Service {
	var (
		service = Service{
			conf: conf,
		}
	)

	bannerRepo, err := banner.NewRepo(conf.BannerDBConfig)
	if err != nil {
		logger.Fatal("failed to create a new repo, can't create region service err causes failed create region repo: %v", err)
	}

	service.repo = bannerRepo

	return service
}

func (s *Service) ListBanner(ctx context.Context) ([]model.BannerOutput, error) {
	query := s.repo.NewQuery()
	result, err := query.ReadAllBanner(ctx)
	if err != nil {
		logger.ErrorWithContext(ctx, "failed get list banner with err [%v]", err)
	}
	return result, nil
}

func (s *Service) CreateBanner(ctx context.Context, input model.BannerInputCreate) (*uuid.UUID, error) {
	command := s.repo.NewCommand(ctx)
	result, err := command.CreateBanner(ctx, input)
	if err != nil {
		if err := command.Rollback(ctx); err != nil {
			logger.ErrorWithContext(ctx, "failed rollback create banner with err [%v]", err)
		}
		logger.ErrorWithContext(ctx, "failed create banner with err [%v]", err)
	}

	if err := command.Commit(ctx); err != nil {
		logger.ErrorWithContext(ctx, "failed commit create banner with err [%v]", err)
	}
	return result, nil
}

func (s *Service) UpdateBanner(ctx context.Context, input model.BannerInputUpdate) (*uuid.UUID, error) {
	command := s.repo.NewCommand(ctx)
	result, err := command.UpdateBanner(ctx, input)
	if err != nil {
		if err := command.Rollback(ctx); err != nil {
			logger.ErrorWithContext(ctx, "failed rollback update banner with id [%v] - err [%v]", input.ID, err)
		}
		logger.ErrorWithContext(ctx, "failed update banner with id [%v] - err [%v]", input.ID, err)
	}

	if err := command.Commit(ctx); err != nil {
		logger.ErrorWithContext(ctx, "failed commit update banner with id [%v] - err [%v]", input.ID, err)
	}
	return result, nil
}

func (s *Service) SaveImageBanner(ctx context.Context, file *multipart.FileHeader) (*UploadPhotoResponse, error) {
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext

	_, imageExtOk := savefile.ImageExt[strings.ReplaceAll(ext, ".", "")]

	if !imageExtOk {
		return nil, werror.Error{
			Code:    "FailedSaveFile",
			Message: "extension not suported",
			Details: map[string]any{
				"ext":                  ext,
				"permitted-extensions": fmt.Sprintf("%v", savefile.ImageExt),
			},
		}
	}

	err := savefile.SaveFile(file, s.conf.BannerImageConfig.Path+"/"+filename)

	if err != nil {
		return nil, err
	}

	result := UploadPhotoResponse{
		Name: filename,
		Url:  s.conf.BannerImageConfig.Url + filename,
	}

	return &result, nil
}
