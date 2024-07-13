#!/bin/bash
# scripts/generate_proto.sh

# Exit on any error
set -e

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Directory containing .proto files
PROTO_DIR="${PROJECT_ROOT}/api/proto"

# Directory to output generated files
GO_OUT_DIR="${PROJECT_ROOT}/api/proto"

# Create the output directory if it doesn't exist
mkdir -p $GO_OUT_DIR

# Generate Go files for each .proto file
for proto_file in $PROTO_DIR/*.proto; do
    # Extract the service name from the file name
    service_name=$(basename "$proto_file" .proto)
    
    # Create a subdirectory for each service
    service_out_dir="${GO_OUT_DIR}/${service_name}"
    mkdir -p $service_out_dir
    
    protoc --go_out=$service_out_dir --go_opt=paths=source_relative \
           --go-grpc_out=$service_out_dir --go-grpc_opt=paths=source_relative \
           -I$PROTO_DIR $proto_file

    echo "Generated files for ${service_name}"
done

echo "Proto files generat