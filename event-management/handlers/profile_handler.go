package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"

	"example.com/learning/event-management/dto"
	"example.com/learning/event-management/models"
	"example.com/learning/event-management/repository"
	"example.com/learning/event-management/storage"
	"github.com/gin-gonic/gin"
)

const (
	profilePictureFormField = "picture"

	// Five mebibytes.
	maxProfilePictureSize int64 = 5 * 1024 * 1024
)

var allowedProfilePictureTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type ProfileHandler struct {
	userRepository repository.UserRepository
	objectStorage  storage.Store
}

func NewProfileHandler(
	userRepository repository.UserRepository,
	objectStorage storage.Store,
) *ProfileHandler {
	return &ProfileHandler{
		userRepository: userRepository,
		objectStorage:  objectStorage,
	}
}

// GetProfile returns the authenticated user's profile.
func (h *ProfileHandler) GetProfile(
	c *gin.Context,
) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	profile, err := h.userRepository.GetProfile(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		h.handleRepositoryError(c, err)
		return
	}

	response, err := h.profileResponse(profile)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to create profile response",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"profile": response,
		},
	)
}

// UpdateProfile updates the authenticated user's name and/or bio.
func (h *ProfileHandler) UpdateProfile(
	c *gin.Context,
) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	var request dto.UpdateProfileRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	if request.Name == nil && request.Bio == nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "at least one of name or bio must be provided",
			},
		)
		return
	}

	normalizeProfileRequest(&request)

	profile, err := h.userRepository.UpdateProfile(
		c.Request.Context(),
		userID,
		request.Name,
		request.Bio,
	)
	if err != nil {
		h.handleRepositoryError(c, err)
		return
	}

	response, err := h.profileResponse(profile)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to create profile response",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"profile": response,
		},
	)
}

// UploadProfilePicture validates and stores a profile picture.
//
// The request must use multipart/form-data and contain a file field named
// "picture".
func (h *ProfileHandler) UploadProfilePicture(
	c *gin.Context,
) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	fileHeader, err := c.FormFile(profilePictureFormField)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "profile picture is required in the picture form field",
			},
		)
		return
	}

	if fileHeader.Size <= 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "profile picture is empty",
			},
		)
		return
	}

	if fileHeader.Size > maxProfilePictureSize {
		c.JSON(
			http.StatusRequestEntityTooLarge,
			gin.H{
				"error": "profile picture must not exceed 5 MB",
			},
		)
		return
	}

	uploadedFile, err := fileHeader.Open()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "failed to open uploaded profile picture",
			},
		)
		return
	}
	defer uploadedFile.Close()

	extension, err := validateProfilePicture(uploadedFile)
	if err != nil {
		c.JSON(
			http.StatusUnsupportedMediaType,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	if _, err := uploadedFile.Seek(0, io.SeekStart); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to process uploaded profile picture",
			},
		)
		return
	}

	objectKey, err := buildProfilePictureKey(
		userID,
		extension,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to generate profile-picture name",
			},
		)
		return
	}

	storedObject, err := h.objectStorage.Put(
		c.Request.Context(),
		objectKey,
		uploadedFile,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to store profile picture",
			},
		)
		return
	}

	oldObjectKey, err := h.userRepository.SetProfilePicture(
		c.Request.Context(),
		userID,
		storedObject.Key,
	)
	if err != nil {
		// The database was not updated, so remove the newly stored file to
		// avoid leaving an unused object behind.
		cleanupErr := h.objectStorage.Delete(
			c.Request.Context(),
			storedObject.Key,
		)
		if cleanupErr != nil {
			log.Printf(
				"failed to clean up new profile picture %q: %v",
				storedObject.Key,
				cleanupErr,
			)
		}

		h.handleRepositoryError(c, err)
		return
	}

	// Database now points to the new object. The old object can be removed.
	if oldObjectKey != "" && oldObjectKey != storedObject.Key {
		if err := h.objectStorage.Delete(
			c.Request.Context(),
			oldObjectKey,
		); err != nil {
			// Do not fail the request because the user's new profile picture
			// is already valid and stored in the database.
			log.Printf(
				"failed to delete previous profile picture %q: %v",
				oldObjectKey,
				err,
			)
		}
	}

	profile, err := h.userRepository.GetProfile(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		h.handleRepositoryError(c, err)
		return
	}

	response, err := h.profileResponse(profile)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to create profile response",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"profile": response,
		},
	)
}

