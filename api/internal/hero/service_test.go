package hero

import (
	"createtodayapi/internal/config"
	"createtodayapi/internal/infra"
	"createtodayapi/internal/logger"
	"github.com/stretchr/testify/require"
	"testing"
)

func NewTestService() *Service {
	conf := config.New("../../.env")
	log := logger.New()

	db, err := infra.InitPostgres(conf.DatabaseDSN)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	postgres := NewPostgresRepo(db)
	memory := NewMemoryRepo()
	emailsService := NewEmailService(conf, memory)
	service := NewService(postgres, conf, emailsService)
	return service
}

func TestGeneratePassword(t *testing.T) {
	t.Parallel()
	service := NewTestService()
	t.Run("should generate a new password", func(t *testing.T) {
		password, err := service.generatePassword()
		require.NoError(t, err)
		require.Equal(t, len(password), 8)
		t.Log(password)
	})
}
