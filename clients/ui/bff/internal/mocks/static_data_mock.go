package mocks

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

func GetRegisteredModelMocks() []openapi.RegisteredModel {
	model1 := openapi.RegisteredModel{
		CustomProperties:         newCustomProperties(),
		Name:                     "Model One",
		Description:              stringToPointer("This model does things and stuff"),
		ExternalId:               stringToPointer("934589798"),
		Id:                       stringToPointer("1"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		Owner:                    stringToPointer("Sherlock Holmes"),
		State:                    stateToPointer(openapi.REGISTEREDMODELSTATE_LIVE),
	}

	model2 := openapi.RegisteredModel{
		CustomProperties:         newCustomProperties(),
		Name:                     "Model Two",
		Description:              stringToPointer("This model does things and stuff"),
		ExternalId:               stringToPointer("345235987"),
		Id:                       stringToPointer("2"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		Owner:                    stringToPointer("John Watson"),
		State:                    stateToPointer(openapi.REGISTEREDMODELSTATE_LIVE),
	}

	model3 := openapi.RegisteredModel{
		CustomProperties:         newCustomProperties(),
		Name:                     "Model Three",
		Description:              stringToPointer("This model does things and stuff"),
		ExternalId:               stringToPointer("345235989"),
		Id:                       stringToPointer("3"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249933"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249933"),
		Owner:                    stringToPointer("M. Oriarty"),
		State:                    stateToPointer(openapi.REGISTEREDMODELSTATE_ARCHIVED),
	}

	return []openapi.RegisteredModel{model1, model2, model3}
}

func GetRegisteredModelListMock() openapi.RegisteredModelList {
	models := GetRegisteredModelMocks()

	return openapi.RegisteredModelList{
		NextPageToken: "abcdefgh",
		PageSize:      2,
		Size:          int32(len(models)),
		Items:         models,
	}
}

func GetModelVersionMocks() []openapi.ModelVersion {
	modelVersion1 := openapi.ModelVersion{
		CustomProperties:         newCustomProperties(),
		Name:                     "Version One",
		Description:              stringToPointer("This version improves stuff and things"),
		ExternalId:               stringToPointer("934589798"),
		Id:                       stringToPointer("1"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		RegisteredModelId:        "1",
		Author:                   stringToPointer("Sherlock Holmes"),
		State:                    stateToPointer(openapi.MODELVERSIONSTATE_LIVE),
	}

	modelVersion2 := openapi.ModelVersion{
		CustomProperties:         newCustomProperties(),
		Name:                     "Version Two",
		Description:              stringToPointer("This version improves stuff and things better"),
		ExternalId:               stringToPointer("934589798"),
		Id:                       stringToPointer("2"),
		CreateTimeSinceEpoch:     stringToPointer("1725282259922"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282259922"),
		RegisteredModelId:        "1",
		Author:                   stringToPointer("Sherlock Holmes"),
		State:                    stateToPointer(openapi.MODELVERSIONSTATE_LIVE),
	}

	modelVersion3 := openapi.ModelVersion{
		CustomProperties:         newCustomProperties(),
		Name:                     "Version Three",
		Description:              stringToPointer("This version improves stuff and things"),
		ExternalId:               stringToPointer("934589799"),
		Id:                       stringToPointer("3"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		RegisteredModelId:        "2",
		Author:                   stringToPointer("Sherlock Holmes"),
		State:                    stateToPointer(openapi.MODELVERSIONSTATE_LIVE),
	}

	modelVersion4 := openapi.ModelVersion{
		CustomProperties:         newCustomProperties(),
		Name:                     "Version Four",
		Description:              stringToPointer("This version didn't improve stuff and things"),
		ExternalId:               stringToPointer("934589791"),
		Id:                       stringToPointer("4"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		RegisteredModelId:        "3",
		Author:                   stringToPointer("Sherlock Holmes"),
		State:                    stateToPointer(openapi.MODELVERSIONSTATE_ARCHIVED),
	}

	return []openapi.ModelVersion{modelVersion1, modelVersion2, modelVersion3, modelVersion4}
}

func GetModelVersionListMock() openapi.ModelVersionList {
	versions := GetModelVersionMocks()

	return openapi.ModelVersionList{
		NextPageToken: "abcdefgh",
		PageSize:      2,
		Items:         versions,
		Size:          2,
	}
}

func GetModelArtifactMocks() []openapi.ModelArtifact {
	artifact1 := openapi.ModelArtifact{
		ArtifactType:             stringToPointer("TYPE_ONE"),
		CustomProperties:         newCustomProperties(),
		Description:              stringToPointer("This artifact can do more than you would expect"),
		ExternalId:               stringToPointer("1000001"),
		Uri:                      stringToPointer("http://localhost/artifacts/1"),
		State:                    stateToPointer(openapi.ARTIFACTSTATE_LIVE),
		Name:                     stringToPointer("Artifact One"),
		Id:                       stringToPointer("1"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		ModelFormatName:          stringToPointer("ONNX"),
		StorageKey:               stringToPointer("key1"),
		StoragePath:              stringToPointer("/artifacts/1"),
		ModelFormatVersion:       stringToPointer("1.0.0"),
		ServiceAccountName:       stringToPointer("service-1"),
	}

	artifact2 := openapi.ModelArtifact{
		ArtifactType:             stringToPointer("TYPE_TWO"),
		CustomProperties:         newCustomProperties(),
		Description:              stringToPointer("This artifact can do more than you would expect, but less than you would hope"),
		ExternalId:               stringToPointer("1000002"),
		Uri:                      stringToPointer("http://localhost/artifacts/2"),
		State:                    stateToPointer(openapi.ARTIFACTSTATE_PENDING),
		Name:                     stringToPointer("Artifact Two"),
		Id:                       stringToPointer("2"),
		CreateTimeSinceEpoch:     stringToPointer("1725282249921"),
		LastUpdateTimeSinceEpoch: stringToPointer("1725282249921"),
		ModelFormatName:          stringToPointer("TensorFlow"),
		StorageKey:               stringToPointer("key2"),
		StoragePath:              stringToPointer("/artifacts/2"),
		ModelFormatVersion:       stringToPointer("1.0.0"),
		ServiceAccountName:       stringToPointer("service-2"),
	}

	return []openapi.ModelArtifact{artifact1, artifact2}
}

func GetModelArtifactListMock() openapi.ModelArtifactList {
	return openapi.ModelArtifactList{
		NextPageToken: "abcdefgh",
		PageSize:      2,
		Items:         GetModelArtifactMocks(),
		Size:          2,
	}
}

func newCustomProperties() *map[string]openapi.MetadataValue {
	result := map[string]openapi.MetadataValue{
		"tensorflow": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"pytorch": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"mll": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"rnn": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"AWS_KEY": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "asdf89asdf098asdfa",
				MetadataType: "MetadataStringValue",
			},
		},
		"AWS_PASSWORD": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "*AadfeDs34adf",
				MetadataType: "MetadataStringValue",
			},
		},
	}

	return &result
}

func catalogCustomProperties() *map[string]openapi.MetadataValue {
	result := map[string]openapi.MetadataValue{
		"tensorflow": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"pytorch": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"mll": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"rnn": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"validated": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "",
				MetadataType: "MetadataStringValue",
			},
		},
		"validated_on": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "RHOAI 2.20,RHAIIS 3.0,RHELAI 1.5",
				MetadataType: "MetadataStringValue",
			},
		},
		"AWS_KEY": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "asdf89asdf098asdfa",
				MetadataType: "MetadataStringValue",
			},
		},
		"AWS_PASSWORD": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "*AadfeDs34adf",
				MetadataType: "MetadataStringValue",
			},
		},
	}

	return &result
}

