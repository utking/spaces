package services_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/utking/spaces/internal/application/domain"
	"github.com/utking/spaces/internal/application/services"
	"github.com/utking/spaces/internal/ports"
)

func TestGetLastOpened(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetLastOpened", mock.Anything, domain.LastOpenedTypeBookmark, "testUser").
		Return("someID", nil).Once()

	lastOpenedService := services.NewLastOpenedService(dbPort)

	userID := "testUser"

	lastOpenedID, err := lastOpenedService.GetLastOpened(
		t.Context(), domain.LastOpenedTypeBookmark, userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if lastOpenedID == "" || lastOpenedID != "someID" {
		t.Errorf("expected last opened ID to be 'someID', got '%s'", lastOpenedID)
	}
}

func TestSetLastOpened(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("SetLastOpened", mock.Anything, domain.LastOpenedTypeBookmark, "testUser", "someID").
		Return(nil).Once()

	lastOpenedService := services.NewLastOpenedService(dbPort)

	userID := "testUser"
	itemID := "someID"

	err := lastOpenedService.SetLastOpened(
		t.Context(), domain.LastOpenedTypeBookmark, userID, itemID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	dbPort.AssertExpectations(t)
}
