pipeline {
    agent any

    environment {
        // You must set the following environment variables
        // ORGANIZATION_NAME
        YOUR_DOCKERHUB_USERNAME = "morettimathieu"
        
        SERVICE_NAME = "hexagun"
        REPOSITORY_TAG="${YOUR_DOCKERHUB_USERNAME}/${SERVICE_NAME}:${BUILD_ID}"
    }

    stages {
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
                                            --destination "${REPOSITORY_TAG}" 
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
