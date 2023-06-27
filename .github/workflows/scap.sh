#!/bin/bash

set -e -o pipefail
set -x

xsd2go=$(pwd)/gocomply_xsd2go

workdir=./scap/
mkdir -p $workdir
cat <<__END__ > $workdir/go.mod
module github.com/gocomply/scap
go 1.17
__END__

pushd $workdir

    # Acquire XSDs
    [ -d .scap_schemas ] || git clone --depth 1 https://github.com/openscap/openscap .scap_schemas

    # Clean-up the workspace
    [ -d pkg/scap/models ] && find pkg/scap/models -name models.go | xargs rm --

    # Generage go code based on XSDs
    $xsd2go convert .scap_schemas/schemas/cpe/2.3/cpe-dictionary_2.3.xsd github.com/gocomply/scap pkg/scap/models
    $xsd2go convert \
            --xmlns-override=http://cpe.mitre.org/language/2.0=cpe_language \
            .scap_schemas/schemas/xccdf/1.2/xccdf_1.2.xsd github.com/gocomply/scap pkg/scap/models
	  $xsd2go convert .scap_schemas/schemas/oval/5.11.3/oval-results-schema.xsd github.com/gocomply/scap pkg/scap/models
	  $xsd2go convert \
            --xmlns-override=http://cpe.mitre.org/language/2.0=cpe_language \
            .scap_schemas/schemas/sds/1.3/scap-source-data-stream_1.3.xsd github.com/gocomply/scap pkg/scap/models

    # Ensure the code can be compiled
    go vet ./...

popd

