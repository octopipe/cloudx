# read -p "Shared infra:" sharedinfra
# read -p "execution id:" executionid
# read -p "action:" action
RAW=$(kubectl get sharedinfra $1 -o json | jq -R -s)
# echo $RAW
go run cmd/runner/*.go $3 $2 "$RAW"