// DeleteProfilePicture removes the authenticated user's current picture.
func (h *ProfileHandler) DeleteProfilePicture(
	c *gin.Context,
) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	oldObjectKey, err := h.userRepository.RemoveProfilePicture(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		h.handleRepositoryError(c, err)
		return
	}

	if oldObjectKey != "" {
		if err := h.objectStorage.Delete(
			c.Request.Context(),
			oldObjectKey,
		); err != nil {
			// PostgreSQL has already been updated. The API state is correct,
			// so log the orphaned-file cleanup failure instead of returning a
			// misleading request failure.
			log.Printf(
				"failed to delete profile picture %q: %v",
				oldObjectKey,
				err,
			)
		}
	}

	profile, err := h.userRepository.GetProfile(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		h.handleRepositoryError(c, err)
		return
	}

	response, err := h.profileResponse(profile)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to create profile response",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "profile picture removed",
			"profile": response,
		},
	)
}

func (h *ProfileHandler) profileResponse(
	profile models.Profile,
) (dto.UserResponse, error) {
	var pictureURL *string

	if profile.ProfilePictureKey != "" {
		value, err := h.objectStorage.URL(
			profile.ProfilePictureKey,
		)
		if err != nil {
			return dto.UserResponse{}, fmt.Errorf(
				"create profile-picture URL: %w",
				err,
			)
		}

		pictureURL = &value
	}

	return dto.UserResponse{
		Id:                profile.UserID,
		Email:             profile.Email,
		Name:              profile.Name,
		Bio:               profile.Bio,
		ProfilePictureURL: pictureURL,
	}, nil
}

func (h *ProfileHandler) handleRepositoryError(
	c *gin.Context,
	err error,
) {
	if errors.Is(err, repository.ErrUserNotFound) {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "user profile not found",
			},
		)
		return
	}

	c.JSON(
		http.StatusInternalServerError,
		gin.H{
			"error": "profile operation failed",
		},
	)
}

func authenticatedUserID(
	c *gin.Context,
) (int, bool) {
	value, exists := c.Get("userId")
	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "authenticated user is missing",
			},
		)
		return 0, false
	}

	userID, ok := value.(int)
	if !ok || userID <= 0 {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "authenticated user is invalid",
			},
		)
		return 0, false
	}

	return userID, true
}

func normalizeProfileRequest(
	request *dto.UpdateProfileRequest,
) {
	if request.Name != nil {
		value := strings.TrimSpace(*request.Name)
		request.Name = &value
	}

	if request.Bio != nil {
		value := strings.TrimSpace(*request.Bio)
		request.Bio = &value
	}
}

func validateProfilePicture(
	file multipart.File,
) (string, error) {
	header := make([]byte, 512)

	readCount, err := file.Read(header)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf(
			"read profile picture: %w",
			err,
		)
	}

	if readCount == 0 {
		return "", errors.New("profile picture is empty")
	}

	contentType := http.DetectContentType(
		header[:readCount],
	)

	extension, allowed := allowedProfilePictureTypes[contentType]
	if !allowed {
		return "", errors.New(
			"only JPEG, PNG, and WebP profile pictures are supported",
		)
	}

	return extension, nil
}

func buildProfilePictureKey(
	userID int,
	extension string,
) (string, error) {
	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf(
			"generate random filename: %w",
			err,
		)
	}

	randomName := hex.EncodeToString(randomBytes)

	return path.Join(
		"users",
		strconv.Itoa(userID),
		"profile",
		randomName+extension,
	), nil
}
