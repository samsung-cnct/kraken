// Configuration variables
github_org             = "samsung-cnct"
quay_org               = "samsung_cnct"
publish_branch         = "master"
release_version        = "${env.RELEASE_VERSION}"
k2_image_tag           = "${env.K2_VERSION}" != "null" ? "${env.K2_VERSION}" : "latest"
release_branch         = "${env.REL_BRANCH}"

podTemplate(label: 'kraken',
    containers: [
        containerTemplate(name: 'jnlp', image: 'quay.io/samsung_cnct/custom-jnlp:0.1', args: '${computer.jnlpmac} ${computer.name}'),
        containerTemplate(name: 'golang', image: 'quay.io/samsung_cnct/kraken-gobuild:1.8.3', ttyEnabled: true, command: 'cat', alwaysPullImage: true),
        containerTemplate(name: 'kraken-tools', image: 'quay.io/samsung_cnct/k2-tools:latest', ttyEnabled: true, command: 'cat', alwaysPullImage: true, resourceRequestMemory: '1Gi', resourceLimitMemory: '1Gi'),
    ],
    volumes: [
        hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
        hostPathVolume(hostPath: '/var/lib/docker/scratch', mountPath: '/var/lib/docker/scratch/'),
        secretVolume(mountPath: '/home/jenkins/kraken-release-token/', secretName: 'kraken-publish-token')
    ])
    {
        node('kraken') {
            customContainer('golang') {

                stage('Checkout') {
                    checkout scm
                    kubesh "mkdir -p go/src/github.com/samsung-cnct/kraken/ && cp -r `ls -A | grep -v \"go\"` main.go go/src/github.com/samsung-cnct/kraken"
                    kubesh "ls -r go/src/github.com/samsung-cnct/kraken/*"
                    git_uri = scm.getRepositories()[0].getURIs()[0].toString()
                    git_branch = scm.getBranches()[0].toString()
                }

                withEnv(["GOPATH=${WORKSPACE}/go/"]) {
                    stage('Test: Unit') {
                        kubesh "cd go/src/github.com/samsung-cnct/kraken/ && gosimple ."
                        kubesh "cd go/src/github.com/samsung-cnct/kraken/ && make deps && make build KLIB_VER=${k2_image_tag}"
                        kubesh "cd go/src/github.com/samsung-cnct/kraken/ && go vet"
                        kubesh "cd go/src/github.com/samsung-cnct/kraken/cmd && go test -v"
                    }

                    stage('Build') {
                        kubesh 'cd go/src/github.com/samsung-cnct/kraken/ && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o kraken'
                    }
                }
            }

            customContainer('kraken-tools') {
                stage('Configure Integration Tests') {
                    // fetches credentials, builds aws and gke config files with appropriate replacements
                    kubesh "go/src/github.com/samsung-cnct/kraken/build-scripts/fetch-credentials.sh /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/"
                    kubesh "go/src/github.com/samsung-cnct/kraken/kraken generate --provider aws /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml"
                    kubesh "go/src/github.com/samsung-cnct/kraken/build-scripts/update-generated-config.sh /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml kca${env.JOB_BASE_NAME}-${env.BUILD_ID} /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}"
                    kubesh "go/src/github.com/samsung-cnct/kraken/kraken generate --provider gke /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml"
                    kubesh "go/src/github.com/samsung-cnct/kraken/build-scripts/update-generated-config.sh /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml kcg${env.JOB_BASE_NAME}-${env.BUILD_ID} /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}"

                }

                try {
                    stage('Test: Cloud') {
                        parallel (
                            "aws": {
                                kubesh "env helm_override_kca`echo ${env.JOB_BASE_NAME}-${env.BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false' + " go/src/github.com/samsung-cnct/kraken/kraken -vvv cluster up --config /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml --output /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/"
                            },
                            "gke": {
                                kubesh "env helm_override_kcg`echo ${env.JOB_BASE_NAME}-${env.BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false' + " go/src/github.com/samsung-cnct/kraken/kraken -vvv cluster up --config /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml --output /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/"
                            }
                        )
                    }
                } finally {
                    stage('Cleanup') {
                        parallel (
                            "aws": {
                                kubesh "env helm_override_kca`echo ${env.JOB_BASE_NAME}-${env.BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false' + " go/src/github.com/samsung-cnct/kraken/kraken -vvv cluster down --config /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/config.yaml --output /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws/ || true"
                                kubesh "rm -rf /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/aws"
                            },
                            "gke": {
                                kubesh "env helm_override_kcg`echo ${env.JOB_BASE_NAME}-${env.BUILD_ID} " + '| tr \'[:upper:]\' \'[:lower:]\' | tr \'-\' \'_\'`=false' + " go/src/github.com/samsung-cnct/kraken/kraken -vvv cluster down --config /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/config.yaml --output /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke/ || true"
                                kubesh "rm -rf /var/lib/docker/scratch/kraken-${env.JOB_BASE_NAME}-${env.BUILD_ID}/gke"
                            }
                        )
                    }
                }
            }

            if (git_branch.contains(publish_branch) && git_uri.contains(github_org)) {
                customContainer('golang') {
                    withEnv(["GOPATH=${WORKSPACE}/go/"]) {
                        stage('Release') {
                            kubesh ". /home/jenkins/kraken-release-token/token && make release VERSION=${release_version} KLIB_VER=${k2_image_tag} REL_BRANCH=${release_branch}"
                        }
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
