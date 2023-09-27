#!/bin/sh

APP=$(find / -type f -name swagger2yaml)
ACCOUNT="$INPUT_ACCOUNT"
HOST="$INPUT_HOST"
INPUT_FILE="$INPUT_INPUT"
OUTPUT_FILE="$INPUT_OUTPUT"
REGION="$INPUT_REGION"
VPC="$INPUT_VPC"

if [ -z "$ACCOUNT" ] || [ -z "$HOST" ] || [ -z "$INPUT_FILE" ] || [ -z "$OUTPUT_FILE" ] || [ -z "$REGION" ] || [ -z "$VPC" ]; then
  echo "One or more required variables are missing."
  exit 1
fi

$APP --account "$ACCOUNT" --host "$HOST" --input "$INPUT_FILE" --output "$OUTPUT_FILE" --region "$REGION" --vpc "$VPC"
