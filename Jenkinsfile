podTemplate(label: 'k2cli', containers: [
    containerTemplate(name: 'jnlp', image: 'quay.io/samsung_cnct/custom-jnlp:0.1', args: '${computer.jnlpmac} ${computer.name}'),
    containerTemplate(name: 'golang', image: 'golang:latest', ttyEnabled: true, command: 'cat'),
    containerTemplate(name: 'k2-tools', image: 'quay.io/samsung_cnct/k2-tools:latest', ttyEnabled: true, command: 'cat', alwaysPullImage: true, resourceRequestMemory: '1Gi', resourceLimitMemory: '1Gi'),
    ], volumes: [
      hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
      hostPathVolume(hostPath: '/var/lib/docker/scratch', mountPath: '/var/lib/docker/scratch/')
    ]) {
        node('k2cli') {
            customContainer('golang') {

                stage('Checkout') {
                    checkout scm
                }

                stage('Test: Unit') {
                    kubesh 'go vet'
                    kubesh 'go get -u github.com/jstemmer/go-junit-report'
                    kubesh 'go test -v ./... 2>&1 | go-junit-report > top_report.xml'
                    kubesh 'go test -v cmd ./... 2>&1 | go-junit-report > cmd_report.xml'
                    junit "*_report.xml"
                }

                stage('Build') {
                    kubesh 'go get -v -d -t ./... || true'
                    kubesh 'GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o k2cli'
                }
            }

            customContainer('k2-tools') {
                stage('Configure Integration Tests') {
                    // fetches credentials, builds aws and gke config files with appropriate replacements
                    kubesh "build-scripts/fetch-credentials.sh /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/"
                    kubesh "./k2cli generate --provider aws /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml"
                    kubesh "build-scripts/update-generated-config.sh /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml kc${env.JOB_BASE_NAME}-${env.BUILD_ID} /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}"
                    kubesh "./k2cli generate --provider gke /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml"
                    kubesh "build-scripts/update-generated-config.sh /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml kc${env.JOB_BASE_NAME}-${env.BUILD_ID} /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}"

                }

                try {
                    stage('Test: Cloud') {
                        parallel ( 
                            "aws": {
                                kubesh "./k2cli -vvv cluster up /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/"
                            }
                            "gke": {
                                kubesh "./k2cli -vvv cluster up /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/"
                            }
                        )
                    }
                } finally {
                    stage('Cleanup') {
                        parallel (
                            "aws": {
                                kubesh "./k2cli -vvv cluster down /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/ || true"
                                kubesh "rm -rf /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws"
                            }
                            "gke": {
                                kubesh "./k2cli -vvv cluster down /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/ || true"
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
  "kubectl exec -i ${env.HOSTNAME} -c ${env.CONTAINER_NAME} -- /bin/sh -c 'cd ${env.WORKSPACE} && ${command}'"
}

def customContainer(String name, Closure body) {
  withEnv(["CONTAINER_NAME=$name"]) {
    body()
  }
}
