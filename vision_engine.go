package ocrworker
import (
        "context"
        "bytes"
        "strings"
        "github.com/couchbaselabs/logg"
        "github.com/thanhtung3512/go-yandex-translate"
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

type Yandex struct {
	Key   string
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
	    locale := texts[0].GetLocale()
	    if locale != "en" && locale!="fr" && locale!="de" && locale!="it" && locale!="ja" && locale!="ko" && locale!="es" && locale!="pt"{
	    	// Create Yandex license
	    	yandexes := []Yandex{}
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T215833Z.53e8a9545764e42a.05662f448bff7951369fc863f696b7ae10dcae40"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T215559Z.80e66b2ce2b6344e.589eb00b3422b43c1379ce6bd0972bdeb53b25ef"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T215457Z.560437ac05817551.42e38032c2be3c057c1dfb19084bf96f73c85603"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T214703Z.78edcdd889060392.1d694d8bf5304c60305eb18c0fa6a05a0fe3da83"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T214749Z.3f9ceb70b0407de1.54639f4722775701f25773159dbe91a3b6173afc"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T214834Z.486edfa3dc719010.0d4e303aad316731a9574266d47f6fcda7c6abbc"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T214930Z.23a6a12078c05a99.9a33442bd218e1ee0a279e3ab7eddc2b6183e1ad"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T215013Z.3163f5743997d9ac.06ca632faae688a988e0982dcd7c83ecac7fbe41"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T214454Z.9cda750a7b3ed883.a067a0dc5ab658795a7c25f43a9a0382c1d52d39"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T200902Z.b5763e9ac3f4602b.8e4952fac4164b34c79ae7a3be4d9b61541e2d4c"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T200651Z.6cc8679b9a6a0604.eb7a5fd6a785901839469779e7bfd40faecbd67a"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T200348Z.b71d027f04df6071.478d93655155eb3652b340be61f9e8e7467ce4e7"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180213T090548Z.258a0dd7469c56cb.c6193105ab25a4257312f83f8e9376841216524c"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180213T090802Z.b53bb4c02b0bd670.bc900c25978e58fd5839e5c15d3c5c32bc3ad848"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20170608T132524Z.03f18d029e3e657a.7a39d518dfbc6b8c53d4163e2dfeaa1f2d4ae9cb"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20170609T181312Z.b1c2b14a6edb6c6f.7ac83a938d6aa4e1e9cadc72c510f4dbc4ceab72"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180127T132223Z.da01b637e1fa0f80.351255755dc70184abc74c8b2c5ecad90ee5061e"})
	    	yandexes = append(yandexes, Yandex{"trnsl.1.1.20180312T200205Z.ecbb13c549ef6e65.1984f46d4fbc162048f9fbe76c92cf245bcd0e38"})

	    	for _, yandex := range yandexes {
	    		tr := translate.New(yandex.Key)
	    		translation, err := tr.Translate("en", text)
				if err != nil {
					// error
				} else {
					text = translation.Result()
					break
				}
			}
	    }

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
