package model

import (
	"time"

	"github.com/e-fish/api/pkg/common/helper/werror"
	"github.com/e-fish/api/pkg/common/infra/orm"
	errorpond "github.com/e-fish/api/pkg/domain/pond/error-pond"
	"github.com/google/uuid"
)

type CreatePondInput struct {
	Name          string              `json:"name"`
	CountryID     uuid.UUID           `gorm:"size:256" json:"countryID"`
	ProvinceID    uuid.UUID           `gorm:"size:256" json:"provinceID"`
	CityID        uuid.UUID           `gorm:"size:256" json:"cityID"`
	DistrictID    uuid.UUID           `gorm:"size:256" json:"districtID"`
	DetailAddress string              `json:"detailAddress"`
	NoteAddress   string              `json:"noteAddress"`
	Type          string              `json:"type"`
	Latitude      float64             `json:"latitude"`
	Longitude     float64             `json:"longitude"`
	Image         string              `json:"image"`
	ListPool      []CreatePoolInput   `json:"listPool"`
	ListBerkas    []CreateBerkasInput `json:"listBerkas"`
}

func (c *CreatePondInput) Validate() error {
	errs := werror.NewError("error validate input")

	if c.Name == "" {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"Name": "empty"}))
	}
	if c.CountryID == uuid.Nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"CountryID": "empty"}))
	}
	if c.ProvinceID == uuid.Nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"ProvinceID": "empty"}))
	}
	if c.CityID == uuid.Nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"CityID": "empty"}))
	}
	if c.DistrictID == uuid.Nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"DistrictID": "empty"}))
	}
	if c.Type == "" {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"Type": "empty"}))
	}

	// if c.Type == TEAM {
	// if c.TeamID == nil {
	// 	errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"Team": "empty"}))
	// }
	if err := ValidateCreateberkasInput(c.ListBerkas); err != nil {
		errs.Add(err)
	}
	// }

	if len(c.ListPool) < 1 {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"Pool": "empty"}))
	}

	err := ValidateCreatePoolInput(c.ListPool)
	if err != nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"err": err}))
	}

	return errs.Return()
}

func (c *CreatePondInput) ToPond(userID, pondID uuid.UUID) Pond {

	return Pond{
		ID:            pondID,
		UserID:        userID,
		Name:          c.Name,
		CountryID:     c.CountryID,
		ProvinceID:    c.ProvinceID,
		CityID:        c.CityID,
		DistrictID:    c.DistrictID,
		DetailAddress: c.DetailAddress,
		NoteAddress:   c.NoteAddress,
		Type:          c.Type,
		Latitude:      c.Latitude,
		Longitude:     c.Longitude,
		Status:        SUBMISION,
		Image:         c.Image,
		ListPool:      ListPoolInputToListPool(userID, pondID, c.ListPool),
		ListBerkas:    ListBerkasInputToListBerkas(userID, pondID, c.ListBerkas),
		OrmModel: orm.OrmModel{
			CreatedAt: time.Time{},
			CreatedBy: userID,
		},
	}
}

type CreateBerkasInput struct {
	Name string `json:"name"`
	File string `json:"file"`
}

func (c *CreateBerkasInput) Validate() error {
	errs := werror.NewError("error validate input")

	if c.Name == "" {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Name": "empty"}))
	}
	if c.File == "" {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"File": "empty"}))
	}

	return errs.Return()
}

func ValidateCreateberkasInput(input []CreateBerkasInput) error {
	errs := werror.NewError("error validate input")

	if len(input) < 1 {
		return errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"berkas": "empty"})
	}

	for _, v := range input {
		if err := v.Validate(); err != nil {
			errs.Add(err)
		}
	}

	return errs.Return()
}

func (c *CreateBerkasInput) ToBerkas(userID uuid.UUID, pondID uuid.UUID) Berkas {
	return Berkas{
		ID:     uuid.New(),
		PondID: pondID,
		Name:   c.Name,
		File:   c.File,
		OrmModel: orm.OrmModel{
			CreatedAt: time.Now(),
			CreatedBy: userID,
		},
	}
}

func ListBerkasInputToListBerkas(userID, pondID uuid.UUID, input []CreateBerkasInput) (newBerkas []Berkas) {
	for _, v := range input {
		newBerkas = append(newBerkas, v.ToBerkas(userID, pondID))
	}
	return
}

type CreatePoolInput struct {
	Name  string  `json:"name"`
	Long  float64 `json:"long"`
	Wide  float64 `json:"wide"`
	Image string  `json:"image"`
}

func (c *CreatePoolInput) Validate() error {
	errs := werror.NewError("error validate input")

	if c.Name == "" {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Name": "empty"}))
	}
	if c.Long == 0 {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Long": "empty"}))
	}
	if c.Wide == 0 {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Wide": "empty"}))
	}
	if c.Image == "" {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Image": "empty"}))
	}

	return errs.Return()
}

