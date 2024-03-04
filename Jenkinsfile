pipeline {
  // Run on an agent where we want to use Go
  agent any

  // Ensure the desired Go version is installed for all stages,
  // using the name defined in the Global Tool Configuration
  tools { go '1.22.0' }

  stages {
    stage('Build') {
      steps {
        // Output will be something like "go version go1.22 darwin/arm64"
        sh 'go version'
      }
    }
  }
}