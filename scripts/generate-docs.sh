#!/usr/bin/env bash
set -euo pipefail

# tfplugindocs introspects the provider schema through the Terraform CLI.
echo "Active terraform: $(terraform version 2>/dev/null | head -1 || echo 'not found on PATH')"

go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate \
  --provider-name guance \
  --examples-dir examples

find docs -name '*.md' -print0 | xargs -0 perl -pi -e 's/[ \t]+$//'