func NewMockSessionContext(parent context.Context) context.Context {
	if parent == nil {
		parent = context.TODO()
	}
	traceId := uuid.NewString()
	ctx := context.WithValue(parent, constants.TraceIdKey, traceId)

	traceLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx = context.WithValue(ctx, constants.TraceLoggerKey, traceLogger)
	return ctx
}

func NewMockSessionContextNoParent() context.Context {
	return NewMockSessionContext(context.TODO())
}

func GenerateMockArtifactList() openapi.ArtifactList {
	var artifacts []openapi.Artifact
	for i := 0; i < 2; i++ {
		artifact := GenerateMockArtifact()
		artifacts = append(artifacts, artifact)
	}

	return openapi.ArtifactList{
		NextPageToken: gofakeit.UUID(),
		PageSize:      int32(gofakeit.Number(1, 20)),
		Size:          int32(len(artifacts)),
		Items:         artifacts,
	}
}

func GenerateMockArtifact() openapi.Artifact {
	modelArtifact := GenerateMockModelArtifact()

	mockData := openapi.Artifact{
		ModelArtifact: &modelArtifact,
	}
	return mockData
}

func GetCatalogModelMocks() []models.CatalogModel {
	sampleModel1 := models.CatalogModel{
		Name:             "repo1/granite-8b-code-instruct",
		Description:      stringToPointer("Granite-8B-Code-Instruct is a 8B parameter model fine tuned from\nGranite-8B-Code-Base on a combination of permissively licensed instruction\ndata to enhance instruction following capabilities including logical\nreasoning and problem-solving skills."),
		Provider:         stringToPointer("provider1"),
		Tasks:            []string{"text-generation", "image-to-text"},
		License:          stringToPointer("apache-2.0"),
		LicenseLink:      stringToPointer("https://www.apache.org/licenses/LICENSE-2.0.txt"),
		Maturity:         stringToPointer("Technology preview"),
		Language:         []string{"ar", "cs", "de", "en", "es", "fr", "it", "ja", "ko", "nl", "pt", "zh"},
		CustomProperties: catalogCustomProperties(),
		Logo:             stringToPointer("data:image/svg+xml;base64,PHN2ZyBpZD0iTGF5ZXJfMSIgZGF0YS1uYW1lPSJMYXllciAxIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxOTIgMTQ1Ij48ZGVmcz48c3R5bGU+LmNscy0xe2ZpbGw6I2UwMDt9PC9zdHlsZT48L2RlZnM+PHRpdGxlPlJlZEhhdC1Mb2dvLUhhdC1Db2xvcjwvdGl0bGU+PHBhdGggZD0iTTE1Ny43Nyw2Mi42MWExNCwxNCwwLDAsMSwuMzEsMy40MmMwLDE0Ljg4LTE4LjEsMTcuNDYtMzAuNjEsMTcuNDZDNzguODMsODMuNDksNDIuNTMsNTMuMjYsNDIuNTMsNDRhNi40Myw2LjQzLDAsMCwxLC4yMi0xLjk0bC0zLjY2LDkuMDZhMTguNDUsMTguNDUsMCwwLDAtMS41MSw3LjMzYzAsMTguMTEsNDEsNDUuNDgsODcuNzQsNDUuNDgsMjAuNjksMCwzNi40My03Ljc2LDM2LjQzLTIxLjc3LDAtMS4wOCwwLTEuOTQtMS43My0xMC4xM1oiLz48cGF0aCBjbGFzcz0iY2xzLTEiIGQ9Ik0xMjcuNDcsODMuNDljMTIuNTEsMCwzMC42MS0yLjU4LDMwLjYxLTE3LjQ2YTE0LDE0LDAsMCwwLS4zMS0zLjQybC03LjQ1LTMyLjM2Yy0xLjcyLTcuMTItMy4yMy0xMC4zNS0xNS43My0xNi42QzEyNC44OSw4LjY5LDEwMy43Ni41LDk3LjUxLjUsOTEuNjkuNSw5MCw4LDgzLjA2LDhjLTYuNjgsMC0xMS42NC01LjYtMTcuODktNS42LTYsMC05LjkxLDQuMDktMTIuOTMsMTIuNSwwLDAtOC40MSwyMy43Mi05LjQ5LDI3LjE2QTYuNDMsNi40MywwLDAsMCw0Mi41Myw0NGMwLDkuMjIsMzYuMywzOS40NSw4NC45NCwzOS40NU0xNjAsNzIuMDdjMS43Myw4LjE5LDEuNzMsOS4wNSwxLjczLDEwLjEzLDAsMTQtMTUuNzQsMjEuNzctMzYuNDMsMjEuNzdDNzguNTQsMTA0LDM3LjU4LDc2LjYsMzcuNTgsNTguNDlhMTguNDUsMTguNDUsMCwwLDEsMS41MS03LjMzQzIyLjI3LDUyLC41LDU1LC41LDc0LjIyYzAsMzEuNDgsNzQuNTksNzAuMjgsMTMzLjY1LDcwLjI4LDQ1LjI4LDAsNTYuNy0yMC40OCw1Ni43LTM2LjY1LDAtMTIuNzItMTEtMjcuMTYtMzAuODMtMzUuNzgiLz48L3N2Zz4="),
		Readme: stringToPointer(`---
pipeline_tag: text-generation
inference: false
license: apache-2.0
library_name: transformers
tags:
- language
- granite-3.1
base_model:
- provider1-granite/granite-3.1-8b-base
---

# Granite-3.1-8B-Instruct

**Model Summary:**
Granite-3.1-8B-Instruct is a 8B parameter long-context instruct model finetuned from Granite-3.1-8B-Base using a combination of open source instruction datasets with permissive license and internally collected synthetic datasets tailored for solving long context problems. This model is developed using a diverse set of techniques with a structured chat format, including supervised finetuning, model alignment using reinforcement learning, and model merging.

- **Developers:** Granite Team, provider1
- **GitHub Repository:** [provider1-granite/granite-3.1-language-models](https://github.com/provider1-granite/granite-3.1-language-models)
- **Website**: [Granite Docs](https://www.provider1.com/granite/docs/)
- **Paper:** [Granite 3.1 Language Models (coming soon)](https://huggingface.co/collections/provider1-granite/granite-31-language-models-6751dbbf2f3389bec5c6f02d) 
- **Release Date**: December 18th, 2024
- **License:** [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0)

**Supported Languages:** 
English, German, Spanish, French, Japanese, Portuguese, Arabic, Czech, Italian, Korean, Dutch, and Chinese. Users may finetune Granite 3.1 models for languages beyond these 12 languages.

**Intended Use:** 
The model is designed to respond to general instructions and can be used to build AI assistants for multiple domains, including business applications.

*Capabilities*
* Summarization
* Text classification
* Text extraction
* Question-answering
* Retrieval Augmented Generation (RAG)
* Code related tasks
* Function-calling tasks
* Multilingual dialog use cases
* Long-context tasks including long document/meeting summarization, long document QA, etc.

**Generation:** 
This is a simple example of how to use Granite-3.1-8B-Instruct model.

Install the following libraries:

` + "```" + `shell
pip install torch torchvision torchaudio
pip install accelerate
pip install transformers
` + "```" + `
Then, copy the snippet from the section that is relevant for your use case.

 ` + "```" + `python
import torch
from transformers import AutoModelForCausalLM, AutoTokenizer

device = "auto"
model_path = "provider1-granite/granite-3.1-8b-instruct"
tokenizer = AutoTokenizer.from_pretrained(model_path)
# drop device_map if running on CPU
model = AutoModelForCausalLM.from_pretrained(model_path, device_map=device)
model.eval()
# change input text as desired
chat = [
    { "role": "user", "content": "Please list one provider1 Research laboratory located in the United States. You should only output its name and location." },
]
chat = tokenizer.apply_chat_template(chat, tokenize=False, add_generation_prompt=True)
# tokenize the text
input_tokens = tokenizer(chat, return_tensors="pt").to(device)
# generate output tokens
output = model.generate(**input_tokens, 
                        max_new_tokens=100)
# decode output tokens into text
output = tokenizer.batch_decode(output)
# print output
print(output)
` + "```" + `

**Evaluation Results:**
<table>
  <caption><b>HuggingFace Open LLM Leaderboard V1</b></caption>
<thead>
  <tr>
    <th style="text-align:left; background-color: #001d6c; color: white;">Models</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">ARC-Challenge</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">Hellaswag</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">MMLU</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">TruthfulQA</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">Winogrande</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">GSM8K</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">Avg</th>
  </tr></thead>
  <tbody>
  <tr>
    <td style="text-align:left; background-color: #DAE8FF; color: black;">Granite-3.1-8B-Instruct</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">62.62</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">84.48</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">65.34</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">66.23</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">75.37</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">73.84</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">71.31</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: #2D2D2D;">Granite-3.1-2B-Instruct</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">54.61</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">75.14</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">55.31</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">59.42</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">67.48</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">52.76</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">60.79</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: #2D2D2D;">Granite-3.1-3B-A800M-Instruct</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">50.42</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">73.01</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">52.19</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">49.71</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">64.87</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">48.97</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">56.53</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: #2D2D2D;">Granite-3.1-1B-A400M-Instruct</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">42.66</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">65.97</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">26.13</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">46.77</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">62.35</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">33.88</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">46.29</td>
  </tr>
</tbody></table>

<table>
  <caption><b>HuggingFace Open LLM Leaderboard V2</b></caption>
<thead>
  <tr>
    <th style="text-align:left; background-color: #001d6c; color: white;">Models</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">IFEval</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">BBH</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">MATH Lvl 5</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">GPQA</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">MUSR</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">MMLU-Pro</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">Avg</th>
  </tr></thead>
  <tbody>
  <tr>
    <td style="text-align:left; background-color: #DAE8FF; color: black;">Granite-3.1-8B-Instruct</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">72.08</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">34.09</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">21.68</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">8.28</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">19.01</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">28.19</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">30.55</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: #2D2D2D;">Granite-3.1-2B-Instruct</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">62.86</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">21.82</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">11.33</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">5.26</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">4.87</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">20.21</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">21.06</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: #2D2D2D;">Granite-3.1-3B-A800M-Instruct</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">55.16</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">16.69</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">10.35</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">5.15</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">2.51</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">12.75</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">17.1</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: #2D2D2D;">Granite-3.1-1B-A400M-Instruct</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">46.86</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">6.18</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">4.08</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">0</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">0.78</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">2.41</td>
    <td style="text-align:center; background-color: #FFFFFF; color: #2D2D2D;">10.05</td>
  </tr>
</tbody></table>

**Model Architecture:**
Granite-3.1-8B-Instruct is based on a decoder-only dense transformer architecture. Core components of this architecture are: GQA and RoPE, MLP with SwiGLU, RMSNorm, and shared input/output embeddings.

<table>
<thead>
  <tr>
    <th style="text-align:left; background-color: #001d6c; color: white;">Model</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">2B Dense</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">8B Dense</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">1B MoE</th>
    <th style="text-align:center; background-color: #001d6c; color: white;">3B MoE</th>
  </tr></thead>
<tbody>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">Embedding size</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">2048</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">4096</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">1024</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">1536</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">Number of layers</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">40</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">40</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">24</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">32</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">Attention head size</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">64</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">128</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">64</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">64</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">Number of attention heads</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">32</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">32</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">16</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">24</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">Number of KV heads</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">8</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">8</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">8</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">8</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">MLP hidden size</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">8192</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">12800</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">512</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">512</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">MLP activation</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">SwiGLU</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">SwiGLU</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">SwiGLU</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">SwiGLU</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">Number of experts</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">—</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">—</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">32</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">40</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">MoE TopK</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">—</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">—</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">8</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">8</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">Initialization std</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">0.1</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">0.1</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">0.1</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">0.1</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">Sequence length</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">128K</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">128K</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">128K</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">128K</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;">Position embedding</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">RoPE</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">RoPE</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">RoPE</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">RoPE</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;"># Parameters</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">2.5B</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">8.1B</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">1.3B</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">3.3B</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;"># Active parameters</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">2.5B</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">8.1B</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">400M</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">800M</td>
  </tr>
  <tr>
    <td style="text-align:left; background-color: #FFFFFF; color: black;"># Training tokens</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">12T</td>
    <td style="text-align:center; background-color: #DAE8FF; color: black;">12T</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">10T</td>
    <td style="text-align:center; background-color: #FFFFFF; color: black;">10T</td>
  </tr>
</tbody></table>

**Training Data:** 
Overall, our SFT data is largely comprised of three key sources: (1) publicly available datasets with permissive license, (2) internal synthetic data targeting specific capabilities including long-context tasks, and (3) very small amounts of human-curated data. A detailed attribution of datasets can be found in the [Granite 3.0 Technical Report](https://github.com/provider1-granite/granite-3.0-language-models/blob/main/paper.pdf), [Granite 3.1 Technical Report (coming soon)](https://huggingface.co/collections/provider1-granite/granite-31-language-models-6751dbbf2f3389bec5c6f02d), and [Accompanying Author List](https://github.com/provider1-granite/granite-3.0-language-models/blob/main/author-ack.pdf).

**Infrastructure:**
We train Granite 3.1 Language Models using provider1's super computing cluster, Blue Vela, which is outfitted with NVIDIA H100 GPUs. This cluster provides a scalable and efficient infrastructure for training our models over thousands of GPUs.

**Ethical Considerations and Limitations:** 
Granite 3.1 Instruct Models are primarily finetuned using instruction-response pairs mostly in English, but also multilingual data covering eleven languages. Although this model can handle multilingual dialog use cases, its performance might not be similar to English tasks. In such case, introducing a small number of examples (few-shot) can help the model in generating more accurate outputs. While this model has been aligned by keeping safety in consideration, the model may in some cases produce inaccurate, biased, or unsafe responses to user prompts. So we urge the community to use this model with proper safety testing and tuning tailored for their specific tasks.

**Resources**
- ⭐️ Learn about the latest updates with Granite: https://www.provider1.com/granite
- 📄 Get started with tutorials, best practices, and prompt engineering advice: https://www.provider1.com/granite/docs/
- 💡 Learn about the latest Granite learning resources: https://provider1.biz/granite-learning-resources

<!-- ## Citation
` + "```" + `
@misc{granite-models,
  author = {author 1, author2, ...},
  title = {},
  journal = {},
  volume = {},
  year = {2024},
  url = {https://arxiv.org/abs/0000.00000},
}
  ` + "```" + ` -->`),
		SourceId:                 stringToPointer("sample-source"),
		LibraryName:              stringToPointer("transformers"),
		CreateTimeSinceEpoch:     stringToPointer("1693526400000"),
		LastUpdateTimeSinceEpoch: stringToPointer("1704067200000"),
	}

	sampleModel2 := models.CatalogModel{
		Name:             "repo1/granite-7b-instruct",
		Description:      stringToPointer("Granite 7B instruction-tuned model for enterprise applications"),
		Provider:         stringToPointer("Red Hat"),
		Tasks:            []string{"text-generation", "image-text-to-text"},
		License:          stringToPointer("apache-2.0"),
		Maturity:         stringToPointer("Generally Available"),
		Language:         []string{"en"},
		SourceId:         stringToPointer("sample-source"),
		CustomProperties: catalogCustomProperties(),
		Logo:             stringToPointer("data:image/svg+xml;base64,PHN2ZyBpZD0iTGF5ZXJfMSIgZGF0YS1uYW1lPSJMYXllciAxIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxOTIgMTQ1Ij48ZGVmcz48c3R5bGU+LmNscy0xe2ZpbGw6I2UwMDt9PC9zdHlsZT48L2RlZnM+PHRpdGxlPlJlZEhhdC1Mb2dvLUhhdC1Db2xvcjwvdGl0bGU+PHBhdGggZD0iTTE1Ny43Nyw2Mi42MWExNCwxNCwwLDAsMSwuMzEsMy40MmMwLDE0Ljg4LTE4LjEsMTcuNDYtMzAuNjEsMTcuNDZDNzguODMsODMuNDksNDIuNTMsNTMuMjYsNDIuNTMsNDRhNi40Myw2LjQzLDAsMCwxLC4yMi0xLjk0bC0zLjY2LDkuMDZhMTguNDUsMTguNDUsMCwwLDAtMS41MSw3LjMzYzAsMTguMTEsNDEsNDUuNDgsODcuNzQsNDUuNDgsMjAuNjksMCwzNi40My03Ljc2LDM2LjQzLTIxLjc3LDAtMS4wOCwwLTEuOTQtMS43My0xMC4xM1oiLz48cGF0aCBjbGFzcz0iY2xzLTEiIGQ9Ik0xMjcuNDcsODMuNDljMTIuNTEsMCwzMC42MS0yLjU4LDMwLjYxLTE3LjQ2YTE0LDE0LDAsMCwwLS4zMS0zLjQybC03LjQ1LTMyLjM2Yy0xLjcyLTcuMTItMy4yMy0xMC4zNS0xNS43My0xNi42QzEyNC44OSw4LjY5LDEwMy43Ni41LDk3LjUxLjUsOTEuNjkuNSw5MCw4LDgzLjA2LDhjLTYuNjgsMC0xMS42NC01LjYtMTcuODktNS42LTYsMC05LjkxLDQuMDktMTIuOTMsMTIuNSwwLDAtOC40MSwyMy43Mi05LjQ5LDI3LjE2QTYuNDMsNi40MywwLDAsMCw0Mi41Myw0NGMwLDkuMjIsMzYuMywzOS40NSw4NC45NCwzOS40NU0xNjAsNzIuMDdjMS43Myw4LjE5LDEuNzMsOS4wNSwxLjczLDEwLjEzLDAsMTQtMTUuNzQsMjEuNzctMzYuNDMsMjEuNzdDNzguNTQsMTA0LDM3LjU4LDc2LjYsMzcuNTgsNTguNDlhMTguNDUsMTguNDUsMCwwLDEsMS41MS03LjMzQzIyLjI3LDUyLC41LDU1LC41LDc0LjIyYzAsMzEuNDgsNzQuNTksNzAuMjgsMTMzLjY1LDcwLjI4LDQ1LjI4LDAsNTYuNy0yMC40OCw1Ni43LTM2LjY1LDAtMTIuNzItMTEtMjcuMTYtMzAuODMtMzUuNzgiLz48L3N2Zz4="),
	}

	sampleModel3 := models.CatalogModel{
		Name:             "repo1/granite-3b-code-base",
		Description:      stringToPointer("Granite 3B code generation model for programming tasks"),
		Provider:         stringToPointer("IBM"),
		Tasks:            []string{"audio-to-text", "text-to-text", "video-to-text"},
		License:          stringToPointer("mit"),
		Maturity:         stringToPointer("Generally Available"),
		Language:         []string{"en"},
		SourceId:         stringToPointer("sample-source"),
		CustomProperties: catalogCustomProperties(),
		Logo:             stringToPointer("data:image/svg+xml;base64,PHN2ZyBpZD0iTGF5ZXJfMSIgZGF0YS1uYW1lPSJMYXllciAxIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxOTIgMTQ1Ij48ZGVmcz48c3R5bGU+LmNscy0xe2ZpbGw6I2UwMDt9PC9zdHlsZT48L2RlZnM+PHRpdGxlPlJlZEhhdC1Mb2dvLUhhdC1Db2xvcjwvdGl0bGU+PHBhdGggZD0iTTE1Ny43Nyw2Mi42MWExNCwxNCwwLDAsMSwuMzEsMy40MmMwLDE0Ljg4LTE4LjEsMTcuNDYtMzAuNjEsMTcuNDZDNzguODMsODMuNDksNDIuNTMsNTMuMjYsNDIuNTMsNDRhNi40Myw2LjQzLDAsMCwxLC4yMi0xLjk0bC0zLjY2LDkuMDZhMTguNDUsMTguNDUsMCwwLDAtMS41MSw3LjMzYzAsMTguMTEsNDEsNDUuNDgsODcuNzQsNDUuNDgsMjAuNjksMCwzNi40My03Ljc2LDM2LjQzLTIxLjc3LDAtMS4wOCwwLTEuOTQtMS43My0xMC4xM1oiLz48cGF0aCBjbGFzcz0iY2xzLTEiIGQ9Ik0xMjcuNDcsODMuNDljMTIuNTEsMCwzMC42MS0yLjU4LDMwLjYxLTE3LjQ2YTE0LDE0LDAsMCwwLS4zMS0zLjQybC03LjQ1LTMyLjM2Yy0xLjcyLTcuMTItMy4yMy0xMC4zNS0xNS43My0xNi42QzEyNC44OSw4LjY5LDEwMy43Ni41LDk3LjUxLjUsOTEuNjkuNSw5MCw4LDgzLjA2LDhjLTYuNjgsMC0xMS42NC01LjYtMTcuODktNS42LTYsMC05LjkxLDQuMDktMTIuOTMsMTIuNSwwLDAtOC40MSwyMy43Mi05LjQ5LDI3LjE2QTYuNDMsNi40MywwLDAsMCw0Mi41Myw0NGMwLDkuMjIsMzYuMywzOS40NSw4NC45NCwzOS40NU0xNjAsNzIuMDdjMS43Myw4LjE5LDEuNzMsOS4wNSwxLjczLDEwLjEzLDAsMTQtMTUuNzQsMjEuNzctMzYuNDMsMjEuNzdDNzguNTQsMTA0LDM3LjU4LDc2LjYsMzcuNTgsNTguNDlhMTguNDUsMTguNDUsMCwwLDEsMS41MS03LjMzQzIyLjI3LDUyLC41LDU1LC41LDc0LjIyYzAsMzEuNDgsNzQuNTksNzAuMjgsMTMzLjY1LDcwLjI4LDQ1LjI4LDAsNTYuNy0yMC40OCw1Ni43LTM2LjY1LDAtMTIuNzItMTEtMjcuMTYtMzAuODMtMzUuNzgiLz48L3N2Zz4="),
	}

	huggingFaceModel1 := models.CatalogModel{
		Name:        "provider2/bert-base-uncased",
		Description: stringToPointer("BERT base model (uncased) - Pretrained model on English language"),
		Provider:    stringToPointer("Google"),
		Tasks:       []string{"audio-to-text", "text-to-text"},
		License:     stringToPointer("apache-2.0"),
		Maturity:    stringToPointer("Generally Available"),
		Language:    []string{"en"},
		SourceId:    stringToPointer("huggingface"),
		LibraryName: stringToPointer("transformers"),
	}

	huggingFaceModel2 := models.CatalogModel{
		Name:        "provider3/gpt2",
		Description: stringToPointer("GPT-2 is a transformers model pretrained on a very large corpus of English data"),
		Provider:    stringToPointer("provider3"),
		Tasks:       []string{"video-to-text"},
		License:     stringToPointer("mit"),
		Maturity:    stringToPointer("Generally Available"),
		Language:    []string{"en"},
		SourceId:    stringToPointer("huggingface"),
		LibraryName: stringToPointer("transformers"),
	}

	huggingFaceModel3 := models.CatalogModel{
		Name:        "huggingface/distilbert-base-uncased",
		Description: stringToPointer("DistilBERT base model (uncased) - A smaller, faster version of BERT"),
		Provider:    stringToPointer("Hugging Face"),
		Tasks:       []string{"fill-mask", "text-classification"},
		License:     stringToPointer("apache-2.0"),
		Maturity:    stringToPointer("Generally Available"),
		Language:    []string{"en"},
		SourceId:    stringToPointer("huggingface"),
		LibraryName: stringToPointer("transformers"),
	}

	otherModel1 := models.CatalogModel{
		Name:        "adminModel2/admin-model-1",
		Description: stringToPointer("sample description"),
		Provider:    stringToPointer("Admin model 1"),
		Tasks:       []string{"code-generation", "instruction-following"},
		License:     stringToPointer("apache-2.0"),
		Maturity:    stringToPointer("Generally Available"),
		Language:    []string{"en"},
		SourceId:    stringToPointer("adminModel2"),
	}

	otherModel2 := models.CatalogModel{
		Name:        "adminModel1/admin-model-2",
		Description: stringToPointer("sample description"),
		Provider:    stringToPointer("Admin model 2"),
		Tasks:       []string{"text-generation", "conversational"},
		License:     stringToPointer("apache-2.0"),
		Maturity:    stringToPointer("Generally Available"),
		Language:    []string{"en"},
		SourceId:    stringToPointer("adminModel1"),
	}

	// added this to test the load more models button
	var additionalRepo1Models []models.CatalogModel
	for i := 1; i <= 20; i++ {
		model := models.CatalogModel{
			Name:                     fmt.Sprintf("repo1/granite-model-%d", i),
			Description:              stringToPointer("Granite-8B-Code-Instruct is a 8B parameter model fine tuned from\nGranite-8B-Code-Base on a combination of permissively licensed instruction\ndata to enhance instruction following capabilities including logical\nreasoning and problem-solving skills."),
			Provider:                 stringToPointer("provider1"),
			Tasks:                    []string{"text-generation"},
			License:                  stringToPointer("apache-2.0"),
			LicenseLink:              stringToPointer("https://www.apache.org/licenses/LICENSE-2.0.txt"),
			Maturity:                 stringToPointer("Technology preview"),
			Language:                 []string{"ar", "cs", "de", "en", "es", "fr", "it", "ja", "ko", "nl", "pt", "zh"},
			Logo:                     stringToPointer("data:image/svg+xml;base64,PHN2ZyBpZD0iTGF5ZXJfMSIgZGF0YS1uYW1lPSJMYXllciAxIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxOTIgMTQ1Ij48ZGVmcz48c3R5bGU+LmNscy0xe2ZpbGw6I2UwMDt9PC9zdHlsZT48L2RlZnM+PHRpdGxlPlJlZEhhdC1Mb2dvLUhhdC1Db2xvcjwvdGl0bGU+PHBhdGggZD0iTTE1Ny43Nyw2Mi42MWExNCwxNCwwLDAsMSwuMzEsMy40MmMwLDE0Ljg4LTE4LjEsMTcuNDYtMzAuNjEsMTcuNDZDNzguODMsODMuNDksNDIuNTMsNTMuMjYsNDIuNTMsNDRhNi40Myw2LjQzLDAsMCwxLC4yMi0xLjk0bC0zLjY2LDkuMDZhMTguNDUsMTguNDUsMCwwLDAtMS41MSw3LjMzYzAsMTguMTEsNDEsNDUuNDgsODcuNzQsNDUuNDgsMjAuNjksMCwzNi40My03Ljc2LDM2LjQzLTIxLjc3LDAtMS4wOCwwLTEuOTQtMS43My0xMC4xM1oiLz48cGF0aCBjbGFzcz0iY2xzLTEiIGQ9Ik0xMjcuNDcsODMuNDljMTIuNTEsMCwzMC42MS0yLjU4LDMwLjYxLTE3LjQ2YTE0LDE0LDAsMCwwLS4zMS0zLjQybC03LjQ1LTMyLjM2Yy0xLjcyLTcuMTItMy4yMy0xMC4zNS0xNS43My0xNi42QzEyNC44OSw4LjY5LDEwMy43Ni41LDk3LjUxLjUsOTEuNjkuNSw5MCw4LDgzLjA2LDhjLTYuNjgsMC0xMS42NC01LjYtMTcuODktNS42LTYsMC05LjkxLDQuMDktMTIuOTMsMTIuNSwwLDAtOC40MSwyMy43Mi05LjQ5LDI3LjE2QTYuNDMsNi40MywwLDAsMCw0Mi41Myw0NGMwLDkuMjIsMzYuMywzOS40NSw4NC45NCwzOS40NU0xNjAsNzIuMDdjMS43Myw4LjE5LDEuNzMsOS4wNSwxLjczLDEwLjEzLDAsMTQtMTUuNzQsMjEuNzctMzYuNDMsMjEuNzdDNzguNTQsMTA0LDM3LjU4LDc2LjYsMzcuNTgsNTguNDlhMTguNDUsMTguNDUsMCwwLDEsMS41MS03LjMzQzIyLjI3LDUyLC41LDU1LC41LDc0LjIyYzAsMzEuNDgsNzQuNTksNzAuMjgsMTMzLjY1LDcwLjI4LDQ1LjI4LDAsNTYuNy0yMC40OCw1Ni43LTM2LjY1LDAtMTIuNzItMTEtMjcuMTYtMzAuODMtMzUuNzgiLz48L3N2Zz4="),
			SourceId:                 stringToPointer("sample-source"),
			LibraryName:              stringToPointer("transformers"),
			CreateTimeSinceEpoch:     stringToPointer("1693526400000"),
			LastUpdateTimeSinceEpoch: stringToPointer("1704067200000"),
			CustomProperties:         catalogCustomProperties(),
		}
		additionalRepo1Models = append(additionalRepo1Models, model)
	}

	allModels := []models.CatalogModel{
		sampleModel1, sampleModel2, sampleModel3,
		huggingFaceModel1, huggingFaceModel2, huggingFaceModel3,
		otherModel1, otherModel2,
	}
	allModels = append(allModels, additionalRepo1Models...)

	return allModels
}

