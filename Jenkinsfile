// Configuration variables
github_org             = "samsung-cnct"
quay_org               = "samsung_cnct"





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
                        kubesh 'cd go/src/github.com/samsung-cnct/k2cli/cmd && go test -v 2>&1 | go-junit-report > k2cli_cmd.xml'
                        junit 'go/src/github.com/samsung-cnct/k2cli/cmd/k2cli_cmd.xml'
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

                // stage('Generate current docs and compare') {
                //     // generates a comparison docs folder and sees if docs need updating
                //     kubesh "go/src/github.com/samsung-cnct/k2cli/k2cli docs test"
                //     kubesh "diff -r test docs || echo The docs are not up to date. Please update your docs with the k2cli docs command. && false"
                    


                // }

                try {
                    stage('Test: Cloud') {
                        parallel ( 
                            "aws": {
                                kubesh "env helm_override_`echo ${JOB_BASE_NAME}-${BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false'
                                kubesh "go/src/github.com/samsung-cnct/k2cli/k2cli -vvv cluster up /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/"
                            },
                            "gke": {
                                kubesh "env helm_override_`echo ${JOB_BASE_NAME}-${BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false'
                                kubesh "go/src/github.com/samsung-cnct/k2cli/k2cli -vvv cluster up /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/"
                            }
                        )
                    }
                } finally {
                    stage('Cleanup') {
                        parallel (
                            "aws": {
                                kubesh "go/src/github.com/samsung-cnct/k2cli/k2cli -vvv cluster down /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/ || true"
                                kubesh "rm -rf /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws"
                            },
                            "gke": {
                                kubesh "go/src/github.com/samsung-cnct/k2cli/k2cli -vvv cluster down /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml --output /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/ || true"
                                kubesh "rm -rf /var/lib/docker/scratch/k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke"
                            }
                        )
                    }
                }
            }
            customContainer('docker') {
            // add a docker rmi/docker purge/etc.
            stage('Build') {
                kubesh "docker rmi quay.io/${quay_org}/k2cli:k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID} || true"
                kubesh "docker rmi quay.io/${quay_org}/k2cli:latest || true"
                kubesh "docker build --no-cache --force-rm -t quay.io/${quay_org}/k2cli:k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID} docker/"
            }

            //only push from master if we are on samsung-cnct fork
            // stage('Publish') {
            //     if (env.BRANCH_NAME == "master" && git_uri.contains(github_org)) {
            //         kubesh "docker tag quay.io/${quay_org}/k2cli:k2cli-${env.JOB_BASE_NAME}-${env.BUILD_ID} quay.io/${quay_org}/k2cli:latest"
            //         kubesh "docker push quay.io/${quay_org}/k2cli:latest"
            //     } else {
            //         echo "Not pushing to docker repo:\n    BRANCH_NAME='${env.BRANCH_NAME}'\n    git_uri='${git_uri}'"
            //     }
            // }
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
