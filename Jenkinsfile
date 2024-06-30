def dev_repository_tag=""

pipeline {
    agent any

    environment {
        // You must set the following environment variables
        // ORGANIZATION_NAME
        YOUR_DOCKERHUB_USERNAME = "morettimathieu"
        
        SERVICE_NAME = "hexagun"
        DOCKER_REPOSITORY = "${YOUR_DOCKERHUB_USERNAME}/${SERVICE_NAME}"
        PROD_REPOSITORY_TAG = "${DOCKER_REPOSITORY}:${BUILD_ID}-prod"
    }
    stages {
        stage('Checkout and Versioning') {
            options { skipDefaultCheckout(true)}
            steps {
                echo "**** scm.branches is ${scm.branches} ****"
                checkout(
                  [ $class: 'GitSCM',
                    branches: scm.branches, // Assumes the multibranch pipeline checkout branch definition is sufficient
                    // extensions: [
                    //   [ $class: 'CloneOption', shallow: true, depth: 1, honorRefspec: true, noTags: true, reference: '/var/lib/git/mwaite/bugs/jenkins-bugs.git'],
                    //   [ $class: 'LocalBranch', localBranch: env.BRANCH_NAME ],
                    //   [ $class: 'PruneStaleBranch' ]
                    // ],
                    // Add honor refspec and reference repo for speed and space improvement
                    gitTool: scm.gitTool,
                    // Default is missing narrow refspec
                    userRemoteConfigs: [ [ url: scm.userRemoteConfigs[0].url ] ]
                    // userRemoteConfigs: scm.userRemoteConfigs // Assumes the multibranch pipeline checkout remoteconfig is sufficient
                  ]
                )                
                script{
                    sh '''
                    ls -la
                    git status
                    VERSION=$(git describe --tags --abbrev=8)
                    dev_repository_tag="${DOCKER_REPOSITORY}:${VERSION%%-*}-${VERSION##*-}"
                '''
                    echo dev_repository_tag
                }
            }
        }     
        stage('Build') {
            agent {
                kubernetes {
                  yaml '''
                    apiVersion: v1
                    kind: Pod
                    metadata:
                        labels:
                        some-label: some-label-value
                    spec:
                        containers:
                        - name: go-builder
                          image: golang:1.22
                          command:
                          - cat
                          tty: true                       
                    '''
                }
            }
            steps {   
                container('go-builder') {

                    // Output will be something like "go version go1.22 darwin/arm64"
                    sh 'ls -la'
                    sh 'go version'
                    sh 'go env -w GOFLAGS="-buildvcs=false"'
                    sh 'go mod download'
                    sh 'go build -o server'
                }                
            }
        }
        stage('Build and Push Image') {
            agent {
                kubernetes {
                  yaml '''
                    apiVersion: v1
                    kind: Pod
                    metadata:
                        labels:
                        some-label: some-label-value
                    spec:
                        containers:
                        - name: kaniko
                          image: gcr.io/kaniko-project/executor:debug
                          imagePullPolicy: Always
                          command:
                          - cat
                          tty: true
                          volumeMounts:      
                            - name: kaniko-secret
                              mountPath: /kaniko/.docker
                        volumes:
                        - name: kaniko-secret
                          secret:
                            secretName: kaniko-secret
                    '''
                }
            }
            steps {
                container('kaniko') {
                    script {
                        if (env.BRANCH_NAME == 'main') {
                            sh '''
                            /kaniko/executor --dockerfile `pwd`/Dockerfile      \
                                            --context `pwd`                    \
                                            --destination "${PROD_REPOSITORY_TAG}" 
                            '''
                        } else if (env.BRANCH_NAME =~ /^dev.*/ ) {
                            sh '''
                            /kaniko/executor --dockerfile `pwd`/Dockerfile      \
                                            --context `pwd`                    \
                                            --destination "${dev_repository_tag}" 
                            '''
                        } else {
                            sh '''
                            /kaniko/executor --dockerfile `pwd`/Dockerfile      \
                                            --context `pwd`                    \
                                            --no-push 
                            '''
                        }
                    }
                }
            }
        }    
    }
}