func GetCatalogModelListMock() models.CatalogModelList {
	allModels := GetCatalogModelMocks()

	return models.CatalogModelList{
		Items:         allModels,
		Size:          int32(len(allModels)),
		PageSize:      int32(10),
		NextPageToken: "10",
	}
}

func GetCatalogSourceMocks() []models.CatalogSource {
	enabled := true
	disabledBool := false

	// Status examples (matching OpenAPI spec)
	availableStatus := "available"
	errorStatus := "error"
	disabledStatus := "disabled"

	invalidCredentialError := "The provided API key is invalid or has expired. Please update your credentials."
	invalidOrgError := "The specified organization 'invalid-org' does not exist or you don't have access to it. Please verify the organization name and ensure you have the necessary permissions to access models from this organization."

	return []models.CatalogSource{
		{
			Id:      "sample-source",
			Name:    "Sample mocked source",
			Enabled: &enabled,
			Labels:  []string{"Sample category 1", "Sample category 2", "Sample category"},
			Status:  &availableStatus,
		},
		{
			Id:     "huggingface",
			Name:   "Hugging Face",
			Labels: []string{"Sample category 2", "Sample category"},
			// Status is nil - represents "Starting" state (no status yet)
			Status: nil,
		},
		{
			Id:      "adminModel1",
			Name:    "Admin model 1",
			Enabled: &enabled,
			Labels:  []string{},
			Status:  &errorStatus,
			Error:   &invalidCredentialError,
		},
		{
			Id:      "adminModel2",
			Name:    "Admin model 2",
			Enabled: &enabled,
			Labels:  []string{"Sample category 1"},
			Status:  &errorStatus,
			Error:   &invalidOrgError,
		},
		{
			Id:     "dora",
			Name:   "Dora source",
			Labels: []string{},
			Status: &availableStatus,
		},
		{
			Id:      "adminModel3",
			Name:    "Admin model 3",
			Enabled: &disabledBool,
			Labels:  []string{},
			Status:  &disabledStatus,
		},
	}
}

