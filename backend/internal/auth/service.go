package auth

import (
	"context"

	"cinema-booking-backend/internal/model"
	"cinema-booking-backend/internal/repository"
)

type Service struct {
	verifier *GoogleVerifier
	issuer   *TokenIssuer
	users    *repository.UserRepository
}

func NewService(
	verifier *GoogleVerifier,
	issuer *TokenIssuer,
	users *repository.UserRepository,
) *Service {
	return &Service{
		verifier: verifier,
		issuer:   issuer,
		users:    users,
	}
}

type LoginResult struct {
	Token string
	User  *model.User
}

func (s *Service) GoogleLogin(
	ctx context.Context,
	idToken string,
) (*LoginResult, error) {
	identity, err := s.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, err
	}

	user, err := s.users.UpsertGoogleUser(
		ctx,
		identity.Sub,
		identity.Email,
		identity.Name,
	)
	if err != nil {
		return nil, err
	}

	token, err := s.issuer.Issue(
		user.ID.Hex(),
		user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
		User:  user,
	}, nil
}	