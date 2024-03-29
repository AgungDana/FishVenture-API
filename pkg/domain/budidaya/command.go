package budidaya

import (
	"context"
	"time"

	"github.com/e-fish/api/pkg/common/helper/ctxutil"
	"github.com/e-fish/api/pkg/common/helper/logger"
	"github.com/e-fish/api/pkg/common/helper/werror"
	"github.com/e-fish/api/pkg/common/infra/orm"
	errorbudidaya "github.com/e-fish/api/pkg/domain/budidaya/error-budidaya"
	"github.com/e-fish/api/pkg/domain/budidaya/model"
	"github.com/e-fish/api/pkg/domain/pond"
	status "github.com/e-fish/api/pkg/domain/pond/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newCommand(ctx context.Context, db *gorm.DB, pondRrepo pond.Repo) Command {
	var (
		dbTxn = orm.BeginTxn(ctx, db)
	)

	return &command{
		dbTxn:     dbTxn.WithContext(ctx),
		query:     newQuery(dbTxn),
		pondQuery: pondRrepo.NewQuery(),
	}
}

type command struct {
	dbTxn     *gorm.DB
	query     Query
	pondQuery pond.Query
}

// UpdateStatusBudidayaWithListPricelist implements Command.
func (c *command) UpdateStatusBudidayaWithListPricelist(ctx context.Context, input model.UpdateBudidayaWithPricelist) (*uuid.UUID, error) {
	var (
		userID, _ = ctxutil.GetUserID(ctx)
	)

	err := input.Validate()
	if err != nil {
		return nil, err
	}

	updatedBudidaya := input.ToBudidaya(userID)

	err = c.dbTxn.Updates(&updatedBudidaya).Error
	if err != nil {
		return nil, errorbudidaya.ErrFailedUpdateBudidaya.AttacthDetail(map[string]any{"error": err})
	}

	pricelist := model.UpdatePricelistInputToPricelist(input.Pricelist, userID)
	err = c.dbTxn.Save(&pricelist).Error
	if err != nil {
		return nil, errorbudidaya.ErrFailedUpdateBudidaya.AttacthDetail(map[string]any{"error": err, "flag": "save pricelist"})
	}

	return &input.BudidayaID, nil
}

// CreateBudidaya implements Command.
func (c *command) CreateBudidaya(ctx context.Context, input model.CreateBudidayaInput) (*uuid.UUID, error) {
	var (
		userID, _ = ctxutil.GetUserID(ctx)
		pondID, _ = ctxutil.GetPondID(ctx)
	)

	data, err := c.pondQuery.GetPondByID(ctx, pondID)
	if err != nil {
		return nil, err
	}

	if data.Status != status.ACTIVED {
		return nil, errorbudidaya.ErrFailedCreateBudidaya.AttacthDetail(map[string]any{"pond-status": data.Status})
	}

	err = input.Validate()
	if err != nil {
		return nil, err
	}

	logger.InfoWithContext(ctx, "###find existing budidaya by pool id for validate budidaya not exist")
	exist, err := c.query.ReadBudidayaActiveByPoolID(ctx, input.PoolID)
	if !errorbudidaya.ErrFoundBudidaya.Is(err) {
		return nil, err
	}

	if exist != nil {
		return nil, errorbudidaya.ErrFailedCreateBudidayaExist.AttacthDetail(map[string]any{"pool": exist.PoolID})
	}

	if input.Code == "" {
		code, err := c.query.ReadBudidayaCodeActive(ctx)
		if err != nil {
			if !errorbudidaya.ErrFoundBudidayaCode.Is(err) {
				return nil, err
			}
		}

		if code == nil {
			newCode, _ := model.GeneratedCodeBudidaya(data.Name, "")
			code = &newCode
		}

		input.Code = *code

	}

	newBudidaya := input.ToBudidaya(userID, pondID)

	err = c.dbTxn.Create(&newBudidaya).Error
	if err != nil {
		return nil, errorbudidaya.ErrFailedCreateBudidaya.AttacthDetail(map[string]any{"error": err})
	}

	return &newBudidaya.ID, nil
}

