package httpserver

import (
	"net/http"

	"github.com/ankorstore/yokai/log"
	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
)

// JsonErrorHandler provides a [echo.HTTPErrorHandler] that outputs errors in JSON format.
// It can also be configured to obfuscate error message (to avoid to leak sensitive details), and to add the error stack to the response.
type JsonErrorHandler struct {
	obfuscate bool
	stack     bool
}

// NewJsonErrorHandler returns a new JsonErrorHandler instance.
func NewJsonErrorHandler(obfuscate bool, stack bool) *JsonErrorHandler {
	return &JsonErrorHandler{
		obfuscate: obfuscate,
		stack:     stack,
	}
}

// Handle handles errors.
func (h *JsonErrorHandler) Handle() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		logger := log.CtxLogger(c.Request().Context())

		if c.Response().Committed {
			return
		}

		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			if httpError.Internal != nil {
				var internalHttpError *echo.HTTPError
				if errors.As(httpError.Internal, &internalHttpError) {
					httpError = internalHttpError
				}
			}
		} else {
			httpError = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		var logRespFields map[string]interface{}

		if h.stack {
			errStack := "n/a"
			if err != nil {
				errStack = errors.New(err).ErrorStack()
			}

			switch m := httpError.Message.(type) {
			case error:
				logRespFields = map[string]interface{}{
					"message": m.Error(),
					"stack":   errStack,
				}
			default:
				logRespFields = map[string]interface{}{
					"message": m,
					"stack":   errStack,
				}
			}
		} else {
			switch m := httpError.Message.(type) {
			case error:
				logRespFields = map[string]interface{}{
					"message": m.Error(),
				}
			default:
				logRespFields = map[string]interface{}{
					"message": m,
				}
			}
		}

		logger.Error().Err(err).Fields(logRespFields).Msg("error handler")

		httpRespFields := logRespFields

		if h.obfuscate {
			httpRespFields["message"] = http.StatusText(httpError.Code)
		}

		var httpRespErr error
		if c.Request().Method == http.MethodHead {
			httpRespErr = c.NoContent(httpError.Code)
		} else {
			httpRespErr = c.JSON(httpError.Code, httpRespFields)
		}

		if httpRespErr != nil {
			logger.Error().Err(httpRespErr).Msg("error handler failure")
		}
	}
}
