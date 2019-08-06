package ocrworker
import (
	"bytes"

	"github.com/tuotoo/qrcode"
)

type QREngine struct {
}

func (m QREngine) ProcessRequest(ocrRequest OcrRequest) (OcrResult, error) {
	fi := bytes.NewReader(ocrRequest.ImgBytes)

	qrmatrix, err := qrcode.Decode(fi)
	if err != nil{
    	failOnError("Failed to qr-decode image: %v", err)
	    //return
	}

	return OcrResult{Text: qrmatrix.Content}, nil
}
