
package ocrworker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbaselabs/logg"
)

type OcrHttpHandler struct {
	RabbitConfig RabbitConfig
}


