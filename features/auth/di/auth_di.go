package di

import (
	"github.com/RodolfoBonis/hermes/core/config"
	"github.com/RodolfoBonis/hermes/core/services"
	"github.com/RodolfoBonis/hermes/features/auth/domain/usecases"
)

func AuthInjection() usecases.AuthUseCase {
	return usecases.AuthUseCase{
		KeycloakClient:     services.AuthClient,
		KeycloakAccessData: config.EnvKeyCloak(),
	}
}
