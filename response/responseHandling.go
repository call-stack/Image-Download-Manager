package response

import (
	"encoding/json"
	"io"
)

type DownloadSuccessResponse struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	INTERNALCODE int    `json:"internal_code"`
	MESSAGE      string `json:"message"`
}

func (er *ErrorResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(er)
}

func (su *DownloadSuccessResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(su)
}
