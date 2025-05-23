package middleware

import (
	"net/http"

	"github.com/hongyuxuan/filebrowser/common/constant"
	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type ValidateuserMiddleware struct {
}

func NewValidateuserMiddleware() *ValidateuserMiddleware {
	return &ValidateuserMiddleware{}
}

func (m *ValidateuserMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Logger := logx.WithContext(r.Context())
		username, role := utils.GetPayload(r.Context())
		if role != constant.ROLE_ADMIN { // admin has all privileges
			if r.Method != http.MethodGet && role == constant.ROLE_READONLY {
				Logger.Errorf("User=%s, role=%s is not allowed for httpmethod=%s", username, role, r.Method)
				httpx.Error(w, errorx.NewError(http.StatusForbidden, "没有此权限", nil))
				return
			}
		}
		next(w, r)
	}
}