func GetCatalogSourceListMock() models.CatalogSourceList {
	allSources := GetCatalogSourceMocks()

	return models.CatalogSourceList{
		Items:         allSources,
		Size:          int32(len(allSources)),
		PageSize:      int32(10),
		NextPageToken: "",
	}
}

func GetCatalogModelArtifactMock() []models.CatalogArtifact {
	return []models.CatalogArtifact{
		{
			ArtifactType:         "model-artifact",
			Uri:                  stringToPointer("oci://registry.sample.io/repo1/modelcar-granite-7b-starter:1.4.0"),
			CreateTimeSinceEpoch: stringToPointer("1693526400000"),

			LastUpdateTimeSinceEpoch: stringToPointer("1704067200000"),
			CustomProperties:         newCustomProperties(),
		},
	}
}

func performanceMetricsCustomProperties(customProperties map[string]openapi.MetadataValue) *map[string]openapi.MetadataValue {
	result := map[string]openapi.MetadataValue{
		"config_id": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "0055d94f6547237dgf324238",
				MetadataType: "MetadataStringValue",
			},
		},
		"totalRPS": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  30,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"ttft_mean": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  35.48818160947744,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"ttft_p90": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  51.55777931213379,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"ttft_p95": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  61.26761436462402,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"ttft_p99": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  72.95823097229004,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"e2e_mean": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  1994.480013381083,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"e2e_p90": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  2644.604682922363,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"e2e_p95": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  2813.79246711731,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"e2e_p99": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  3117.565155029297,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"tps_mean": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  1785.325259154939,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"tps_p90": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  3318.278481012658,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"tps_p95": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  4934.475294117647,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"tps_p99": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  11781.75280898876,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"itl_mean": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  7.685877762379151,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"itl_p90": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  7.778410935521725,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"itl_p95": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  7.812754891135476,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"itl_p99": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  7.9100158577958,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"requests_per_second": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  7,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"max_input_tokens": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  1024,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"max_output_tokens": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  1,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"hardware_type": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "H100",
				MetadataType: "MetadataStringValue",
			},
		},
		"hardware_count": {
			MetadataIntValue: &openapi.MetadataIntValue{
				IntValue:     "2",
				MetadataType: "MetadataIntValue",
			},
		},
		"framework": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "vllm",
				MetadataType: "MetadataStringValue",
			},
		},
		"framework_version": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "v0.1.1",
				MetadataType: "MetadataStringValue",
			},
		},
		"docker_image": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "vllm/vllm-openai:v0.1.1",
				MetadataType: "MetadataStringValue",
			},
		},
		"entrypoint": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue: "\npython3\n",
			},
		},
		"inserted_at": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "2025-05-07T00:00:00.000Z",
				MetadataType: "MetadataStringValue",
			},
		},
		"created_at": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "2025-05-07T00:00:00.000Z",
				MetadataType: "MetadataStringValue",
			},
		},
		"updated_at": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue:  "2025-05-14T12:08:25.402Z",
				MetadataType: "MetadataStringValue",
			},
		},
		"mean_input_tokens": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  511.5445458496306,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"mean_output_tokens": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  255.8678835289005,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"model_hf_repo_name": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue: "provider1-granite/granite-3.1-8b-instruct",
			},
		},
	}
	for key, value := range customProperties {
		result[key] = value
	}
	return &result
}

