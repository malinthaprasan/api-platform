/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package model

import (
	"time"
)

// API represents an API entity in the platform
type API struct {
	UUID             string              `json:"uuid,omitempty" db:"uuid"`
	Name             string              `json:"name" db:"name"`
	DisplayName      string              `json:"displayName,omitempty" db:"display_name"`
	Description      string              `json:"description,omitempty" db:"description"`
	Context          string              `json:"context" db:"context"`
	Version          string              `json:"version" db:"version"`
	Provider         string              `json:"provider,omitempty" db:"provider"`
	ProjectID        string              `json:"projectId" db:"project_id"` // FK to Project.UUID
	CreatedAt        time.Time           `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt        time.Time           `json:"updatedAt,omitempty" db:"updated_at"`
	LifeCycleStatus  string              `json:"lifeCycleStatus,omitempty" db:"lifecycle_status"`
	HasThumbnail     bool                `json:"hasThumbnail,omitempty" db:"has_thumbnail"`
	IsDefaultVersion bool                `json:"isDefaultVersion,omitempty" db:"is_default_version"`
	IsRevision       bool                `json:"isRevision,omitempty" db:"is_revision"`
	RevisionedAPIID  string              `json:"revisionedApiId,omitempty" db:"revisioned_api_id"`
	RevisionID       int                 `json:"revisionId,omitempty" db:"revision_id"`
	Type             string              `json:"type,omitempty" db:"type"`
	Transport        []string            `json:"transport,omitempty" db:"transport"`
	MTLS             *MTLSConfig         `json:"mtls,omitempty"`
	Security         *SecurityConfig     `json:"security,omitempty"`
	CORS             *CORSConfig         `json:"cors,omitempty"`
	BackendServices  []BackendService    `json:"backend-services,omitempty"`
	APIRateLimiting  *RateLimitingConfig `json:"api-rate-limiting,omitempty"`
	Operations       []Operation         `json:"operations,omitempty"`
}

// TableName returns the table name for the API model
func (API) TableName() string {
	return "apis"
}

// MTLSConfig represents mutual TLS configuration
type MTLSConfig struct {
	Enabled                    bool   `json:"enabled,omitempty"`
	EnforceIfClientCertPresent bool   `json:"enforceIfClientCertPresent,omitempty"`
	VerifyClient               bool   `json:"verifyClient,omitempty"`
	ClientCert                 string `json:"clientCert,omitempty"`
	ClientKey                  string `json:"clientKey,omitempty"`
	CACert                     string `json:"caCert,omitempty"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	Enabled bool            `json:"enabled,omitempty"`
	APIKey  *APIKeySecurity `json:"apiKey,omitempty"`
	OAuth2  *OAuth2Security `json:"oauth2,omitempty"`
}

// APIKeySecurity represents API key security configuration
type APIKeySecurity struct {
	Enabled bool   `json:"enabled,omitempty"`
	Header  string `json:"header,omitempty"`
	Query   string `json:"query,omitempty"`
	Cookie  string `json:"cookie,omitempty"`
}

// OAuth2Security represents OAuth2 security configuration
type OAuth2Security struct {
	GrantTypes *OAuth2GrantTypes `json:"grantTypes,omitempty"`
	Scopes     []string          `json:"scopes,omitempty"`
}

// OAuth2GrantTypes represents OAuth2 grant types configuration
type OAuth2GrantTypes struct {
	AuthorizationCode *AuthorizationCodeGrant `json:"authorizationCode,omitempty"`
	Implicit          *ImplicitGrant          `json:"implicit,omitempty"`
	Password          *PasswordGrant          `json:"password,omitempty"`
	ClientCredentials *ClientCredentialsGrant `json:"clientCredentials,omitempty"`
}

// AuthorizationCodeGrant represents authorization code grant configuration
type AuthorizationCodeGrant struct {
	Enabled     bool   `json:"enabled,omitempty"`
	CallbackURL string `json:"callbackUrl,omitempty"`
}

// ImplicitGrant represents implicit grant configuration
type ImplicitGrant struct {
	Enabled     bool   `json:"enabled,omitempty"`
	CallbackURL string `json:"callbackUrl,omitempty"`
}

// PasswordGrant represents password grant configuration
type PasswordGrant struct {
	Enabled bool `json:"enabled,omitempty"`
}

