package room

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/m1k1o/neko/server/pkg/utils"
)

// TODO: Extract file uploading to custom utility.

// multipartFormMaxMemory controls how much of the form is kept in memory.
const multipartFormMaxMemory = 32 << 20

// multipartFormMaxSize is the hard request-body limit. ParseMultipartForm's
// maxMemory argument alone does not limit the total amount written to disk.
const multipartFormMaxSize = 32 << 20

func parseUploadForm(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, multipartFormMaxSize)
	return r.ParseMultipartForm(multipartFormMaxMemory)
}

func writeUploadedFile(dir string, header *multipart.FileHeader) (string, error) {
	filename := filepath.Base(header.Filename)
	if filename == "." || filename == string(filepath.Separator) || filename == "" {
		return "", os.ErrInvalid
	}

	target := filepath.Join(dir, filename)
	srcFile, err := header.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		_ = dstFile.Close()
		return "", err
	}
	if err := dstFile.Close(); err != nil {
		return "", err
	}

	return target, nil
}

func (h *RoomHandler) uploadDrop(w http.ResponseWriter, r *http.Request) error {
	if !h.desktop.IsUploadDropEnabled() {
		return utils.HttpBadRequest("upload drop is disabled")
	}

	err := parseUploadForm(w, r)
	if err != nil {
		return utils.HttpBadRequest("failed to parse multipart form").WithInternalErr(err)
	}

	//nolint
	defer r.MultipartForm.RemoveAll()

	X, err := strconv.Atoi(r.FormValue("x"))
	if err != nil {
		return utils.HttpBadRequest("no X coordinate received").WithInternalErr(err)
	}

	Y, err := strconv.Atoi(r.FormValue("y"))
	if err != nil {
		return utils.HttpBadRequest("no Y coordinate received").WithInternalErr(err)
	}

	req_files := r.MultipartForm.File["files"]
	if len(req_files) == 0 {
		return utils.HttpBadRequest("no files received")
	}

	dir, err := os.MkdirTemp("", "neko-drop-*")
	if err != nil {
		return utils.HttpInternalServerError().
			WithInternalErr(err).
			WithInternalMsg("unable to create temporary directory")
	}
	defer os.RemoveAll(dir)

	files := []string{}
	for _, req_file := range req_files {
		filePath, err := writeUploadedFile(dir, req_file)
		if err != nil {
			return utils.HttpInternalServerError().
				WithInternalErr(err).
				WithInternalMsg("unable to store uploaded file")
		}

		files = append(files, filePath)
	}

	if !h.desktop.DropFiles(X, Y, files) {
		return utils.HttpInternalServerError().
			WithInternalMsg("unable to drop files")
	}

	return utils.HttpSuccess(w)
}

func (h *RoomHandler) uploadDialogPost(w http.ResponseWriter, r *http.Request) error {
	if !h.desktop.IsFileChooserDialogEnabled() {
		return utils.HttpBadRequest("file chooser dialog is disabled")
	}

	err := parseUploadForm(w, r)
	if err != nil {
		return utils.HttpBadRequest("failed to parse multipart form").WithInternalErr(err)
	}

	//nolint
	defer r.MultipartForm.RemoveAll()

	req_files := r.MultipartForm.File["files"]
	if len(req_files) == 0 {
		return utils.HttpBadRequest("no files received")
	}

	if !h.desktop.IsFileChooserDialogOpened() {
		return utils.HttpUnprocessableEntity("file chooser dialog is not open")
	}

	dir, err := os.MkdirTemp("", "neko-dialog-*")
	if err != nil {
		return utils.HttpInternalServerError().
			WithInternalErr(err).
			WithInternalMsg("unable to create temporary directory")
	}
	defer os.RemoveAll(dir)

	for _, req_file := range req_files {
		_, err := writeUploadedFile(dir, req_file)
		if err != nil {
			return utils.HttpInternalServerError().
				WithInternalErr(err).
				WithInternalMsg("unable to store uploaded file")
		}
	}

	if err := h.desktop.HandleFileChooserDialog(dir); err != nil {
		return utils.HttpInternalServerError().
			WithInternalErr(err).
			WithInternalMsg("unable to handle file chooser dialog")
	}

	return utils.HttpSuccess(w)
}

func (h *RoomHandler) uploadDialogClose(w http.ResponseWriter, r *http.Request) error {
	if !h.desktop.IsFileChooserDialogEnabled() {
		return utils.HttpBadRequest("file chooser dialog is disabled")
	}

	if !h.desktop.IsFileChooserDialogOpened() {
		return utils.HttpUnprocessableEntity("file chooser dialog is not open")
	}

	h.desktop.CloseFileChooserDialog()

	return utils.HttpSuccess(w)
}