func GetCatalogPerformanceMetricsArtifactMock(itemCount int32) []models.CatalogArtifact {
	artifacts := []models.CatalogArtifact{
		{
			ArtifactType:             *stringToPointer("metrics-artifact"),
			MetricsType:              stringToPointer("performance-metrics"),
			CreateTimeSinceEpoch:     stringToPointer("1693526400000"),
			LastUpdateTimeSinceEpoch: stringToPointer("1704067200000"),
			CustomProperties: performanceMetricsCustomProperties(map[string]openapi.MetadataValue{
				"use_case": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "chatbot",
						MetadataType: "MetadataStringValue",
					},
				},
			}),
		},
		{
			ArtifactType:             *stringToPointer("metrics-artifact"),
			MetricsType:              stringToPointer("performance-metrics"),
			CreateTimeSinceEpoch:     stringToPointer("1693526400000"),
			LastUpdateTimeSinceEpoch: stringToPointer("1704067200000"),
			CustomProperties: performanceMetricsCustomProperties(map[string]openapi.MetadataValue{
				"hardware_type": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "RTX 4090",
						MetadataType: "MetadataStringValue",
					},
				},
				"hardware_count": {
					MetadataIntValue: &openapi.MetadataIntValue{
						IntValue:     "33",
						MetadataType: "MetadataIntValue",
					},
				},
				"requests_per_second": {
					MetadataDoubleValue: &openapi.MetadataDoubleValue{
						DoubleValue:  10,
						MetadataType: "MetadataDoubleValue",
					},
				},
				"ttft_mean": {
					MetadataDoubleValue: &openapi.MetadataDoubleValue{
						DoubleValue:  67.14892749816,
						MetadataType: "MetadataDoubleValue",
					},
				},
				"use_case": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "rag",
						MetadataType: "MetadataStringValue",
					},
				},
			}),
		},
		{
			ArtifactType:             *stringToPointer("metrics-artifact"),
			MetricsType:              stringToPointer("performance-metrics"),
			CreateTimeSinceEpoch:     stringToPointer("1693526400000"),
			LastUpdateTimeSinceEpoch: stringToPointer("1704067200000"),
			CustomProperties: performanceMetricsCustomProperties(map[string]openapi.MetadataValue{
				"hardware_type": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "A100",
						MetadataType: "MetadataStringValue",
					},
				},
				"hardware_count": {
					MetadataIntValue: &openapi.MetadataIntValue{
						IntValue:     "40",
						MetadataType: "MetadataIntValue",
					},
				},
				"requests_per_second": {
					MetadataDoubleValue: &openapi.MetadataDoubleValue{
						DoubleValue:  15,
						MetadataType: "MetadataDoubleValue",
					},
				},
				"ttft_mean": {
					MetadataDoubleValue: &openapi.MetadataDoubleValue{
						DoubleValue:  42.123791232,
						MetadataType: "MetadataDoubleValue",
					},
				},
				"use_case": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "code_fixing",
						MetadataType: "MetadataStringValue",
					},
				},
			}),
		},
		{
			ArtifactType:             *stringToPointer("metrics-artifact"),
			MetricsType:              stringToPointer("performance-metrics"),
			CreateTimeSinceEpoch:     stringToPointer("1693526400000"),
			LastUpdateTimeSinceEpoch: stringToPointer("1704067200000"),
			CustomProperties: performanceMetricsCustomProperties(map[string]openapi.MetadataValue{
				"hardware_type": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "A100",
						MetadataType: "MetadataStringValue",
					},
				},
				"hardware_count": {
					MetadataIntValue: &openapi.MetadataIntValue{
						IntValue:     "8",
						MetadataType: "MetadataIntValue",
					},
				},
				"requests_per_second": {
					MetadataDoubleValue: &openapi.MetadataDoubleValue{
						DoubleValue:  25,
						MetadataType: "MetadataDoubleValue",
					},
				},
				"ttft_mean": {
					MetadataDoubleValue: &openapi.MetadataDoubleValue{
						DoubleValue:  28.5,
						MetadataType: "MetadataDoubleValue",
					},
				},
				"use_case": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "long_rag",
						MetadataType: "MetadataStringValue",
					},
				},
			}),
		},
	}
	artifacts = artifacts[:itemCount]
	return artifacts
}

