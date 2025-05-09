package pleasant

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrPrereqNotMet        = errors.New("error: not all prerequisites met")
	ErrNotFound            = errors.New("error: the requested resource was not found")
	ErrBadRequest          = errors.New("error: bad request")
	ErrUnauthorized        = errors.New("error: unauthorized, please log in again")
	ErrInvalidCredentials  = errors.New("error: invalid or incomplete credentials")
	ErrPathStartIncorrect  = errors.New("error: path must start with 'Root/'")
	ErrInvalidResourceType = errors.New("error: invalid resource type, must be 'entry' or 'folder'")
	ErrNoResult            = errors.New("error: no matching entries or folders")
	ErrParentNotFound      = errors.New("error: parent folder not found")
	ErrAmbiguousResult     = errors.New("error: ambiguous result, multiple matching entries or folders")
	ErrLastPathComp        = errors.New("error: last path component is empty")
	ErrDuplicateEntry      = errors.New("error: duplicate entry found, skipping creation")
	ErrDuplicateFolder     = errors.New("error: duplicate folder found, skipping creation")
	ErrArchiveNotEnabled   = errors.New("error: entry/folder/accessrowid does not exist or archiving is possibly disabled")
)

func generateError(res *http.Response) error {
	defer res.Body.Close()

	body, _ := decodeBody(res.Body)

	switch res.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	default:
		return fmt.Errorf("error: HTTP %v %v: %v", res.StatusCode, http.StatusText(res.StatusCode), body)
	}
}
