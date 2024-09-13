package users

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/thrillee/triq/internals/common"
	"github.com/thrillee/triq/internals/media"
	"github.com/thrillee/triq/internals/schemas"
	"github.com/thrillee/triq/internals/security"
)

type UserServiceImp struct {
	limit        int
	searchFields []string
	repo         schemas.Repository[User]
	model        User
}

func NewUserService() UserServiceImp {
	user := User{}
	return UserServiceImp{
		limit:        100,
		searchFields: []string{"Username"},
		repo:         user.QueryRepo(),
		model:        user,
	}
}

func (i UserServiceImp) ChangePassword(ctx context.Context, u *User, d *ChangePasswordPayload) error {
	if err := common.Validate(d); err != nil {
		return err
	}
	if !security.CheckPasswordHash(d.OldPassword, u.Password) {
		return common.NewSessionError(common.BAD_REQUEST, fmt.Errorf("Invalid Password"))
	}

	u.SetPassword(ctx, d.NewPassword)

	go publishEvent(ctx, USER_CHANGED_PASSWORD, u)
	return nil
}

func (i UserServiceImp) UploadDisplayPicutre(ctx context.Context, user *User, fh *multipart.FileHeader) error {
	fileURL, err := media.HandleMediaUpload("accounts/display-picture/", fh)
	if err != nil {
		return common.NewSessionError(common.BAD_REQUEST, err)
	}

	user.DisplayPicture = schemas.MyString(fileURL)
	_, err = user.QueryRepo().Save(ctx, user)
	if err != nil {
		return common.NewSessionError(common.INTERNVAL_SERVER_ERROR, err)
	}
	return nil
}

func (i UserServiceImp) EditUser(ctx context.Context, accountRef string, d *NewUserPayload) (*User, error) {
	if err := common.Validate(d); err != nil {
		return nil, err
	}

	var user User
	result := i.repo.Query(ctx).Model(User{AccountRef: accountRef}).First(&user)
	err := schemas.HandleDBError(result, "Account")
	if err != nil {
		return nil, common.NewSessionError(common.BAD_REQUEST, err)
	}

	if user.Username != d.Username && i.model.UsernameExists(ctx, d.Username) {
		return nil, common.NewSessionError(common.BAD_REQUEST, fmt.Errorf("Username %s already exists", d.Username))
	}

	user.FullName = d.Fullname
	user.Username = d.Username
	user.Phone = d.Phone

	usr, err := i.repo.Save(ctx, &user)
	if err != nil {
		return nil, common.NewSessionError(common.INTERNVAL_SERVER_ERROR, err)
	}
	return usr, nil
}
