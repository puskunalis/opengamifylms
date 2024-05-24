package opengamifylms

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ConstError string

var _ error = (*ConstError)(nil)

func (err ConstError) Error() string {
	return string(err)
}

type Claims struct {
	ID              int64
	IsAdministrator bool
	Email           string
	FullName        string
}

func getAllClaimsFromToken(h http.Header, jwtSecretKey string) (Claims, error) {
	var claims Claims

	// Extract the JWT token from the Authorization header
	tokenString := h.Get("Authorization")
	if tokenString == "" {
		return claims, errors.New("missing token")
	}

	// Remove the "Bearer " prefix from the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Provide the secret key for token verification
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return claims, err
	}

	// Extract the instructor ID from the token claims
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return claims, errors.New("invalid token")
	}

	id, ok := mapClaims["id"].(float64)
	if !ok {
		return claims, errors.New("id claim not found in token")
	}
	claims.ID = int64(id)

	isAdministrator, ok := mapClaims["is_administrator"].(bool)
	if !ok {
		return claims, errors.New("is_administrator claim not found in token")
	}
	claims.IsAdministrator = isAdministrator

	email, ok := mapClaims["email"].(string)
	if !ok {
		return claims, errors.New("email claim not found in token")
	}
	claims.Email = email

	fullName, ok := mapClaims["full_name"].(string)
	if !ok {
		return claims, errors.New("full_name claim not found in token")
	}
	claims.FullName = fullName

	return claims, nil
}

type UpdateOrderRequest struct {
	Order map[int32]int32 `json:"order"`
}

func (req UpdateOrderRequest) updateOrderRequestToSlices() ([]int32, []int32) {
	ids := make([]int32, 0, len(req.Order))
	orders := make([]int32, 0, len(req.Order))

	for id, order := range req.Order {
		ids = append(ids, id)
		orders = append(orders, order)
	}

	return ids, orders
}