func accuracyMetricsCustomProperties() *map[string]openapi.MetadataValue {
	result := map[string]openapi.MetadataValue{
		"overall_average": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  0.584329,
				MetadataType: "MetadataDoubleValue",
			},
		},
		"arc_v1": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue:  0.674673,
				MetadataType: "MetadataDoubleValue",
			},
		},
	}
	return &result
}

func GetCatalogAccuracyMetricsArtifactMock() []models.CatalogArtifact {
	return []models.CatalogArtifact{
		{
			ArtifactType:             *stringToPointer("metrics-artifact"),
			MetricsType:              stringToPointer("accuracy-metrics"),
			CreateTimeSinceEpoch:     stringToPointer("1693526400000"),
			LastUpdateTimeSinceEpoch: stringToPointer("1704067200000"),
			CustomProperties:         accuracyMetricsCustomProperties(),
		},
	}
}
func GetModelArtifactListMockWithItems(items []models.CatalogArtifact, pageSize int32) models.CatalogModelArtifactList {
	return models.CatalogModelArtifactList{
		Items:         items,
		Size:          int32(len(items)),
		PageSize:      pageSize,
		NextPageToken: "",
	}
}

func GetCatalogModelArtifactListMock() models.CatalogModelArtifactList {
	allArtifactMock := GetCatalogModelArtifactMock()
	return GetModelArtifactListMockWithItems(allArtifactMock, 10)
}

