podTemplate(label: 'k2cli', containers: [
    containerTemplate(name: 'jnlp', image: 'quay.io/samsung_cnct/custom-jnlp:0.1', args: '${computer.jnlpmac} ${computer.name}'),
    containerTemplate(name: 'golang', image: 'quay.io/guineveresaenger/guinsci:latest', ttyEnabled: true, command: 'cat'),
    containerTemplate(name: 'k2-tools', image: 'quay.io/samsung_cnct/k2-tools:latest', ttyEnabled: true, command: 'cat', alwaysPullImage: true, resourceRequestMemory: '1Gi', resourceLimitMemory: '1Gi'),
    ], volumes: [
      hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
      hostPathVolume(hostPath: '/var/lib/docker/scratch', mountPath: '/var/lib/docker/scratch/')
    ]) {
        node('k2cli') {
            customContainer('golang') {

                stage('Checkout') {
                    checkout scm
                    kubesh "mkdir -p go/src/github.com/samsung-cnct/k2cli/ && cp -r `ls -A | grep -v \"go\"` main.go go/src/github.com/samsung-cnct/k2cli"
                    kubesh "ls -r go/src/github.com/samsung-cnct/k2cli/*"
                    git_uri = scm.getRepositories()[0].getURIs()[0].toString()
                    kubesh "echo ${git_uri}"
                }

                withEnv(["GOPATH=${WORKSPACE}/go/"]) {
                    stage('Test: Unit') {
                        kubesh 'cd go/src/github.com/samsung-cnct/k2cli/ && gosimple .'
                        kubesh 'cd go/src/github.com/samsung-cnct/k2cli/ && make deps && make build'
                        kubesh 'cd go/src/github.com/samsung-cnct/k2cli/ && go vet'
                        kubesh 'cd go/src/github.com/samsung-cnct/k2cli/cmd && go test -v'
                    }

                    stage('Build') {
                        kubesh 'cd go/src/github.com/samsung-cnct/k2cli/ && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o k2cli'
                    }
                }
            }

            customContainer('k2-tools') {
                stage('Configure Integration Tests') {
                    // fetches credentials, builds aws and gke config files with appropriate replacements
                    kubesh "go/src/github.com/samsung-cnct/k2cli/build-scripts/fetch-credentials.sh /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/"
                    kubesh "go/src/github.com/samsung-cnct/k2cli/k2cli generate --provider aws /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml"
                    kubesh "go/src/github.com/samsung-cnct/k2cli/build-scripts/update-generated-config.sh /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml kca${env.JOB_BASE_NAME}-${env.BUILD_ID} /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}"
                    kubesh "go/src/github.com/samsung-cnct/k2cli/k2cli generate --provider gke /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml"
                    kubesh "go/src/github.com/samsung-cnct/k2cli/build-scripts/update-generated-config.sh /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml kcg${env.JOB_BASE_NAME}-${env.BUILD_ID} /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}"

                }

                try {
                    stage('Test: Cloud') {
                        parallel ( 
                            "aws": {
                                kubesh "env helm_override_kca`echo ${env.JOB_BASE_NAME}-${env.BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false' + " go/src/github.com/samsung-cnct/k2cli/k2cli -vvv cluster up /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/"
                            },
                            "gke": {
                                kubesh "env helm_override_kcg`echo ${env.JOB_BASE_NAME}-${env.BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false' + " go/src/github.com/samsung-cnct/k2cli/k2cli -vvv cluster up /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/"
                            }
                        )
                    }
                } finally {
                    stage('Cleanup') {
                        parallel (
                            "aws": {
                                kubesh "env helm_override_kca`echo ${env.JOB_BASE_NAME}-${env.BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false' + " go/src/github.com/samsung-cnct/k2cli/k2cli -vvv cluster down /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/ || true"
                                kubesh "rm -rf /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws"
                            },
                            "gke": {
                                kubesh "env helm_override_kcg`echo ${env.JOB_BASE_NAME}-${env.BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false' + " go/src/github.com/samsung-cnct/k2cli/k2cli -vvv cluster down /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/ || true"
                                kubesh "rm -rf /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke"
                            }
                        )
                    }
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
  "kubectl exec -i ${env.HOSTNAME} -c ${env.CONTAINER_NAME} -- /bin/sh -c 'cd ${env.WORKSPACE} && export GOPATH=${env.GOPATH} && ${command}'"
}

def customContainer(String name, Closure body) {
  withEnv(["CONTAINER_NAME=$name"]) {
    body()
  }
}
