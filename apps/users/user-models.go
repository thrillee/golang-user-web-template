package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/thrillee/triq/internals/schemas"
	"github.com/thrillee/triq/internals/security"
	"gorm.io/gorm"
)

type User struct {
	schemas.BaseModel

	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"-"`

	AccountRef string `json:"account_ref" gorm:"uniqueIndex"`

	FullName       schemas.MyString `json:"full_name"`
	DisplayPicture schemas.MyString `json:"display_picture"`
	Username       string           `json:"username"`

	Active     bool `json:"active"`
	IsVerified bool `json:"is_verified"`

	LastLogin sql.NullTime `json:"last_login"`
}

func (u User) GetID() interface{} {
	return u.ID
}

func (u User) GetModelName() interface{} {
	return "User"
}

func (u User) QueryRepo() schemas.Repository[User] {
	return schemas.NewBaseRepository[User](&User{})
}

func (u User) Query(ctx context.Context) *gorm.DB {
	return u.QueryRepo().Query(ctx)
}

func (u *User) FindByAccountRef(ctx context.Context, accountRef string) *User {
	var user User
	u.QueryRepo().Query(ctx).Model(User{AccountRef: accountRef}).First(&user)
	return &user
}

func (u *User) UsernameExists(ctx context.Context, username string) bool {
	return u.QueryRepo().Exists(ctx, u.QueryRepo().Query(ctx).Model(User{Username: username}))
}

func (u *User) EmailExists(ctx context.Context, email string) bool {
	return u.QueryRepo().Exists(ctx, u.QueryRepo().Query(ctx).Model(User{Email: email}))
}

func (u *User) PhoneExists(ctx context.Context, phone string) bool {
	return u.QueryRepo().Exists(ctx, u.QueryRepo().Query(ctx).Model(User{Phone: phone}))
}

func (u *User) SetPassword(ctx context.Context, password string) (bool, error) {
	hash, err := security.HashPassword(password)
	if err != nil {
		return false, err
	}
	u.Password = hash
	u.QueryRepo().Edit(ctx, u.GetID(), u)
	return true, nil
}

func (u *User) toString() string {
	return fmt.Sprintf("User[ID=%v,Fullname=%s,Email=%s,Mobile=%s,AccountRef=%s]",
		u.ID, u.FullName, u.Email, u.Phone, u.AccountRef)
}