func GetCatalogPerformanceMetricsArtifactListMock(itemCount int32) models.CatalogModelArtifactList {
	allArtifactMock := GetCatalogPerformanceMetricsArtifactMock(itemCount)
	return GetModelArtifactListMockWithItems(allArtifactMock, 10)
}

func GetCatalogAccuracyMetricsArtifactListMock() models.CatalogModelArtifactList {
	allArtifactMock := GetCatalogAccuracyMetricsArtifactMock()
	return GetModelArtifactListMockWithItems(allArtifactMock, 10)
}

const (
	FilterOptionTypeString = "string"
	FilterOptionTypeNumber = "number"
)

func float32Ptr(i float32) *float32 {
	return &i
}

func GetFilterOptionMocks() map[string]models.FilterOption {
	filterOptions := make(map[string]models.FilterOption)

	filterOptions["provider"] = models.FilterOption{
		Type: FilterOptionTypeString,
		Values: []interface{}{
			"Red Hat", "IBM", "Google",
		},
	}

	filterOptions["license"] = models.FilterOption{
		Type: FilterOptionTypeString,
		Values: []interface{}{
			"apache-2.0",
			"mit",
		},
	}

	filterOptions["tasks"] = models.FilterOption{
		Type: FilterOptionTypeString,
		Values: []interface{}{
			"audio-to-text", "image-to-text", "image-text-to-text", "text-generation", "text-to-text", "video-to-text",
		},
	}

	// String type filter for programming languages supported
	filterOptions["language"] = models.FilterOption{
		Type: FilterOptionTypeString,
		Values: []interface{}{
			"ar", "cs", "de", "en", "es", "fr", "it", "ja", "ko", "nl", "pt", "zh",
		},
	}

	// String type filter for use cases
	filterOptions["use_case"] = models.FilterOption{
		Type: FilterOptionTypeString,
		Values: []interface{}{
			"chatbot", "code_fixing", "long_rag", "rag",
		},
	}

	filterOptions["ttft_mean"] = models.FilterOption{
		Type: FilterOptionTypeNumber,
		Range: &models.FilterRange{
			Min: float32Ptr(1),
			Max: float32Ptr(100),
		},
	}

	return filterOptions
}

