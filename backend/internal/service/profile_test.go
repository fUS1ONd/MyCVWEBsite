package service

import (
	"context"
	"errors"
	"testing"

	"personal-web-platform/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProfileRepository is a mock implementation of ProfileRepository
type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) GetProfile(ctx context.Context) (*domain.Profile, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Profile), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockProfileRepository) UpdateProfile(ctx context.Context, req *domain.UpdateProfileRequest) (*domain.Profile, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Profile), args.Error(1) //nolint:errcheck // mock method
}

func TestProfileService_GetProfile(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockProfileRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - profile found",
			setupMock: func(m *MockProfileRepository) {
				m.On("GetProfile", mock.Anything).Return(&domain.Profile{
					ID:          1,
					Name:        "John Doe",
					Description: "Software Engineer",
					Activity:    "Building great products",
					Contacts: domain.Contacts{
						Email:  "john@example.com",
						GitHub: "https://github.com/johndoe",
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - profile not found",
			setupMock: func(m *MockProfileRepository) {
				m.On("GetProfile", mock.Anything).Return(nil, nil)
			},
			wantErr:     true,
			errContains: "profile not found",
		},
		{
			name: "error - repository error",
			setupMock: func(m *MockProfileRepository) {
				m.On("GetProfile", mock.Anything).Return(nil, errors.New("database error"))
			},
			wantErr:     true,
			errContains: "failed to get profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProfileRepository)
			tt.setupMock(mockRepo)

			service := NewProfileService(mockRepo)
			profile, err := service.GetProfile(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, profile)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, "John Doe", profile.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProfileService_UpdateProfile(t *testing.T) {
	validRequest := &domain.UpdateProfileRequest{
		Name:        "Jane Doe",
		Description: "Senior Software Engineer",
		Activity:    "Leading development teams",
		Contacts: domain.Contacts{
			Email:    "jane@example.com",
			GitHub:   "https://github.com/janedoe",
			LinkedIn: "https://linkedin.com/in/janedoe",
		},
	}

	tests := []struct {
		name        string
		request     *domain.UpdateProfileRequest
		setupMock   func(*MockProfileRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:    "success - profile updated",
			request: validRequest,
			setupMock: func(m *MockProfileRepository) {
				m.On("UpdateProfile", mock.Anything, validRequest).Return(&domain.Profile{
					ID:          1,
					Name:        "Jane Doe",
					Description: "Senior Software Engineer",
					Activity:    "Leading development teams",
					Contacts: domain.Contacts{
						Email:    "jane@example.com",
						GitHub:   "https://github.com/janedoe",
						LinkedIn: "https://linkedin.com/in/janedoe",
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - validation failed (empty name)",
			request: &domain.UpdateProfileRequest{
				Name:        "",
				Description: "Description",
				Activity:    "Activity",
				Contacts:    domain.Contacts{Email: "test@test.com"},
			},
			setupMock:   func(m *MockProfileRepository) {},
			wantErr:     true,
			errContains: "validation failed",
		},
		{
			name: "error - validation failed (invalid email)",
			request: &domain.UpdateProfileRequest{
				Name:        "Name",
				Description: "Description",
				Activity:    "Activity",
				Contacts:    domain.Contacts{Email: "invalid-email"},
			},
			setupMock:   func(m *MockProfileRepository) {},
			wantErr:     true,
			errContains: "validation failed",
		},
		{
			name:    "error - repository error",
			request: validRequest,
			setupMock: func(m *MockProfileRepository) {
				m.On("UpdateProfile", mock.Anything, validRequest).Return(nil, errors.New("database error"))
			},
			wantErr:     true,
			errContains: "failed to update profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProfileRepository)
			tt.setupMock(mockRepo)

			service := NewProfileService(mockRepo)
			profile, err := service.UpdateProfile(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, profile)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, tt.request.Name, profile.Name)
				assert.Equal(t, tt.request.Description, profile.Description)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
