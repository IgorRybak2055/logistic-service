// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-05-03 13:07:37.586232067 +0300 +03 m=+0.172406389

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "support@logistic.com"
        },
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/company/register": {
            "post": {
                "description": "Create a new company.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Company"
                ],
                "summary": "Registration company of logistic-service service",
                "parameters": [
                    {
                        "description": "company",
                        "name": "company",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Company"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "response structure: {message:\"answer\"}",
                        "schema": {
                            "$ref": "#/definitions/models.Account"
                        }
                    },
                    "400": {
                        "description": "response structure: {error:\"error message\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "response structure: {error:\"error message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/deliveries": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Creating new delivery.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Delivery"
                ],
                "summary": "Creating new delivery.",
                "parameters": [
                    {
                        "description": "delivery",
                        "name": "delivery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Delivery"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "response structure: {message:delivery}",
                        "schema": {
                            "$ref": "#/definitions/models.Delivery"
                        }
                    },
                    "400": {
                        "description": "response structure: {error:\"error message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/login": {
            "post": {
                "description": "Login in ragger with email and password",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "summary": "Login in ragger",
                "parameters": [
                    {
                        "type": "string",
                        "description": "account email",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "account password len(password) \u003e 6",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "response structure: {message:\"answer\"}",
                        "schema": {
                            "$ref": "#/definitions/models.Account"
                        }
                    },
                    "400": {
                        "description": "response structure: {error:\"error message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/new_password": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Setting new password for account",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "summary": "Set new password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "new password for account",
                        "name": "password",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "confirm_password new password",
                        "name": "confirm_password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "400": {
                        "description": "response structure: {error:\"error message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/register": {
            "post": {
                "description": "Create a new user with the input name, email \u0026 password.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "summary": "Registration user of contact service",
                "parameters": [
                    {
                        "description": "account",
                        "name": "account",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Account"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "response structure: {message:\"answer\"}",
                        "schema": {
                            "$ref": "#/definitions/models.Account"
                        }
                    },
                    "400": {
                        "description": "response structure: {error:\"error message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/restore_password": {
            "get": {
                "description": "Returned html page for setting new password",
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "Account"
                ],
                "summary": "Restore password",
                "responses": {
                    "200": {},
                    "500": {
                        "description": "response structure: {error:\"error message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/token": {
            "get": {
                "description": "Generate Token for access to ragger",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Token"
                ],
                "summary": "Generate Token in ragger",
                "parameters": [
                    {
                        "type": "string",
                        "description": "last generated refresh_token",
                        "name": "refresh_token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "response structure: {message:\"access_token:access, refresh_token:refresh\"}",
                        "schema": {
                            "$ref": "#/definitions/string"
                        }
                    },
                    "401": {
                        "description": "response structure: {error:\"error message\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Health check ragger service",
                "produces": [
                    "application/json"
                ],
                "summary": "Ragger health check",
                "responses": {
                    "200": {
                        "description": "response structure: {status:\"UP\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Account": {
            "type": "object",
            "properties": {
                "company_id": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "tokens": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "models.Company": {
            "type": "object",
            "properties": {
                "bank_detail": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "kind": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "models.Delivery": {
            "type": "object",
            "properties": {
                "cargo": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "price": {
                    "type": "number"
                },
                "shipment_date": {
                    "type": "string"
                },
                "shipment_place": {
                    "type": "string"
                },
                "trailer_type": {
                    "type": "string"
                },
                "unloading_place": {
                    "type": "string"
                },
                "volume_cargo": {
                    "type": "number"
                },
                "weight_cargo": {
                    "type": "number"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8888",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Logistic API",
	Description: "This is a sample service ...",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
