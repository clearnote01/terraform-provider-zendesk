VERSION := 0.0.4

concat: provider/dev_provider.tf resources/**/resource.tf
	cat provider/dev_provider.tf resources/**/resource.tf > main.tf
install: ../main.go ../zendesk/*.go
	mkdir -p build
	go build -o build/terraform-provider-zendesk_v$(VERSION) ../.
plan: concat install dev.tfrc
	export TF_LOG=trace
	TF_CLI_CONFIG_FILE=dev.tfrc terraform plan
apply: concat install dev.tfrc
	export TF_LOG=trace
	TF_CLI_CONFIG_FILE=dev.tfrc terraform apply -auto-approve
destroy: concat install dev.tfrc
	export TF_LOG=trace
	TF_CLI_CONFIG_FILE=dev.tfrc terraform destroy