func GetFilterOptionsListMock() models.FilterOptionsList {
	filterOptions := GetFilterOptionMocks()

	return models.FilterOptionsList{
		Filters: &filterOptions,
	}
}

func CreateSampleCatalogSource(id string, name string, catalogType string, enabled bool) models.CatalogSourceConfig {
	defaultCatalog := id == "sample-source"

	sourceConfig := models.CatalogSourceConfig{
		Name:      name,
		Id:        id,
		Type:      catalogType,
		Enabled:   &enabled,
		Labels:    []string{"source-1"},
		IsDefault: &defaultCatalog,
	}

	if !defaultCatalog {
		sourceConfig.IncludedModels = []string{"rhelai1/modelcar-granite-7b-starter"}
		sourceConfig.ExcludedModels = []string{"model-a:1.0", "model-b:*"}
	}

	switch catalogType {
	case "yaml":
		sourceConfig.Yaml = stringToPointer("models:\n  - name: model1")
	case "hf":
		// Use different organizations for the failed sources
		if id == "adminModel2" {
			sourceConfig.AllowedOrganization = stringToPointer("invalid-org")
		} else {
			sourceConfig.AllowedOrganization = stringToPointer("org1")
		}
		sourceConfig.ApiKey = stringToPointer("apikey")
	}

	return sourceConfig
}

func BoolPtr(b bool) *bool {
	return &b
}

func GetModelsWithInclusionStatusListMocks() []models.CatalogSourcePreviewModel {
	return []models.CatalogSourcePreviewModel{
		{
			Name:     "sample-source/granite",
			Included: true,
		},
		{
			Name:     "sample-source/model-1",
			Included: true,
		},
		{
			Name:     "sample-source/model-2",
			Included: true,
		},
		{
			Name:     "sample-source/model-3",
			Included: true,
		},
		{
			Name:     "sample-source/model-4",
			Included: true,
		},
		{
			Name:     "sample-source/model-5",
			Included: true,
		},
		{
			Name:     "sample-source/model-6",
			Included: false,
		},
		{
			Name:     "adminModel1/model-1",
			Included: true,
		},
		{
			Name:     "adminModel1/model-2",
			Included: true,
		},
		{
			Name:     "adminModel1/model-3",
			Included: true,
		},
		{
			Name:     "adminModel1/model-4",
			Included: true,
		},
		{
			Name:     "adminModel1/model-5",
			Included: true,
		},
		{
			Name:     "adminModel1/model-6",
			Included: true,
		},
		{
			Name:     "adminModel1/model-7",
			Included: true,
		},
		{
			Name:     "adminModel1/model-8",
			Included: true,
		},
		{
			Name:     "adminModel1/model-9",
			Included: true,
		},
		{
			Name:     "adminModel1/model-10",
			Included: false,
		},
		{
			Name:     "adminModel1/model-11",
			Included: false,
		},
		{
			Name:     "adminModel1/model-12",
			Included: false,
		},
		{
			Name:     "adminModel1/model-13",
			Included: false,
		},
	}
}

func GetCatalogSourcePreviewSummaryMock() models.CatalogSourcePreviewSummary {
	return models.CatalogSourcePreviewSummary{
		TotalModels:    20,
		IncludedModels: 15,
		ExcludedModels: 5,
	}
}

func CreateCatalogSourcePreviewMock() models.CatalogSourcePreviewResult {
	catalogModelPreview := GetModelsWithInclusionStatusListMocks()
	catalogSourcePreviewSummary := GetCatalogSourcePreviewSummaryMock()

	return models.CatalogSourcePreviewResult{
		Items:         catalogModelPreview,
		Summary:       catalogSourcePreviewSummary,
		NextPageToken: "",
		PageSize:      int32(10),
		Size:          int32(len(catalogModelPreview)),
	}
}
