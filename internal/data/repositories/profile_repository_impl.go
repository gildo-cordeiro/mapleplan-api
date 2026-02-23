package repositories

import (
"context"

"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/profile"
ports "github.com/gildo-cordeiro/mapleplan-api/internal/ports/repositories"
"gorm.io/gorm"
)

type ProfileRepositoryImpl struct {
db *gorm.DB
}

func (p *ProfileRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
if ctx == nil {
return p.db
}
if tx, ok := ctx.Value("tx_key").(*gorm.DB); ok && tx != nil {
return tx.WithContext(ctx)
}
return p.db.WithContext(ctx)
}

func NewGormProfileRepository(db *gorm.DB) ports.ProfileRepository {
return &ProfileRepositoryImpl{db: db}
}

func (p *ProfileRepositoryImpl) FindByID(ctx context.Context, id string) (*profile.ImmigrationProfile, error) {
var prof profile.ImmigrationProfile
err := p.getDB(ctx).Where("id = ?", id).First(&prof).Error
if err != nil {
return nil, err
}
return &prof, nil
}

func (p *ProfileRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*profile.ImmigrationProfile, error) {
var profiles []*profile.ImmigrationProfile
err := p.getDB(ctx).Where("user_id = ?", userID).Find(&profiles).Error
if err != nil {
return nil, err
}
return profiles, nil
}

func (p *ProfileRepositoryImpl) Save(ctx context.Context, prof *profile.ImmigrationProfile) error {
return p.getDB(ctx).Create(prof).Error
}

func (p *ProfileRepositoryImpl) Update(ctx context.Context, id string, prof *profile.ImmigrationProfile) error {
foundProf, err := p.FindByID(ctx, id)
if err != nil {
return err
}

foundProf.Name = prof.Name
foundProf.UserID = prof.UserID

return p.getDB(ctx).Save(foundProf).Error
}

func (p *ProfileRepositoryImpl) Delete(ctx context.Context, id string) error {
return p.getDB(ctx).Where("id = ?", id).Delete(&profile.ImmigrationProfile{}).Error
}

// ProfileMemberRepositoryImpl
type ProfileMemberRepositoryImpl struct {
db *gorm.DB
}

func (pm *ProfileMemberRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
if ctx == nil {
return pm.db
}
if tx, ok := ctx.Value("tx_key").(*gorm.DB); ok && tx != nil {
return tx.WithContext(ctx)
}
return pm.db.WithContext(ctx)
}

func NewGormProfileMemberRepository(db *gorm.DB) ports.ProfileMemberRepository {
return &ProfileMemberRepositoryImpl{db: db}
}

func (pm *ProfileMemberRepositoryImpl) FindByID(ctx context.Context, id string) (*profile.ProfileMember, error) {
var member profile.ProfileMember
err := pm.getDB(ctx).Where("id = ?", id).First(&member).Error
if err != nil {
return nil, err
}
return &member, nil
}

func (pm *ProfileMemberRepositoryImpl) FindByProfileID(ctx context.Context, profileID string) ([]*profile.ProfileMember, error) {
var members []*profile.ProfileMember
err := pm.getDB(ctx).Where("profile_id = ?", profileID).Find(&members).Error
if err != nil {
return nil, err
}
return members, nil
}

func (pm *ProfileMemberRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*profile.ProfileMember, error) {
var members []*profile.ProfileMember
err := pm.getDB(ctx).Where("user_id = ?", userID).Find(&members).Error
if err != nil {
return nil, err
}
return members, nil
}

func (pm *ProfileMemberRepositoryImpl) Save(ctx context.Context, m *profile.ProfileMember) error {
return pm.getDB(ctx).Create(m).Error
}

func (pm *ProfileMemberRepositoryImpl) Update(ctx context.Context, id string, m *profile.ProfileMember) error {
foundMember, err := pm.FindByID(ctx, id)
if err != nil {
return err
}

foundMember.Role = m.Role

return pm.getDB(ctx).Save(foundMember).Error
}

func (pm *ProfileMemberRepositoryImpl) Delete(ctx context.Context, id string) error {
return pm.getDB(ctx).Where("id = ?", id).Delete(&profile.ProfileMember{}).Error
}