// ClientCredentialsGrant represents client credentials grant configuration
type ClientCredentialsGrant struct {
	Enabled bool `json:"enabled,omitempty"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	Enabled          bool   `json:"enabled,omitempty"`
	AllowOrigins     string `json:"allowOrigins,omitempty"`
	AllowMethods     string `json:"allowMethods,omitempty"`
	AllowHeaders     string `json:"allowHeaders,omitempty"`
	ExposeHeaders    string `json:"exposeHeaders,omitempty"`
	MaxAge           int    `json:"maxAge,omitempty"`
	AllowCredentials bool   `json:"allowCredentials,omitempty"`
}

// BackendService represents a backend service configuration
type BackendService struct {
	Name           string                `json:"name,omitempty"`
	IsDefault      bool                  `json:"isDefault,omitempty"`
	Endpoints      []BackendEndpoint     `json:"endpoints,omitempty"`
	Timeout        *TimeoutConfig        `json:"timeout,omitempty"`
	Retries        int                   `json:"retries,omitempty"`
	LoadBalance    *LoadBalanceConfig    `json:"loadBalance,omitempty"`
	CircuitBreaker *CircuitBreakerConfig `json:"circuitBreaker,omitempty"`
}

// BackendEndpoint represents a backend endpoint
type BackendEndpoint struct {
	URL         string             `json:"url,omitempty"`
	Description string             `json:"description,omitempty"`
	HealthCheck *HealthCheckConfig `json:"healthCheck,omitempty"`
	Weight      int                `json:"weight,omitempty"`
	MTLS        *MTLSConfig        `json:"mtls,omitempty"`
}

// HealthCheckConfig represents health check configuration
type HealthCheckConfig struct {
	Enabled            bool `json:"enabled,omitempty"`
	Interval           int  `json:"interval,omitempty"`
	Timeout            int  `json:"timeout,omitempty"`
	UnhealthyThreshold int  `json:"unhealthyThreshold,omitempty"`
	HealthyThreshold   int  `json:"healthyThreshold,omitempty"`
}

// TimeoutConfig represents timeout configuration
type TimeoutConfig struct {
	Connect int `json:"connect,omitempty"`
	Read    int `json:"read,omitempty"`
	Write   int `json:"write,omitempty"`
}

// LoadBalanceConfig represents load balancing configuration
type LoadBalanceConfig struct {
	Algorithm string `json:"algorithm,omitempty"`
	Failover  bool   `json:"failover,omitempty"`
}

// CircuitBreakerConfig represents circuit breaker configuration
type CircuitBreakerConfig struct {
	Enabled            bool `json:"enabled,omitempty"`
	MaxConnections     int  `json:"maxConnections,omitempty"`
	MaxPendingRequests int  `json:"maxPendingRequests,omitempty"`
	MaxRequests        int  `json:"maxRequests,omitempty"`
	MaxRetries         int  `json:"maxRetries,omitempty"`
}

// RateLimitingConfig represents rate limiting configuration
type RateLimitingConfig struct {
	Enabled           bool   `json:"enabled,omitempty"`
	RateLimitCount    int    `json:"rateLimitCount,omitempty"`
	RateLimitTimeUnit string `json:"rateLimitTimeUnit,omitempty"`
	StopOnQuotaReach  bool   `json:"stopOnQuotaReach,omitempty"`
}

// Operation represents an API operation
type Operation struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Request     *OperationRequest `json:"request,omitempty"`
}

// OperationRequest represents operation request details
type OperationRequest struct {
	Method           string                `json:"method,omitempty"`
	Path             string                `json:"path,omitempty"`
	BackendServices  []BackendRouting      `json:"backend-services,omitempty"`
	Authentication   *AuthenticationConfig `json:"authentication,omitempty"`
	RequestPolicies  []Policy              `json:"requestPolicies,omitempty"`
	ResponsePolicies []Policy              `json:"responsePolicies,omitempty"`
}

// BackendRouting represents backend routing configuration
type BackendRouting struct {
	Name   string `json:"name,omitempty"`
	Weight int    `json:"weight,omitempty"`
}

// AuthenticationConfig represents authentication configuration for operations
type AuthenticationConfig struct {
	Required bool     `json:"required,omitempty"`
	Scopes   []string `json:"scopes,omitempty"`
}

// Policy represents a request or response policy
type Policy struct {
	Name   string                 `json:"name,omitempty"`
	Params map[string]interface{} `json:"params,omitempty"`
}
