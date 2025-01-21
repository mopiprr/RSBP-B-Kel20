package service

import (
	"context"

	"github.com/mci-its/backend-service/data-layer/entity"
	"github.com/mci-its/backend-service/data-layer/repository"
	"github.com/mci-its/backend-service/helpers/dto"
	helpers "github.com/mci-its/backend-service/helpers/function"
)

type (
	LogService interface {
		CreateLog(ctx context.Context, req dto.LogCreateRequest, userPkid int64) (dto.LogCreateResponse, error)
		GetLogByUser(ctx context.Context, userPkid int64) (dto.GetLogResponse, error)
		// GetLogByOffice(ctx context.Context, id int) (*entity.LogActivity, error)
		// GetLogByTargetTable(ctx context.Context, targetTable string) (*entity.LogActivity, error)
	}

	logService struct {
		logRepo  repository.LogRepository
		userRepo repository.UserRepository
	}
)

func NewLogService(logRepo repository.LogRepository, userRepo repository.UserRepository) LogService {
	return &logService{
		logRepo:  logRepo,
		userRepo: userRepo,
	}
}

func (s *logService) CreateLog(ctx context.Context, req dto.LogCreateRequest, userPkid int64) (dto.LogCreateResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userPkid)
	if err != nil {
		return dto.LogCreateResponse{}, err
	}

	ip := helpers.GetIpAdress()
	host := helpers.GetHostName()

	log := entity.LogActivity{
		UserPkid:    userPkid,
		IpAddress:   ip,
		ActionType:  req.ActionType,
		ActionTime:  req.ActionTime,
		TargetTable: req.TargetTable,
		TargetPkid:  req.TargetPkid,
		OfficePkid:  user.OfficePkid,
		CreatedBy:   req.CreatedBy,
		CreatedHost: host,
		UpdatedBy:   "",
		UpdatedHost: "",
		DeletedBy:   "",
		DeletedHost: "",
		IsDeleted:   false,
	}

	logCreate, err := s.logRepo.CreateLog(ctx, nil, log)
	if err != nil {
		return dto.LogCreateResponse{}, err
	}

	return dto.LogCreateResponse{
		Pkid:        logCreate.Pkid,
		UserPkid:    logCreate.UserPkid,
		IpAddress:   logCreate.IpAddress,
		OfficePkid:  logCreate.OfficePkid,
		ActionType:  logCreate.ActionType,
		ActionTime:  logCreate.ActionTime,
		TargetTable: logCreate.TargetTable,
		TargetPkid:  logCreate.TargetPkid,
	}, nil
}

func (s *logService) GetLogByUser(ctx context.Context, userPkid int64) (dto.GetLogResponse, error) {
	log, err := s.logRepo.GetLog(ctx, nil, userPkid)
	if err != nil {
		return dto.GetLogResponse{}, err
	}

	return dto.GetLogResponse{
		UserPkid:    log.UserPkid,
		IpAddress:   log.IpAddress,
		OfficePkid:  log.OfficePkid,
		ActionType:  log.ActionType,
		ActionTime:  log.ActionTime,
		TargetTable: log.TargetTable,
		TargetPkid:  log.TargetPkid,
	}, nil
}
