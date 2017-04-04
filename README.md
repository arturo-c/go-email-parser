# go-email-parser
Go Email Parse and return user defined json fields.

## Install
### Plain GO
- Dependencies `glide install`
- Build `go build`
- Run `./go-email-parser`
- Check `curl -XPOST "http://localhost:8080/parse"`

### Docker
- `docker build -t go-email-parser:v1 .`
- `docker run -p 8080:8080 -it go-email-parser`

## Test
- `go test -v`

## Deploy
- Development `kubectl create -f deploy/minikube.yml`
- Production  `kubectl create -f deploy/aws.yml`

## TODO
- Return specified fields in params input.
- Add Jenkinsfile to automate test builds in Jenkins Pipeline.
- Change log output to logstash.
