package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Swagger struct {
	SwaggerVersion string `json:"swagger"`
	Info           struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Version     string `json:"version"`
	} `json:"info"`
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
	Host     string   `json:"host"`
	Consumes []string `json:"consumes"`
	Produces []string `json:"produces"`
	Paths    map[string]map[string]struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
		OperationId string `json:"operationId"`
		Responses   map[string]struct {
			Description string `json:"description"`
			Schema      struct {
				Ref string `json:"$ref"`
			} `json:"schema"`
		} `json:"responses"`
		Tags []string `json:"tags"`
	} `json:"paths"`
	Definitions map[string]struct {
		Type       string `json:"type"`
		Properties map[string]struct {
			Type   string `json:"type"`
			Format string `json:"format,omitempty"`
			Items  struct {
				Ref string `json:"$ref"`
			} `json:"items,omitempty"`
		} `json:"properties,omitempty"`
	} `json:"definitions"`
}

type SwaggerYAML struct {
	SwaggerVersion string `yaml:"swagger"`
	Info           struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Version     string `yaml:"version"`
	} `yaml:"info"`
	Tags []struct {
		Name string `yaml:"name"`
	} `yaml:"tags"`
	BasePath string   `yaml:"basePath"`
	Consumes []string `yaml:"consumes"`
	Produces []string `yaml:"produces"`
	Paths    map[string]map[string]struct {
		Summary     string `yaml:"summary"`
		Description string `yaml:"description"`
		OperationId string `yaml:"operationId"`
		Responses   map[string]struct {
			Description string `yaml:"description"`
			Schema      struct {
				Ref string `yaml:"$ref"`
			} `yaml:"schema"`
		} `yaml:"responses"`
		Tags                         []string `yaml:"tags"`
		XAmazonApigatewayIntegration struct {
			ConnectionId string `yaml:"connectionId"`
			HttpMethod   string `yaml:"httpMethod"`
			Uri          string `yaml:"uri"`
			Responses    map[string]struct {
				StatusCode string `yaml:"statusCode"`
			} `yaml:"responses"`
			PassThroughBehavior string            `yaml:"passthroughBehavior"`
			ConnectionType      string            `yaml:"connectionType"`
			RequestParameters   map[string]string `yaml:"requestParameters"`
			Type                string            `yaml:"type"`
		} `yaml:"x-amazon-apigateway-integration"`
	} `yaml:"paths"`
	Definitions map[string]struct {
		Type       string `yaml:"type"`
		Properties map[string]struct {
			Type   string `yaml:"type"`
			Format string `yaml:"format,omitempty"`
			Items  struct {
				Ref string `yaml:"$ref"`
			} `yaml:"items,omitempty"`
		} `yaml:"properties,omitempty"`
	} `yaml:"definitions"`
	XAmazonApigatewayPolicy struct {
		Version   string `yaml:"Version"`
		Statement []struct {
			Effect    string `yaml:"Effect"`
			Principal string `yaml:"Principal"`
			Action    string `yaml:"Action"`
			Resource  string `yaml:"Resource"`
		} `yaml:"Statement"`
	} `yaml:"x-amazon-apigateway-policy"`
}

type SwaggerXAmazonApigatewayIntegration struct {
	ConnectionId string `yaml:"connectionId"`
	HttpMethod   string `yaml:"httpMethod"`
	Uri          string `yaml:"uri"`
	Responses    map[string]struct {
		StatusCode string `yaml:"statusCode"`
	} `yaml:"responses"`
	PassThroughBehavior string            `yaml:"passthroughBehavior"`
	ConnectionType      string            `yaml:"connectionType"`
	RequestParameters   map[string]string `yaml:"requestParameters"`
	Type                string            `yaml:"type"`
}

type SwaggerYAMLXAmazonApigatewayPolicy struct {
	Version   string `yaml:"Version"`
	Statement []struct {
		Effect    string `yaml:"Effect"`
		Principal string `yaml:"Principal"`
		Action    string `yaml:"Action"`
		Resource  string `yaml:"Resource"`
	} `yaml:"Statement"`
}

