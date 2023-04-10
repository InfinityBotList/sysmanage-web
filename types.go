package main

type Config struct {
	AllowedUsers []string `yaml:"allowed_users"`
	InfraFolder  string   `yaml:"infra_folder"`
	DPSecret     string   `yaml:"dp_secret"`
}
