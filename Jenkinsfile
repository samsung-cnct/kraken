podTemplate(label: 'k2cli', containers: [
    containerTemplate(name: 'jnlp', image: 'quay.io/samsung_cnct/custom-jnlp:0.1', args: '${computer.jnlpmac} ${computer.name}'),
    containerTemplate(name: 'golang', image: 'golang:latest', ttyEnabled: true, command: 'cat'),
    ], volumes: [
      hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
      hostPathVolume(hostPath: '/var/lib/docker/scratch', mountPath: '/var/lib/docker/scratch/')
    ]) {
        node('k2cli') {
            customContainer('golang') {

                stage('hello!') {
                    echo 'hello world!'
                }

                stage('checkout') {
                    checkout scm
                    kubesh 'go version'
                }

                stage('build') {
                    kubesh 'go get -v -d -t ./... || true'
                    kubesh 'GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o k2cli'
                }

                stage('anything') {
                    kubesh 'touch /var/lib/docker/scratch/poop && ls /var/lib/docker/scratch'
                }

                stage('aws config generation') {
                    kubesh './k2cli generate /var/lib/docker/scratch/config.yaml'
                }

                stage('cat config file') {
                    kubesh 'cat /var/lib/docker/scratch/config.yaml'
                }

                stage('update generated aws config') {
                    kubesh "build-scripts/update-generated-config.sh /var/lib/docker/scratch/config.yaml ${env.JOB_BASE_NAME}-${env.BUILD_ID}"
                }

                stage("read config file again") {
                    kubesh 'cat /var/lib/docker/scratch/config.yaml'
                }



            }

        }
    }
def kubesh(command) {
  if (env.CONTAINER_NAME) {
    if ((command instanceof String) || (command instanceof GString)) {
      command = kubectl(command)
    }

    if (command instanceof LinkedHashMap) {
      command["script"] = kubectl(command["script"])
    }
  }

  sh(command)
}

def kubectl(command) {
  "kubectl exec -i ${env.HOSTNAME} -c ${env.CONTAINER_NAME} -- /bin/sh -c 'cd ${env.WORKSPACE} && ${command}'"
}

def customContainer(String name, Closure body) {
  withEnv(["CONTAINER_NAME=$name"]) {
    body()
  }
}
