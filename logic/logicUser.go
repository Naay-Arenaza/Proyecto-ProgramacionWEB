package logic

import (
	db "ProyectoFinanzas/db/sqlc" // Usamos el alias db para las estructuras generadas
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("usuario no encontrado o credenciales inválidas")

const SessionDuration = 15 * time.Minute

type UserCapaLogica struct {
	repo UsuarioRepository
}

func NewUserLogic(r UsuarioRepository) *UserCapaLogica {
	return &UserCapaLogica{
		repo: r,
	}
}

func (l *UserCapaLogica) Authenticate(ctx context.Context, email, password string) (db.Usuario, error) {
	user, err := l.repo.GetUsuarioMail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Usuario{}, ErrUserNotFound
		}
		return db.Usuario{}, err
	}
	if user.Contraseña == password {
		return user, nil
	}

	return db.Usuario{}, ErrUserNotFound
}

func (l *UserCapaLogica) CreateUserSession(ctx context.Context, userID int32) (string, time.Time, error) {
	sessionToken := uuid.NewString() // Generar un token único
	expiresAt := time.Now().Add(SessionDuration)

	arg := db.CreateSessionParams{
		SessionToken: sessionToken,
		IDUsuario:    userID,
		Expires:      expiresAt,
	}

	_, err := l.repo.CreateSession(ctx, arg)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error al guardar sesión en DB: %w", err)
	}

	return sessionToken, expiresAt, nil
}

func (l *UserCapaLogica) GetUserSession(ctx context.Context, sessionToken string) (db.Session, error) {
	session, err := l.repo.GetSession(ctx, sessionToken)
	if err != nil {
		return db.Session{}, ErrUserNotFound
	}

	if session.Expires.Before(time.Now()) {
		return db.Session{}, ErrUserNotFound
	}

	return session, nil
}

func (l *UserCapaLogica) DeleteUserSession(ctx context.Context, sessionToken string) error {
	return l.repo.DeleteSession(ctx, sessionToken)
}
