package ocrworker
import (
        "context"
        "bytes"
        "github.com/couchbaselabs/logg"
        "encoding/json"

        vision "cloud.google.com/go/vision/apiv1"

        pb "google.golang.org/genproto/googleapis/cloud/vision/v1"

  		nlu "github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
)

const VISION_ENGINE_RESPONSE = "vision engine decoder response"

type VisionEngine struct {
	Annotations 		[]*pb.EntityAnnotation
	Keywords 			*nlu.AnalysisResults
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
		  },
		},
	)
	failOnError("Failed to get IBM API response: %v", err)
	keywords := req.GetAnalyzeResult(response)
	return keywords, nil
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

    return texts, nil
}

func (m VisionEngine) ProcessRequest(ocrRequest OcrRequest) (OcrResult, error) {

    // Creates a VISION client.
    texts, err := Vision(ocrRequest)
    text := texts[0].GetDescription()
    
	visionEngine := VisionEngine{Annotations: texts}

    // Creates a IBM API NLP client
    licenses := []License{}
	l1 := License{"fbb08187-7a55-4eb8-a859-b6b4b4127e5b", "QcOZNFoifw2e"}
	licenses = append(licenses, l1)

	for _, license := range licenses {
		keywords, err := NaturalLanguageUnderstanding(text, license)
		failOnError("Failed to use NLU: %v", err)
		visionEngine.Keywords = keywords
	}

	// Put together
    //visionEngine := VisionEngine{Annotations: texts, Keywords: keywords}

    res, err := json.Marshal(visionEngine)
    failOnError("Failed to convert to json: %v", err)

	return OcrResult{Text: string(res)}, nil
	//return OcrResult{Text: VISION_ENGINE_RESPONSE}, nil
}
