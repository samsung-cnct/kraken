podTemplate(label: 'k2cli', containers: [
    containerTemplate(name: 'jnlp', image: 'quay.io/samsung_cnct/custom-jnlp:0.1', args: '${computer.jnlpmac} ${computer.name}'),
    containerTemplate(name: 'golang', image: 'golang:latest', ttyEnabled: true, command: 'cat'),
    containerTemplate(name: 'docker', image: 'docker', command: 'cat', ttyEnabled: true)
    ], volumes: [
      hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
      hostPathVolume(hostPath: '/var/lib/docker/scratch', mountPath: '/mnt/scratch'),
      secretVolume(mountPath: '/home/jenkins/.docker/', secretName: 'samsung-cnct-quay-robot-dockercfg')
    ]) {
        node('k2cli') {
            container('golang') {

                stage('hello!') {
                    echo 'hello world!'
                }

                stage('checkout') {
                    checkout scm
                    sh 'go version'
                }

                stage('build') {
                    sh 'go get -v -d -t ./... || true'
                    sh 'GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o k2cli'
                }

                stage('fetch credentials') {
                    sh 'build-scripts/fetch-credentials.sh'
                }

                stage('aws config generation') {
                    sh './k2cli generate'
                }

                stage('cat config file') {
                    sh 'cat cluster/aws/config.yaml'
                }

                stage('update generated aws config') {
                    sh "build-scripts/update-generated-config.sh cluster/aws/config.yaml ${env.JOB_BASE_NAME}-${env.BUILD_ID}"
                }

                stage("read config file again") {
                    sh 'cat cluster/aws/config.yaml'
                }

            }
            // container('k2-tools'){
            //
            //     stage('checkout') {
            //         checkout scm
            //     }
            //
            //     stage('aws config generation') {
            //         echo WORKSPACE
            //         sh "cd ${env.WORKSPACE} && ./k2cli generate"
            //
            //     }
            // }

        }
    }
