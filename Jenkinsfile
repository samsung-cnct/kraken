podTemplate(label: 'k2cli', containers: [
    containerTemplate(name: 'jnlp', image: 'jenkinsci/jnlp-slave:2.62-alpine', args: '${computer.jnlpmac} ${computer.name}'),
    containerTemplate(name: 'golang', image: 'golang:latest', ttyEnabled: true, command: 'cat')
    ]) {
        node('k2cli') {
            container('golang'){

                stage('hello!') {
                    echo 'hello world!'
                }

                stage('checkout'){
                    checkout scm
                    sh 'go version'


                }

                stage('build'){
                    sh 'go get -v -d -t ./... || true'
                    sh 'GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v'
                }

                stage('aws config generation') {
                    sh 'k2cli generate'
                    sh 'cat ${HOME}/.kraken/config.yaml'
                }

            }

        }
    }
