package ocrworker
import (
        "context"
        "bytes"
        "strings"
        "github.com/couchbaselabs/logg"
        "encoding/json"

        vision "cloud.google.com/go/vision/apiv1"

        pb "google.golang.org/genproto/googleapis/cloud/vision/v1"

  		nlu "github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
)

type VisionEngine struct {
	Annotations 		[]*pb.EntityAnnotation
	Tags 				*nlu.AnalysisResults
}

type License struct {
	User   string 
    Pass   string 
}

func failOnError(msg string, err error) {
    if err != nil {
        logg.LogTo("OCR_HTTP", msg, err)
    }
}

func NaturalLanguageUnderstanding(text string, license License) (*nlu.AnalysisResults, error) {
	req, err := nlu.NewNaturalLanguageUnderstandingV1(&nlu.NaturalLanguageUnderstandingV1Options{
		URL: "https://gateway.watsonplatform.net/natural-language-understanding/api",
		Version: "2019-07-12",
		Username: license.User,
		Password: license.Pass,
	})
	failOnError("Failed to send IBM API request: %v", err)

	sentiment := true
	emotion := true
	limit := int64(1000)

	response, err := req.Analyze(
		&nlu.AnalyzeOptions{
			Text: &text,
			Features: &nlu.Features{
				Keywords: &nlu.KeywordsOptions{
					Sentiment: &sentiment,
					Emotion: &emotion,
					Limit: &limit,
				},
				Entities: &nlu.EntitiesOptions{
					Sentiment: &sentiment,
					Limit: &limit,
				},
			},
		},
	)
	failOnError("Failed to get IBM API response: %v", err)
	tags := req.GetAnalyzeResult(response)
	return tags, err
}

func Vision(ocrRequest OcrRequest) ([]*pb.EntityAnnotation, error) {
	ctx := context.Background()
    client, err := vision.NewImageAnnotatorClient(ctx)
    failOnError("Failed to create client: %v", err)
    defer client.Close()

    image, err := vision.NewImageFromReader(bytes.NewReader(ocrRequest.ImgBytes))
    failOnError("Failed to create image: %v", err)

    texts, err := client.DetectTexts(ctx, image, nil, 1000)
    failOnError("Failed to detect text: %v", err)

    return texts, err
}

func (m VisionEngine) ProcessRequest(ocrRequest OcrRequest) (OcrResult, error) {

    // Creates a VISION client.
    texts, err := Vision(ocrRequest)
    
	visionEngine := VisionEngine{Annotations: texts}

	if err == nil && len(texts) > 0 {
	    text := texts[0].GetDescription()
	    // Creates a IBM API NLP client
	    licenses := []License{}
		licenses = append(licenses, License{"fbb08187-7a55-4eb8-a859-b6b4b4127e5b", "QcOZNFoifw2e"})
		licenses = append(licenses, License{"a02d69a0-6f92-4755-95d3-cafa5a9be186", "iEeWsosMO2uE"})
		licenses = append(licenses, License{"1aaf739b-c1bc-4559-a9f4-d883bb04715d", "S3NCvzi57pJC"})
		licenses = append(licenses, License{"ac4eca2c-6f85-4f0f-be0b-98199c87c4b9", "v3yBpL8k2ytd"})
		licenses = append(licenses, License{"20b6cc50-1639-4eca-a5f1-b66329c46a63", "AnrVkyhbbkzq"})
		licenses = append(licenses, License{"c193cd8d-58b3-4aa5-a45c-4496cb9ce7b2", "kbofmE2GNChM"})
		licenses = append(licenses, License{"6d1fdb3f-42a8-4ebb-bb6b-9d6e6c52eb29", "zmOwL5rt5WIH"})
		licenses = append(licenses, License{"779333db-6241-45d0-9e59-78e9687e22e5", "OAP8UsDBdqht"})
		licenses = append(licenses, License{"7119f207-8191-4c63-8047-da0ce13cda55", "OWvAvYYsjCZu"})
		for _, license := range licenses {
			tags, err := NaturalLanguageUnderstanding(text, license)
			if err != nil{
				if strings.Contains(err.Error(),"Code: 422") || strings.Contains(err.Error(),"Code: 400") || strings.Contains(err.Error(),"unsupported text language") || strings.Contains(err.Error(),"unknown language detected") {
					break
				}
			} else {
				visionEngine.Tags = tags
				break
			}
		}
	}else{
		visionEngine.Tags = nil
	}
    
    res, err := json.Marshal(visionEngine)
    failOnError("Failed to convert to json: %v", err)

	return OcrResult{Text: string(res)}, nil
}
