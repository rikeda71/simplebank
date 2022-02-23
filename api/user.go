package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/s14t284/simplebank/db/ent"
	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	UserName          string    `json:"username,omitempty"`
	FullName          string    `json:"full_name,omitempty"`
	Email             string    `json:"email,omitempty"`
	PasswordChangedAt time.Time `json:"password_changed_at,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})
	if err != nil {
		// ent を使っているのでこっちに引っかかる
		if ent.IsConstraintError(err) || ent.IsValidationError(err) {
			if strings.Contains(err.Error(), "violate") {
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		// ORMに頼っていない場合は以下に引っかかる
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		UserName:          user.ID,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
