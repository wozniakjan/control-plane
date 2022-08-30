#!/usr/bin/env bash

set -xeuo pipefail

if [[ -n "$(git diff)" || -n "$(git diff --cached)" ]]; then
    echo "Refusing to run on dirty repository" 1>&2
    exit 1
fi


DIR=$(dirname "${BASH_SOURCE[0]}")/..

current="$(git rev-parse --abbrev-ref HEAD)"
pr=$(gh pr list --author "@me" | awk "/$current/{print(\$1)}")
if [[ "$pr" == "" ]]; then
    echo "No PR for branch $current" 1>&2
    exit 2
fi


while read component; do
   case $component in
       kyma-environment-broker)
            awk -i inplace '/    kyma_environment_broker:/{b=1} /      version: "/{if(b != 1){print($0)}else{print("      version: \"'PR-$pr'\"")};b=0} !/      version: "/{print($0)}' "$DIR/resources/kcp/charts/kyma-environment-broker/values.yaml"
            #awk -i inplace '/    kyma_environments_cleanup_job:/{b=1} /      version: "/{if(b != 1){print($0)}else{print("      version: \"'PR-$pr'\"")};b=0} !/      version: "/{print($0)}' "$DIR/resources/kcp/values.yaml"
            #awk -i inplace '/    kyma_environments_subaccount_cleanup_job:/{b=1} /      version: "/{if(b != 1){print($0)}else{print("      version: \"'PR-$pr'\"")};b=0} !/      version: "/{print($0)}' "$DIR/resources/kcp/values.yaml"
            #awk -i inplace '/    kyma_environments_subscription_cleanup_job:/{b=1} /      version: "/{if(b != 1){print($0)}else{print("      version: \"'PR-$pr'\"")};b=0} !/      version: "/{print($0)}' "$DIR/resources/kcp/values.yaml"
           ;;
       schema-migrator)
            awk -i inplace '/    schema_migrator:/{b=1} /      version: "/{if(b != 1){print($0)}else{print("      version: \"'PR-$pr'\"")};b=0} !/      version: "/{print($0)}' "$DIR/resources/kcp/values.yaml"
           ;;
   esac
done < <(git diff-tree --no-commit-id --name-only -r ..main | awk -F / '{print $2}' | sort | uniq )

git commit -am "Bump kyma_environment_broker images to PR-$pr"
