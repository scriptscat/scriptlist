package script_svc

import (
	"context"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"

	api "github.com/scriptscat/scriptlist/internal/api/script"
)

type AccessInviteSvc interface {
	// InviteCodeList 邀请码列表
	InviteCodeList(ctx context.Context, req *api.InviteCodeListRequest) (*api.InviteCodeListResponse, error)
	// CreateInviteLink 创建邀请链接
	CreateInviteLink(ctx context.Context, entity *script_entity.ScriptInvite) (*script_entity.ScriptInvite, error)
	// CreateInviteCode 创建邀请码
	CreateInviteCode(ctx context.Context, req *api.CreateInviteCodeRequest) (*api.CreateInviteCodeResponse, error)
	// DeleteInviteCode 删除邀请码
	DeleteInviteCode(ctx context.Context, req *api.DeleteInviteCodeRequest) (*api.DeleteInviteCodeResponse, error)
	// AuditInviteCode 审核邀请码
	AuditInviteCode(ctx context.Context, req *api.AuditInviteCodeRequest) (*api.AuditInviteCodeResponse, error)
	// AcceptInvite 接受邀请
	AcceptInvite(ctx context.Context, req *api.AcceptInviteRequest) (*api.AcceptInviteResponse, error)
	// GroupInviteCodeList 群组邀请码列表
	GroupInviteCodeList(ctx context.Context, req *api.GroupInviteCodeListRequest) (*api.GroupInviteCodeListResponse, error)
	// CreateGroupInviteCode 创建群组邀请码
	CreateGroupInviteCode(ctx context.Context, req *api.CreateGroupInviteCodeRequest) (*api.CreateGroupInviteCodeResponse, error)
	// InviteCodeInfo 邀请码信息
	InviteCodeInfo(ctx context.Context, req *api.InviteCodeInfoRequest) (*api.InviteCodeInfoResponse, error)
}

type accessInviteSvc struct {
}

var defaultAccessInvite = &accessInviteSvc{}

func AccessInvite() AccessInviteSvc {
	return defaultAccessInvite
}

// InviteCodeList 邀请码列表
func (a *accessInviteSvc) InviteCodeList(ctx context.Context, req *api.InviteCodeListRequest) (*api.InviteCodeListResponse, error) {
	script := Script().CtxScript(ctx)
	list, total, err := script_repo.ScriptInvite().FindAccessPage(ctx, script.ID, req.PageRequest)
	if err != nil {
		return nil, err
	}
	resp := &api.InviteCodeListResponse{
		PageResponse: httputils.PageResponse[*api.InviteCode]{
			List:  make([]*api.InviteCode, 0),
			Total: total,
		},
	}
	for _, v := range list {
		v, err := a.toInviteCode(ctx, v)
		if err != nil {
			return nil, err
		}
		if v != nil {
			resp.List = append(resp.List, v)
		}
	}
	return resp, nil
}

func (a *accessInviteSvc) toInviteCode(ctx context.Context, invite *script_entity.ScriptInvite) (*api.InviteCode, error) {
	ret := &api.InviteCode{
		ID:           invite.ID,
		Code:         invite.Code,
		UserID:       0,
		Username:     "",
		IsAudit:      invite.IsAudit == consts.YES,
		InviteStatus: invite.InviteStatus,
		Expiretime:   invite.Expiretime,
		Createtime:   invite.Createtime,
	}
	if invite.UserID > 0 {
		user, err := user_repo.User().Find(ctx, invite.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, nil
		}
		ret.UserID = user.ID
		ret.Username = user.Username
	}
	if ret.InviteStatus == script_entity.InviteStatusUnused && ret.Expiretime > 0 && ret.Expiretime < time.Now().Unix() {
		ret.InviteStatus = script_entity.InviteStatusExpired
	}

	return ret, nil
}

// CreateInviteLink 创建邀请链接
func (a *accessInviteSvc) CreateInviteLink(ctx context.Context, entity *script_entity.ScriptInvite) (*script_entity.ScriptInvite, error) {
	invite := &script_entity.ScriptInvite{
		ScriptID:     entity.ScriptID,
		Code:         utils.RandString(32, utils.Letter),
		CodeType:     script_entity.InviteCodeTypeLink,
		GroupID:      entity.GroupID,
		Type:         entity.Type,
		UserID:       entity.UserID,
		IsAudit:      consts.NO,
		InviteStatus: script_entity.InviteStatusUnused,
		Status:       consts.ACTIVE,
		Createtime:   time.Now().Unix(),
		Updatetime:   time.Now().Unix(),
	}
	if err := script_repo.ScriptInvite().Create(ctx, invite); err != nil {
		return nil, err
	}
	return invite, nil
}

