package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"crypto/sha256"

	"github.com/golang-jwt/jwt/v5"
)

// Возвращает ключ для подписи JWT
// Берём пароль из TODO_PASSWORD
// Используем запасной ключ, если пароль не задан.
func getJWTKey() []byte {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return []byte("fallback_key")
	}
	return []byte(password)
}

// Cоздаём хэш пароля для хранения в JWT-токене
func jwtHash(password string) string {
	hash := sha256.Sum256([]byte(password))
		return fmt.Sprintf("%x", hash)
}

type CustomClaims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// Создаём обработчик аутентификации
func signinHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON с паролем
	var request struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	
	password := os.Getenv("TODO_PASSWORD")

	// Сравниваем пароли
	if request.Password != password {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": "invalid password"})
		return
	}

	// Создаём JWT-токен
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := CustomClaims{
		PasswordHash: jwtHash(password),
		RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTKey())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
			return
	}

	// Возвращаем токен
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": tokenString}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// Создаём функцию для проверки аутентификации
// Переименовал переменную jwt, потому что был конфликт с пакетом
func auth(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    
    var jwtToken string
		
    // Получаем куку
    cookie, err := r.Cookie("token")
    if err == nil {
        jwtToken = cookie.Value
    } else {
    // Возвращаем ошибку если куки не было
    http.Error(w, "Authentication required", http.StatusUnauthorized)
        return
    }

    var valid bool
	
    // Здесь код для валидации и проверки JWT-токена
    claims := &CustomClaims{}
    token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
        return getJWTKey(), nil
    })

    valid = err == nil && token.Valid

    // Сравниваем хэш пароля из токена с текущим паролем
    if valid {
        currentHash := jwtHash(os.Getenv("TODO_PASSWORD"))
        if claims.PasswordHash != currentHash {
            valid = false
    }
    }

    if !valid {
    // Возвращаем ошибку авторизации 401
    http.Error(w, "Authentication required", http.StatusUnauthorized)
        return
    }
        next(w, r)
    })
}


