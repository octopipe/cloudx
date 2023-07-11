# read -p "Shared infra:" infra
# read -p "execution id:" executionid
# read -p "action:" action
# RAW=$(kubectl get infra $1 -o json | jq -R -s)
# echo $RAW
go run cmd/runner/*.go $1 $2