func ValidateCreatePoolInput(input []CreatePoolInput) error {
	errs := werror.NewError("error validate input")

	for _, v := range input {
		err := v.Validate()
		if err != nil {
			errs.Add(err)
		}
	}

	return errs.Return()
}

func (c *CreatePoolInput) ToPool(userID uuid.UUID, pondID uuid.UUID) Pool {
	return Pool{
		ID:       uuid.New(),
		PondID:   pondID,
		Name:     c.Name,
		Long:     c.Long,
		Wide:     c.Wide,
		Image:    c.Image,
		OrmModel: orm.OrmModel{CreatedAt: time.Now(), CreatedBy: userID},
	}
}

func ListPoolInputToListPool(userID, pondID uuid.UUID, input []CreatePoolInput) (newPool []Pool) {
	for _, v := range input {
		newPool = append(newPool, v.ToPool(userID, pondID))
	}
	return
}

type UpdatePondInput struct {
	Name          string    `json:"name"`
	CountryID     uuid.UUID `gorm:"size:256" json:"countryID"`
	ProvinceID    uuid.UUID `gorm:"size:256" json:"provinceID"`
	CityID        uuid.UUID `gorm:"size:256" json:"cityID"`
	DistrictID    uuid.UUID `gorm:"size:256" json:"districtID"`
	DetailAddress string    `json:"detailAddress"`
	NoteAddress   string    `json:"noteAddress"`
	Type          string    `json:"type"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Image         string    `json:"image"`
}

func (u *UpdatePondInput) ToPond(userID, pondID uuid.UUID) Pond {
	var (
		today = time.Now()
	)

	return Pond{
		ID:            pondID,
		Name:          u.Name,
		CountryID:     u.CountryID,
		ProvinceID:    u.ProvinceID,
		CityID:        u.CityID,
		DistrictID:    u.DistrictID,
		DetailAddress: u.DetailAddress,
		NoteAddress:   u.NoteAddress,
		Type:          u.Type,
		Latitude:      u.Latitude,
		Longitude:     u.Longitude,
		Image:         u.Image,
		OrmModel: orm.OrmModel{
			UpdatedAt: &today,
			UpdatedBy: &userID,
		},
	}
}

type UpdatePondStatus struct {
	PondID  uuid.UUID `json:"pondID"`
	Status  string    `json:"status"`
	Reasons string    `json:"reasons"`
}

func (u *UpdatePondStatus) Validate() error {
	errs := werror.NewError("failed validate input update pond status")

	if u.PondID == uuid.Nil {
		errs.Add(errorpond.ErrValidateInputbUpdateStatus.AttacthDetail(map[string]any{"PondID": "empty"}))
	}
	if u.Status == "" {
		errs.Add(errorpond.ErrValidateInputbUpdateStatus.AttacthDetail(map[string]any{"Status": "empty"}))
	}

	if u.Status == DISABLED && u.Reasons == "" {
		errs.Add(errorpond.ErrValidateInputbUpdateStatus.AttacthDetail(map[string]any{"reasons": "empty"}))
	}

	return errs.Return()
}

func (u *UpdatePondStatus) ToPond(userID uuid.UUID) Pond {
	var (
		today = time.Now()
	)

	return Pond{
		ID:      u.PondID,
		Status:  u.Status,
		Reasons: u.Reasons,
		OrmModel: orm.OrmModel{
			UpdatedAt: &today,
			UpdatedBy: &userID,
		},
	}
}

type Resubmission struct {
	Name          string              `json:"name"`
	CountryID     uuid.UUID           `gorm:"size:256" json:"countryID"`
	ProvinceID    uuid.UUID           `gorm:"size:256" json:"provinceID"`
	CityID        uuid.UUID           `gorm:"size:256" json:"cityID"`
	DistrictID    uuid.UUID           `gorm:"size:256" json:"districtID"`
	DetailAddress string              `json:"detailAddress"`
	NoteAddress   string              `json:"noteAddress"`
	Type          string              `json:"type"`
	Latitude      float64             `json:"latitude"`
	Longitude     float64             `json:"longitude"`
	Image         string              `json:"image"`
	ListPool      []UpdatePoolInput   `json:"listPool"`
	ListBerkas    []UpdateBerkasInput `json:"listBerkas"`
}

func (c *Resubmission) Validate() error {
	errs := werror.NewError("error validate input")

	if c.Name == "" {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"Name": "empty"}))
	}
	if c.CountryID == uuid.Nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"CountryID": "empty"}))
	}
	if c.ProvinceID == uuid.Nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"ProvinceID": "empty"}))
	}
	if c.CityID == uuid.Nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"CityID": "empty"}))
	}
	if c.DistrictID == uuid.Nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"DistrictID": "empty"}))
	}
	if c.Type == "" {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"Type": "empty"}))
	}

	if c.Type == TEAM {
		if err := ValidateUpdateberkasInput(c.ListBerkas); err != nil {
			errs.Add(err)
		}
	}

	if len(c.ListPool) < 1 {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"Pool": "empty"}))
	}

	err := ValidateUpdatePoolInput(c.ListPool)
	if err != nil {
		errs.Add(errorpond.ErrValidateInputPond.AttacthDetail(map[string]any{"err": err}))
	}

	return errs.Return()
}

type UpdatePoolInput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Long      float64   `json:"long"`
	Wide      float64   `json:"wide"`
	Image     string    `json:"image"`
	IsDeleted bool      `json:"isDeleted"`
}

type UpdateBerkasInput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	File      string    `json:"file"`
	IsDeleted bool      `json:"isDeleted"`
}

func (c *UpdateBerkasInput) Validate() error {
	errs := werror.NewError("error validate input")

	if c.Name == "" {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Name": "empty"}))
	}
	if c.File == "" {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"File": "empty"}))
	}

	return errs.Return()
}

func ValidateUpdateberkasInput(input []UpdateBerkasInput) error {
	errs := werror.NewError("error validate input")

	for idx := range input {
		if err := input[idx].Validate(); err != nil {
			errs.Add(err)
		}
	}

	return errs.Return()
}

func (c *UpdateBerkasInput) ToBerkas(userID uuid.UUID, pondID uuid.UUID) Berkas {
	var (
		now    = time.Now()
		berkas = Berkas{
			ID:     c.ID,
			PondID: pondID,
			Name:   c.Name,
			File:   c.File,
			OrmModel: orm.OrmModel{
				UpdatedAt: &now,
				UpdatedBy: &userID,
			},
		}
	)

	if c.IsDeleted {
		berkas.DeletedAt = &now
		berkas.DeletedBy = &userID
	}
	return berkas
}

func UpdateListBerkasInputToListBerkas(userID, pondID uuid.UUID, input []UpdateBerkasInput) (newBerkas []Berkas) {
	for idx := range input {
		if input[idx].ID == uuid.Nil {
			input[idx].ID = uuid.New()
		}
		newBerkas = append(newBerkas, input[idx].ToBerkas(userID, pondID))
	}
	return
}

func (c *UpdatePoolInput) Validate() error {
	errs := werror.NewError("error validate input")

	if c.Name == "" {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Name": "empty"}))
	}
	if c.Long == 0 {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Long": "empty"}))
	}
	if c.Wide == 0 {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Wide": "empty"}))
	}
	if c.Image == "" {
		errs.Add(errorpond.ErrValidateInputbBerkas.AttacthDetail(map[string]any{"Image": "empty"}))
	}

	return errs.Return()
}

func ValidateUpdatePoolInput(input []UpdatePoolInput) error {
	errs := werror.NewError("error validate input")

	for _, v := range input {
		err := v.Validate()
		if err != nil {
			errs.Add(err)
		}
	}

	return errs.Return()
}

func (c *UpdatePoolInput) ToPool(userID uuid.UUID, pondID uuid.UUID) Pool {
	var (
		now = time.Now()

		pool = Pool{
			ID:     c.ID,
			PondID: pondID,
			Name:   c.Name,
			Long:   c.Long,
			Wide:   c.Wide,
			Image:  c.Image,
			OrmModel: orm.OrmModel{
				UpdatedAt: &now,
				UpdatedBy: &userID,
			},
		}
	)

	if c.IsDeleted {
		pool.DeletedAt = &now
		pool.DeletedBy = &userID
	}
	return pool
}

func UpdateListPoolInputToListPool(userID, pondID uuid.UUID, input []UpdatePoolInput) (newPool []Pool) {
	for idx := range input {
		if input[idx].ID == uuid.Nil {
			input[idx].ID = uuid.New()
		}
		newPool = append(newPool, input[idx].ToPool(userID, pondID))
	}
	return
}

func (c *Resubmission) ToPond(userID, pondID uuid.UUID) Pond {
	now := time.Now()
	return Pond{
		ID:            pondID,
		UserID:        userID,
		Name:          c.Name,
		CountryID:     c.CountryID,
		ProvinceID:    c.ProvinceID,
		CityID:        c.CityID,
		DistrictID:    c.DistrictID,
		DetailAddress: c.DetailAddress,
		NoteAddress:   c.NoteAddress,
		Type:          c.Type,
		Latitude:      c.Latitude,
		Longitude:     c.Longitude,
		Status:        SUBMISION,
		Image:         c.Image,
		OrmModel: orm.OrmModel{
			UpdatedAt: &now,
			UpdatedBy: &userID,
		},
	}
}
