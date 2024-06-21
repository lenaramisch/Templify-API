# How to use this app

## Build api spec

`rm -f ${PWD}/api-spec/bundle.yml`
`docker run --rm -v $PWD/api-spec/:/spec ghcr.io/redocly/cli:v1.0.0 lint openapi.yml`
`docker run --rm -v $PWD/api-spec/:/spec ghcr.io/redocly/cli:v1.0.0 bundle openapi.yml >> ./api-spec/bundle.yml`
