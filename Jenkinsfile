podTemplate(label: 'k2cli', containers: [
    containerTemplate(name: 'jnlp', image: 'jenkinsci/jnlp-slave:2.62-alpine', args: '${computer.jnlpmac} ${computer.name}'),
    containerTemplate(name: 'golang', image: 'golang:latest', ttyEnabled: true, command: 'cat'),
    containerTemplate(name: 'k2-tools', image: 'quay.io/samsung_cnct/k2-tools:latest', ttyEnabled: true, command: 'cat', alwaysPullImage: true, resourceRequestMemory: '1Gi', resourceLimitMemory: '1Gi')
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

            }
            container('k2-tools'){

                stage('checkout') {
                    checkout scm
                }

                stage('aws config generation') {
                    echo WORKSPACE
                    sh 'cd WORKSPACE && k2cli generate'

                }
            }

        }
    }

// def customContainer(String name, Closure body) {
//   withEnv(["CONTAINER_NAME=$name"]) {
//     body()
//   }
// }
