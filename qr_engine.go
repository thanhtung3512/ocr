package ocrworker
import (
	"bytes"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type QREngine struct {
}

func (m QREngine) ProcessRequest(ocrRequest OcrRequest) (OcrResult, error) {
	bmp := bytes.NewReader(ocrRequest.ImgBytes)
	// prepare BinaryBitmap
	//bmp, _ := gozxing.NewBinaryBitmapFromImage(fi)
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil{
    	failOnError("Failed to qr-decode image: %v", err)
	    return OcrResult{Text: err.Error()}, nil
	}
	

	return OcrResult{Text: result.GetText()}, nil
}
