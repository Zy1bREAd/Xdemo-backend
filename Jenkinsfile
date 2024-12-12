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
        // 判断分支来决定环境的CI/CD
        stage('Check Branch') {
            steps {
                script {
                    def branchName = env.BRANCH_NAME
                    if (branchName == 'main') {
                        developPipeline()
                    }else {
                        echo "不支持${branchName}分支构建"
                    }
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
    stage('Checkout GitHub and Pull Code') {
        steps {
            // 从 GitHub 仓库检出代码
            checkout([$class: 'GitSCM', 
                    branches: [[name: '*/main']], 
                    userRemoteConfigs: [[url: 'https://github.com/Zy1bREAd/Xdemo-backend.git']]])
            // 拉取代码
            git credentialsId: 'GitHub-Token', url: "${GITHUB_REPO_URL}"
        }
    }
    stage('Build On Image') {
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
    stage('Deploy To Server') {
        steps {
            // 运行测试用例，同样根据项目类型修改
            sh "echo Autodeploy"
            script {
                sshCommand remote: remote,command: "ls -alth"
            }
        }
    }
}

def productionPipeline() {
    stage('Not Support') {
        steps {
            sh 'echo "还没有支持正式环境的pipeline"'
        }
    }
}
