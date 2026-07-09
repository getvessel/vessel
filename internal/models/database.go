package models

import "time"

type Database struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	Engine        string    `json:"engine"`
	Version       string    `json:"version"`
	Port          int       `json:"port"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	DatabaseName  string    `json:"databaseName"`
	VolumePath    string    `json:"volumePath"`
	ContainerID   string    `json:"containerId"`
	Status        string    `json:"status"`
	InternalDNS   string    `json:"internalDns"`
	ExternalDNS   string    `json:"externalDns"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type CreateDatabaseRequest struct {
	ProjectID     string `json:"projectId"`
	EnvironmentID string `json:"environmentId"`
	Name          string `json:"name"`
	Engine        string `json:"engine"`
	Version       string `json:"version"`
	Port          int    `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	DatabaseName  string `json:"databaseName"`
	VolumePath    string `json:"volumePath"`
}

type Storage struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	APIPort       int       `json:"apiPort"`
	ConsolePort   int       `json:"consolePort"`
	AccessKey     string    `json:"accessKey"`
	SecretKey     string    `json:"secretKey,omitempty"`
	BucketName    string    `json:"bucketName"`
	VolumePath    string    `json:"volumePath"`
	ContainerID   string    `json:"containerId"`
	Status        string    `json:"status"`
	InternalDNS   string    `json:"internalDns"`
	ExternalDNS   string    `json:"externalDns"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
