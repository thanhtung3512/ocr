package ocrworker
import (
	"bytes"
	"image"
	_ "image/jpeg"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type QREngine struct {
}

func (m QREngine) ProcessRequest(ocrRequest OcrRequest) (OcrResult, error) {
	fi := bytes.NewReader(ocrRequest.ImgBytes)
	img, _, err := image.Decode(fi)
	if err != nil{
    	failOnError("Failed to qr-decode image: %v", err)
	    return OcrResult{Text: err.Error()}, nil
	}
	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil{
    	failOnError("Failed to qr-decode image: %v", err)
	    return OcrResult{Text: err.Error()}, nil
	}
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil{
    	failOnError("Failed to qr-decode image: %v", err)
	    return OcrResult{Text: err.Error()}, nil
	}
	

	return OcrResult{Text: result.GetText()}, nil
}
