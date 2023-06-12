package gormimpl

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/db/table"
	"gorm.io/gorm"
)

func NewApplet(db *gorm.DB) *Applet {
	return &Applet{
		db: db,
	}
}

type Applet struct {
	db *gorm.DB
}

func (o *Applet) Create(ctx context.Context, applets ...*table.Applet) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&applets).Error)
}

func (o *Applet) Del(ctx context.Context, ids []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("id in (?)", ids).Delete(&table.Applet{}).Error)
}

func (o *Applet) Update(ctx context.Context, id string, data map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&table.Applet{}).Where("id = ?", id).Updates(data).Error)
}

func (o *Applet) Take(ctx context.Context, id string) (*table.Applet, error) {
	var a table.Applet
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("id = ?", id).Take(&a).Error)
}

func (o *Applet) Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*table.Applet, error) {
	return ormutil.GormSearch[table.Applet](o.db.WithContext(ctx), []string{"name", "id", "app_id", "version"}, keyword, page, size)
}

func (o *Applet) FindOnShelf(ctx context.Context) ([]*table.Applet, error) {
	var ms []*table.Applet
	return ms, errs.Wrap(o.sort(o.db).Where("status = ?", constant.StatusOnShelf).Find(&ms).Error)
}

func (o *Applet) FindID(ctx context.Context, ids []string) ([]*table.Applet, error) {
	var ms []*table.Applet
	return ms, errs.Wrap(o.sort(o.db).Where("id in (?)", ids).Find(&ms).Error)
}

func (o *Applet) sort(db *gorm.DB) *gorm.DB {
	return db.Order("priority desc, create_time desc")
}