// CreateFishSpecies implements Command.
func (c *command) CreateFishSpecies(ctx context.Context, input model.CreateFishSpeciesInput) (*uuid.UUID, error) {
	var (
		userID, _ = ctxutil.GetUserID(ctx)
	)

	err := input.Validate()
	if err != nil {
		return nil, err
	}

	newFishSpecies := input.ToFishSpecies(userID)

	err = c.dbTxn.Create(&newFishSpecies).Error
	if err != nil {
		return nil, errorbudidaya.ErrFailedCreateBudidaya.AttacthDetail(map[string]any{"error": err})
	}

	return &newFishSpecies.ID, nil
}

// CreateMultiplePricelistBudidaya implements Command.
func (c *command) CreateMultiplePricelistBudidaya(ctx context.Context, input model.CreateMultiplePriceListInput) ([]*uuid.UUID, error) {
	var (
		userID, _ = ctxutil.GetUserID(ctx)
		uid       = []*uuid.UUID{}
	)

	err := input.Validate()
	if err != nil {
		return nil, err
	}

	exist, err := c.query.ReadBudidayaByID(ctx, input.BudidayaID)
	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, errorbudidaya.ErrFoundBudidaya
	}

	if input.EstDate.Before(exist.DateOfSeed) {
		return nil, errorbudidaya.ErrFailedUpdateBudidayaEstDate.AttacthDetail(map[string]any{"Sowing-period": exist.DateOfSeed, "Harvest-estimate": input.EstDate})
	}

	newPricelist := input.ToMultiplePriceList(userID)

	err = c.dbTxn.Create(&newPricelist).Error
	if err != nil {
		return nil, errorbudidaya.ErrFailedCreateBudidaya.AttacthDetail(map[string]any{"error": err})
	}

	_, err = c.UpdateStatusBudidaya(ctx, model.UpdateBudidayaStatusInput{
		ID:        input.BudidayaID,
		EstTonase: input.EstTonase,
		Status:    model.PANEN,
		EstDate:   input.EstDate,
	})

	if err != nil {
		return nil, err
	}

	for _, v := range newPricelist {
		uid = append(uid, &v.ID)
	}

	return uid, nil
}

// UpdateStatusBudidaya implements Command.
func (c *command) UpdateStatusBudidaya(ctx context.Context, input model.UpdateBudidayaStatusInput) (*uuid.UUID, error) {
	var (
		userID, _ = ctxutil.GetUserID(ctx)
	)

	newStatus := input.ToBudidaya(userID)

	err := c.dbTxn.Updates(&newStatus).Error
	if err != nil {
		return nil, errorbudidaya.ErrFailedUpdateBudidaya.AttacthDetail(map[string]any{"error": err})
	}

	return &newStatus.ID, nil
}

// Commit implements Command.
func (c *command) Commit(ctx context.Context) error {
	if err := orm.CommitTxn(ctx); err != nil {
		return errorbudidaya.ErrCommit.AttacthDetail(map[string]any{"errors": err})
	}
	return nil
}

// Rollback implements Command.
func (c *command) Rollback(ctx context.Context) error {
	if err := orm.RollbackTxn(ctx); err != nil {
		return errorbudidaya.ErrRollback.AttacthDetail(map[string]any{"errors": err})
	}
	return nil
}

// UpdateBudidayaSoldQty implements Command.
func (c *command) UpdateBudidayaSoldQty(ctx context.Context, input model.UpdateBudidayaSoldQty) (*uuid.UUID, error) {
	var (
		userID, _ = ctxutil.GetUserID(ctx)
		today     = time.Now()
	)

	exist, err := c.query.lock().ReadBudidayaByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, errorbudidaya.ErrFoundBudidaya.AttacthDetail(map[string]any{"error": "budidaya empty"})
	}

	exist.Sold = exist.Sold + input.SoldQty
	if input.IsCancel {
		exist.Sold = exist.Sold - input.SoldQty

	}

	if !input.IsCancel && (input.SoldQty > exist.Stock) {
		return nil, werror.Error{
			Code:    "FailedUpdateOrder",
			Message: "Order Estimate Exceeded Capacity",
		}
	}

	err = c.dbTxn.Where("deleted_at IS NULL AND id = ?", input.ID).Updates(&model.Budidaya{
		Sold: exist.Sold,
		OrmModel: orm.OrmModel{
			UpdatedAt: &today,
			UpdatedBy: &userID,
		},
	}).Error
	if err != nil {
		return nil, errorbudidaya.ErrFailedUpdateBudidaya.AttacthDetail(map[string]any{"error": err})
	}

	return &input.ID, nil
}
