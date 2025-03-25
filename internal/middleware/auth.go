package middleware

import (
	"context"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

//---------------------------------------------------------------------------------------
//                          МИДДЛВЕР ДЛЯ ЗАЩИЩЕННЫХ МАРШРУТОВ
//---------------------------------------------------------------------------------------

// AuthMiddleware - это функция, возвращающая httprouter.Handle.
// Она принимает "next" - конечный обработчик, который будет вызван,
// только если в middleware проверка токена прошла успешно.
//
// Благодаря этому мы можем оборачивать любые маршруты,
// и они автоматически становятся защищёнными, требующими валидный access-токен.

func Auth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cfg := config.GetConfig()

		// 1. Пытаемся извлечь куку "access_token"
		accessCookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Необходима авторизация (нет access_token)", http.StatusUnauthorized)
			return
		}

		// 2. Валидируем access-токен
		claims, err := service.ValidateAccessToken(cfg, accessCookie.Value)
		if err != nil {
			http.Error(w, "Невалидный или просроченный access-токен", http.StatusUnauthorized)
			return
		}

		// 3. Если токен валидный, можем сохранить user_id в контексте request,
		//    чтобы передать информацию дальше в защищённый обработчик.
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		r = r.WithContext(ctx)

		// 4. Вызываем "next" (защищённый маршрут), передавая ему обновлённый request с контекстом.
		next(w, r, ps)
	}
}