func main() {
	var account, host, inputPath, outputPath, region, vpc string

	flag.StringVar(&account, "account", "", "AWS account ID")
	flag.StringVar(&host, "host", "", "Cloud host URL")
	flag.StringVar(&inputPath, "input", "swagger.json", "Path to the input Swagger JSON file")
	flag.StringVar(&outputPath, "output", "swagger.yaml", "Path to the output YAML file")
	flag.StringVar(&region, "region", "", "AWS region")
	flag.StringVar(&vpc, "vpc", "", "AWS VPC ID")
	flag.Parse()

	// read the input Swagger JSON file
	jsonFile, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Error reading Swagger JSON file: %v", err)
	}

	// parse the JSON into Swagger struct
	var swagger Swagger
	if err := json.Unmarshal(jsonFile, &swagger); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// create a SwaggerYAML struct for YAML conversion
	var swaggerYAML SwaggerYAML
	jsonBytes, err := json.Marshal(swagger)
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}
	if err := yaml.Unmarshal(jsonBytes, &swaggerYAML); err != nil {
		log.Fatalf("Error unmarshaling JSON to YAML struct: %v", err)
	}

	// add AWS policy tags
	resource := fmt.Sprintf("arn:aws:execute-api:%s:%s:*/*/*/*", region, account)
	swaggerYAML.XAmazonApigatewayPolicy = SwaggerYAMLXAmazonApigatewayPolicy{
		Version: "2012-10-17",
		Statement: []struct {
			Effect    string `yaml:"Effect"`
			Principal string `yaml:"Principal"`
			Action    string `yaml:"Action"`
			Resource  string `yaml:"Resource"`
		}{
			{
				Effect:    "Deny",
				Principal: "*",
				Action:    "execute-api:Invoke",
				Resource:  resource,
			},
		},
	}

	for route, methods := range swaggerYAML.Paths {
		for method, data := range methods {
			// replace "default" with "500" in responses
			_, defaultResponseExists := data.Responses["default"]
			if defaultResponseExists {
				// extract the existing "Ref" value
				existingRef := data.Responses["default"].Schema.Ref

				// remove the existing default response
				delete(data.Responses, "default")

				// ddd a new response with status code "500"
				data.Responses["500"] = struct {
					Description string `yaml:"description"`
					Schema      struct {
						Ref string `yaml:"$ref"`
					} `yaml:"schema"`
				}{
					Description: "Internal Server Error (No Retry)",
					Schema: struct {
						Ref string `yaml:"$ref"`
					}{
						Ref: existingRef, // preserve the original "Ref" value
					},
				}
			}

			// update the x-amazon-apigateway-integration fields
			data.XAmazonApigatewayIntegration.ConnectionId = vpc
			data.XAmazonApigatewayIntegration.HttpMethod = method
			data.XAmazonApigatewayIntegration.Uri = host + route
			data.XAmazonApigatewayIntegration.Responses = map[string]struct {
				StatusCode string `yaml:"statusCode"`
			}{
				"200": {
					StatusCode: "200",
				},
				// TODO: add other status codes as needed
			}
			data.XAmazonApigatewayIntegration.PassThroughBehavior = "when_no_match"
			data.XAmazonApigatewayIntegration.ConnectionType = "VPC_LINK"

			requestParameters := make(map[string]string)

			// regular expression pattern to match path parameters
			pathParamPattern := "{(.*?)}"
			re := regexp.MustCompile(pathParamPattern)
			pathParams := re.FindAllString(route, -1)

			// add them to request parameters
			for _, pathParam := range pathParams {
				paramName := strings.Trim(pathParam, "{}")
				requestParamKey := "integration.request.path." + paramName
				requestParamValue := "method.request.path." + paramName
				requestParameters[requestParamKey] = requestParamValue
			}

			data.XAmazonApigatewayIntegration.RequestParameters = requestParameters
			data.XAmazonApigatewayIntegration.Type = "http"

			swaggerYAML.Paths[route][method] = data
		}
	}

	// marshal the SwaggerYAML struct to YAML
	yamlData, err := yaml.Marshal(&swaggerYAML)
	if err != nil {
		log.Fatalf("Error marshaling to YAML: %v", err)
	}

	// write the YAML to the output file
	if err := os.WriteFile(outputPath, yamlData, os.ModePerm); err != nil {
		log.Fatalf("Error writing YAML to output file: %v", err)
	}

	fmt.Printf("Conversion complete. YAML saved to %s\n", outputPath)
}
