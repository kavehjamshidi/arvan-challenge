package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kavehjamshidi/arvan-challenge/domain"
	contract2 "github.com/kavehjamshidi/arvan-challenge/service/upload/contract"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
	"net/http"
)

type UploadHandler struct {
	uploadService contract2.UploadService
}

func NewUploadHandler(
	uploadService contract2.UploadService,
) UploadHandler {
	return UploadHandler{
		uploadService: uploadService,
	}
}

func (u UploadHandler) HandleUpload(c *fiber.Ctx) error {
	userID := c.Get("user-id")

	multipartFile, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response[string]{
			Message: MsgValidationError,
			Error:   "could not find file in form data",
		})
	}

	fileDescriptor, err := multipartFile.Open()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response[string]{
			Message: MsgValidationError,
			Error:   "could not open file data",
		})
	}
	defer fileDescriptor.Close()

	fileID := c.FormValue("file_id")
	if fileID == "" {
		return c.Status(http.StatusBadRequest).JSON(response[string]{
			Message: MsgValidationError,
			Error:   "file_id is required",
		})
	}

	file := domain.File{
		Data:   fileDescriptor,
		Size:   multipartFile.Size,
		FileID: fileID,
		UserID: userID,
	}

	err = u.uploadService.UploadFile(c.Context(), file)
	if err != nil {
		return u.handleUploadError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response[*string]{
		Message: MsgSuccess,
		Error:   nil,
	})
}

func (u UploadHandler) handleUploadError(c *fiber.Ctx, err error) error {
	if errors.Is(err, utils.ErrFileIDAlreadyExists) {
		return c.Status(http.StatusConflict).JSON(response[string]{
			Error:   err.Error(),
			Message: MsgFailed,
		})
	}

	if errors.Is(err, utils.ErrUsageLimitExceeded) {
		return c.Status(http.StatusForbidden).JSON(response[string]{
			Error:   err.Error(),
			Message: MsgForbidden,
		})
	}

	return c.Status(http.StatusInternalServerError).JSON(response[string]{
		Error:   err.Error(),
		Message: MsgFailed,
	})
}