func (a *accessInviteSvc) createInviteCode(ctx context.Context, count int32, entity *script_entity.ScriptInvite) ([]string, error) {
	if count > 100 {
		count = 100
	}
	codes := make([]string, 0)
	if err := db.Ctx(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = db.WithContextDB(ctx, tx)
		for i := 0; i < int(count); i++ {
			invite := &script_entity.ScriptInvite{
				ScriptID:     entity.ScriptID,
				Code:         utils.RandString(16, utils.Letter),
				CodeType:     script_entity.InviteCodeTypeCode,
				GroupID:      entity.GroupID,
				Type:         entity.Type,
				UserID:       0,
				IsAudit:      entity.IsAudit,
				InviteStatus: script_entity.InviteStatusUnused,
				Status:       consts.ACTIVE,
				Expiretime:   entity.Expiretime,
				Createtime:   time.Now().Unix(),
				Updatetime:   time.Now().Unix(),
			}
			if err := script_repo.ScriptInvite().Create(ctx, invite); err != nil {
				return err
			}
			codes = append(codes, invite.Code)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return codes, nil
}

// CreateInviteCode 创建邀请码
func (a *accessInviteSvc) CreateInviteCode(ctx context.Context, req *api.CreateInviteCodeRequest) (*api.CreateInviteCodeResponse, error) {
	script := Script().CtxScript(ctx)
	isAudit := int32(consts.NO)
	if req.Audit {
		isAudit = consts.YES
	}
	entity := &script_entity.ScriptInvite{
		ScriptID: script.ID,
		IsAudit:  isAudit,
		Type:     script_entity.InviteTypeAccess,
	}
	if req.Days != 0 {
		entity.Expiretime = time.Now().Add(time.Duration(req.Days) * 24 * time.Hour).Unix()
	}
	codes, err := a.createInviteCode(ctx, req.Count, entity)
	if err != nil {
		return nil, err
	}
	return &api.CreateInviteCodeResponse{
		Code: codes,
	}, nil
}

// DeleteInviteCode 删除邀请码
func (a *accessInviteSvc) DeleteInviteCode(ctx context.Context, req *api.DeleteInviteCodeRequest) (*api.DeleteInviteCodeResponse, error) {
	script := Script().CtxScript(ctx)
	if err := script_repo.ScriptInvite().Delete(ctx, script.ID, req.CodeID); err != nil {
		return nil, err
	}
	return &api.DeleteInviteCodeResponse{}, nil
}

// AuditInviteCode 审核邀请码
func (a *accessInviteSvc) AuditInviteCode(ctx context.Context, req *api.AuditInviteCodeRequest) (*api.AuditInviteCodeResponse, error) {
	script := Script().CtxScript(ctx)
	invite, err := script_repo.ScriptInvite().Find(ctx, script.ID, req.CodeID)
	if err != nil {
		return nil, err
	}
	if invite == nil || invite.CodeType != script_entity.InviteCodeTypeCode {
		return nil, i18n.NewNotFoundError(ctx, code.AccessInviteNotFound)
	}
	if invite.IsAudit != consts.YES {
		return nil, i18n.NewNotFoundError(ctx, code.AccessInviteNotAudit)
	}
	err = db.Ctx(ctx).Transaction(func(tx *gorm.DB) error {
		// 加入access
		ctx = db.WithContextDB(ctx, tx)
		switch invite.InviteStatus {
		case script_entity.InviteStatusPending:
			if req.Status == 1 {
				invite.InviteStatus = script_entity.InviteStatusUsed
				switch invite.Type {
				case script_entity.InviteTypeAccess:
					// 加入access
					if err := Access().AddAccess(ctx, &script_entity.ScriptAccess{
						ScriptID:     script.ID,
						LinkID:       invite.UserID,
						Type:         script_entity.AccessTypeUser,
						Role:         script_entity.AccessRoleGuest,
						InviteStatus: script_entity.AccessInviteStatusAccept,
						Status:       consts.ACTIVE,
						Expiretime:   0,
					}); err != nil {
						return err
					}
				case script_entity.InviteTypeGroup:
					// 加入群组
					if err := Group().AddMemberInternal(ctx, &script_entity.ScriptGroupMember{
						ScriptID:     invite.ScriptID,
						GroupID:      invite.GroupID,
						UserID:       invite.UserID,
						InviteStatus: script_entity.AccessInviteStatusAccept,
						Status:       consts.ACTIVE,
						Expiretime:   0,
					}); err != nil {
						return err
					}
				}
			} else {
				invite.InviteStatus = script_entity.InviteStatusReject
			}
		default:
			return i18n.NewNotFoundError(ctx, code.AccessInviteNotPending)
		}
		if err := script_repo.ScriptInvite().Update(ctx, invite); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &api.AuditInviteCodeResponse{}, nil
}

// AcceptInvite 接受邀请
func (a *accessInviteSvc) AcceptInvite(ctx context.Context, req *api.AcceptInviteRequest) (*api.AcceptInviteResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	invite, err := script_repo.ScriptInvite().FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if err := invite.Check(ctx); err != nil {
		return nil, err
	}
	if invite.CodeType == script_entity.InviteCodeTypeLink {
		// 邀请链接
		err = db.Ctx(ctx).Transaction(func(tx *gorm.DB) error {
			ctx = db.WithContextDB(ctx, tx)
			switch invite.Type {
			case script_entity.InviteTypeAccess:
				// 加入access
				// 搜索access记录
				access, err := script_repo.ScriptAccess().Find(ctx, invite.ScriptID, invite.UserID)
				if err != nil {
					return err
				}
				if err := access.Check(ctx); err != nil {
					return err
				}
				if user.UID != access.LinkID {
					return i18n.NewNotFoundError(ctx, code.AccessInviteUserError)
				}
				// 修改状态
				if req.Accept {
					access.InviteStatus = script_entity.AccessInviteStatusAccept
				} else {
					access.InviteStatus = script_entity.AccessInviteStatusReject
				}
				if err := Access().AddAccess(ctx, access); err != nil {
					logger.Ctx(ctx).Error("邀请链接邀请失败",
						zap.Int64("script_id", invite.ScriptID), zap.Int64("access_id", invite.UserID),
						zap.Int64("user_id", user.UID), zap.String("code", req.Code), zap.Error(err),
						zap.String("code", req.Code), zap.Error(err))
					return err
				}
			case script_entity.InviteTypeGroup:
				// 加入群组
				// 搜索群组记录
				member, err := script_repo.ScriptGroupMember().Find(ctx, invite.ScriptID, invite.UserID)
				if err != nil {
					return err
				}
				if err := member.Check(ctx); err != nil {
					return err
				}
				if user.UID != member.UserID {
					return i18n.NewNotFoundError(ctx, code.AccessInviteUserError)
				}
				// 修改状态
				if req.Accept {
					member.InviteStatus = script_entity.AccessInviteStatusAccept
				} else {
					member.InviteStatus = script_entity.AccessInviteStatusReject
				}
				if err := Group().AddMemberInternal(ctx, member); err != nil {
					return err
				}
			default:
				return i18n.NewNotFoundError(ctx, code.AccessInviteInvalid)
			}
			// 更新邀请码状态
			invite.InviteStatus = script_entity.InviteStatusUsed
			invite.UserID = user.UID
			invite.Updatetime = time.Now().Unix()
			if err := script_repo.ScriptInvite().Update(ctx, invite); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return &api.AcceptInviteResponse{}, nil
	}
	// 邀请码
	if !req.Accept {
		// 邀请码不能拒绝
		return nil, i18n.NewNotFoundError(ctx, code.AccessInviteInvalid)
	}
	err = db.Ctx(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = db.WithContextDB(ctx, tx)
		// 判断是否还要审核
		if invite.IsAudit == consts.YES {
			// 修改状态为等待审核
			invite.UserID = user.UID
			invite.InviteStatus = script_entity.InviteStatusPending
			if err := script_repo.ScriptInvite().Update(ctx, invite); err != nil {
				return err
			}
			return nil
		}
		switch invite.Type {
		case script_entity.InviteTypeAccess:
			// 加入access
			if err := Access().AddAccess(ctx, &script_entity.ScriptAccess{
				ScriptID:     invite.ScriptID,
				LinkID:       user.UID,
				Type:         script_entity.AccessTypeUser,
				Role:         script_entity.AccessRoleGuest,
				InviteStatus: script_entity.AccessInviteStatusAccept,
				Status:       consts.ACTIVE,
				Expiretime:   0,
			}); err != nil {
				return err
			}
		case script_entity.InviteTypeGroup:
			if err := Group().AddMemberInternal(ctx, &script_entity.ScriptGroupMember{
				ScriptID:     invite.ScriptID,
				GroupID:      invite.GroupID,
				UserID:       user.UID,
				InviteStatus: script_entity.AccessInviteStatusAccept,
				Status:       consts.ACTIVE,
				Expiretime:   0,
			}); err != nil {
				return err
			}
		default:
			return i18n.NewNotFoundError(ctx, code.AccessInviteInvalid)
		}
		// 更新邀请码状态
		invite.InviteStatus = script_entity.InviteStatusUsed
		invite.UserID = user.UID
		invite.Updatetime = time.Now().Unix()
		if err := script_repo.ScriptInvite().Update(ctx, invite); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &api.AcceptInviteResponse{}, nil
}

// GroupInviteCodeList 群组邀请码列表
func (a *accessInviteSvc) GroupInviteCodeList(ctx context.Context, req *api.GroupInviteCodeListRequest) (*api.GroupInviteCodeListResponse, error) {
	script := Script().CtxScript(ctx)
	list, total, err := script_repo.ScriptInvite().FindGroupPage(ctx, script.ID, req.GroupID, req.PageRequest)
	if err != nil {
		return nil, err
	}
	resp := &api.GroupInviteCodeListResponse{
		PageResponse: httputils.PageResponse[*api.InviteCode]{
			List:  make([]*api.InviteCode, 0),
			Total: total,
		},
	}
	for _, v := range list {
		v, err := a.toInviteCode(ctx, v)
		if err != nil {
			return nil, err
		}
		if v != nil {
			resp.List = append(resp.List, v)
		}
	}
	return resp, nil
}

// CreateGroupInviteCode 创建群组邀请码
func (a *accessInviteSvc) CreateGroupInviteCode(ctx context.Context, req *api.CreateGroupInviteCodeRequest) (*api.CreateGroupInviteCodeResponse, error) {
	script := Script().CtxScript(ctx)
	isAudit := int32(consts.NO)
	if req.Audit {
		isAudit = consts.YES
	}
	entity := &script_entity.ScriptInvite{
		ScriptID: script.ID,
		GroupID:  req.GroupID,
		IsAudit:  isAudit,
		Type:     script_entity.InviteTypeGroup,
	}
	if req.Days != 0 {
		entity.Expiretime = time.Now().Add(time.Duration(req.Days) * 24 * time.Hour).Unix()
	}
	codes, err := a.createInviteCode(ctx, req.Count, entity)
	if err != nil {
		return nil, err
	}
	return &api.CreateGroupInviteCodeResponse{
		Code: codes,
	}, nil
}

// InviteCodeInfo 邀请码信息
func (a *accessInviteSvc) InviteCodeInfo(ctx context.Context, req *api.InviteCodeInfoRequest) (*api.InviteCodeInfoResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	invite, err := script_repo.ScriptInvite().FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, i18n.NewNotFoundError(ctx, code.AccessInviteNotFound)
	}
	script, err := script_repo.Script().Find(ctx, invite.ScriptID)
	if err != nil {
		return nil, err
	}
	if err := script.CheckOperate(ctx); err != nil {
		return nil, err
	}
	scriptInfo, err := Script().ToScript(ctx, script, false, "")
	if err != nil {
		return nil, err
	}
	resp := &api.InviteCodeInfoResponse{
		CodeType:     invite.CodeType,
		InviteStatus: invite.GetInviteStatus(),
		Type:         invite.Type,
		IsAudit:      invite.IsAudit == consts.YES,
		Script:       scriptInfo,
	}
	if invite.Type == script_entity.InviteTypeAccess {
		if invite.CodeType == script_entity.InviteCodeTypeLink {
			access, err := script_repo.ScriptAccess().Find(ctx, invite.ScriptID, invite.UserID)
			if err != nil {
				return nil, err
			}
			if err := access.Check(ctx); err != nil {
				return nil, err
			}
			if user.UID != access.LinkID {
				return nil, i18n.NewNotFoundError(ctx, code.AccessInviteUserError)
			}
			resp.Access = &api.InviteCodeInfoAccess{
				Role: access.Role,
			}
		} else {
			resp.Access = &api.InviteCodeInfoAccess{
				Role: script_entity.AccessRoleGuest,
			}
		}
	} else {
		group, err := script_repo.ScriptGroup().Find(ctx, invite.ScriptID, invite.GroupID)
		if err != nil {
			return nil, err
		}
		if err := group.Check(ctx); err != nil {
			return nil, err
		}
		resp.Group = &api.InviteCodeInfoGroup{
			Name:        group.Name,
			Description: group.Description,
		}
		// 链接的话还要检查邀请用户与当前用户是否一致
		if invite.CodeType == script_entity.InviteCodeTypeLink {
			member, err := script_repo.ScriptGroupMember().Find(ctx, invite.ScriptID, invite.UserID)
			if err != nil {
				return nil, err
			}
			if err := member.Check(ctx); err != nil {
				return nil, err
			}
			if user.UID != member.UserID {
				return nil, i18n.NewNotFoundError(ctx, code.AccessInviteUserError)
			}
		}
	}
	return resp, nil
}
