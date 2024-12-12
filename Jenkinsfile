pipeline {
    agent any

    // 定义环境变量
    environment {
        // 例如设置项目相关的变量
        PROJECT_NAME = "OceanWang"
        DOCKER_IMAGE_NAME = "xdemoapp"
        DOCKER_IMAGE_TAG = "main"
        HARBOR_URL = "124.220.17.5:8018"
        HARBOR_PROJECT = "xdemo"
        GITHUB_REPO_URL = "https://github.com/Zy1bREAd/Xdemo-backend.git"
        TARGET_SERVER_IP = "10.0.20.5"
        TARGET_SERVER_USER = "ubuntu"
    }

    // 触发构建的条件，这里是当 GitHub 仓库有推送（push）事件时触发
    // triggers {
    //     githubPush()
    // }

    // 构建步骤
    stages {
        stage('Checkout GitHub Branch and Pull Code') {
            steps {
                // 从 GitHub 仓库检出代码
                checkout([$class: 'GitSCM', 
                        branches: [[name: '*/main']], 
                        userRemoteConfigs: [[url: 'https://github.com/Zy1bREAd/Xdemo-backend.git']]])
            }
        }
        stage('Build On Image For Develop') {
            when {
                branch 'main'
            }
            steps {
                sh "docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ."
            }
        }
        stage('Build On Image For Production') {
            when {
                branch 'prod'
            }
            steps {
                sh "docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ."
            }
        }
        stage('Push Image') {
            // 推送镜像到Harbor
            steps {
                sh "docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${HARBOR_URL}/${HARBOR_PROJECT}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"
                sh "docker push ${HARBOR_URL}/${HARBOR_PROJECT}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"
            }
        }
        stage('Deploy To Develop Env') {
            when {
                branch 'main'
            }
            steps {
                // 运行测试用例，同样根据项目类型修改
                sh "echo Autodeploy"
                script {
                    sshCommand remote: remote,command: "ls -alth"
                }
            }
        }
        stage('Deploy To Production Env') {
            when {
                branch 'prod'
            }
            steps {
                // 运行测试用例，同样根据项目类型修改
                sh "echo Autodeploy"
                script {
                    sshCommand remote: remote,command: "ls -alth"
                }
            }
        }
    }

    // 构建后操作，如发送通知等
    // post {
    //     success {
    //         // 构建成功时发送邮件通知等操作，需要配置 Jenkins 的邮件插件等相关设置
    //         emailext subject: 'Build Success: ${PROJECT_NAME}', 
    //                 body: 'The build of ${PROJECT_NAME} was successful.', 
    //                 to: 'your-email@example.com'
    //     }
    //     failure {
    //         // 构建失败时发送邮件通知等操作
    //         emailext subject: 'Build Failure: ${PROJECT_NAME}', 
    //                 body: 'The build of ${PROJECT_NAME} has failed.', 
    //                 to: 'your-email@example.com'
    //     }
    // }
}

// 定义多分支流水线
def developPipeline() {
    
}
