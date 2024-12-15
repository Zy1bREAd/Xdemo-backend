pipeline {
    agent any

    // 定义环境变量
    environment {
        // 例如设置项目相关的变量
        PROJECT_NAME = "OceanWang"
        CONTAINER_NAME = "xdemo_app"
        HARBOR_URL = "oceanwang.hub"
        HARBOR_PROJECT = "library"
        GITHUB_REPO_URL = "https://github.com/Zy1bREAd/Xdemo-backend.git"
        DEVELOP_SERVER_IP = "10.0.20.5"
        DEVELOP_SERVER_USER = "ubuntu"
        DEVELOP_SERVER_CRED_ID = "ssh-for-password-10.0.20.5"
    }

    // 构建步骤
    stages {
        stage('Checkout GitHub Branch and Pull Code') {
            steps {
                script {
                    if (env.GIT_BRANCH){
                        DOCKER_IMAGE_NAME = "xdemoapp"
                        DOCKER_IMAGE_TAG = "${env.GIT_BRANCH}"
                    }else if (env.GIT_TAG){
                        DOCKER_IMAGE_NAME = "xdemoapp"
                        DOCKER_IMAGE_TAG = "${env.GIT_TAG}"
                    }else {
                        echo "无法识别当前分支或标签"
                    }
                }
            }
        }
        // 登录Harbor镜像仓库
        stage('Login Image Registry') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'harbor_robot_account', passwordVariable: 'harbor_robot_token', usernameVariable: 'harbor_robot_account')]) {
                    sh "sudo docker login ${HARBOR_URL} -u ${harbor_robot_account} -p ${harbor_robot_token}"
                }
            }
        }
        // 构建镜像在dev环境
        stage('Build On Image For Develop') {
            when {
                anyOf {
                    expression { "${env.GIT_BRANCH}" =~ /(main)/ }
                    expression { "${env.GIT_TAG}" =~ /(^v.*)/ }
                }
            }
            steps {
                
                sh "sudo docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ."
            }
        }

        // 构建镜像在prod环境
        stage('Build On Image For Production') {
            when {
                anyOf {
                    expression { "${env.GIT_BRANCH}" =~ /(^release-v.*)/ }
                    expression { "${env.GIT_TAG}" =~ /(^release-v.*)/ }
                }
            }
            steps {
                
                sh "echo 'Enter Prod Pipeline for Build image.'"
            }
        }
        stage('Push Image') {
            // 推送镜像到Harbor
            steps {
                sh "sudo docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${HARBOR_URL}/${HARBOR_PROJECT}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"
                sh "sudo docker push ${HARBOR_URL}/${HARBOR_PROJECT}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"
            }
        }
        stage('Deploy To Develop Env') {
            when {
                anyOf {
                    expression { "${env.GIT_BRANCH}" =~ /(main)/ }
                    expression { "${env.GIT_TAG}" =~ /(^v.*)/ }
                }
            }
            steps {
                script {
                    def remote = [:]
                    remote.name = 'develop-server-01'
                    remote.host = "${DEVELOP_SERVER_IP}"
                    remote.allowAnyHosts = true
                    withCredentials([usernamePassword(credentialsId: 'harbor_robot_account', passwordVariable: 'harbor_robot_token', usernameVariable: 'harbor_robot_account'), usernamePassword(credentialsId: 'ssh-for-password-10.0.20.5', passwordVariable: 'dev_server_pwd', usernameVariable: 'dev_server_user')]) {
                        // 设置ssh server的login info
                        remote.user = "${dev_server_user}"
                        remote.password = "${dev_server_pwd}"
                        // 登录Harbor
                        sshCommand remote: remote, command: "sudo docker login ${HARBOR_URL} -u ${harbor_robot_account} -p ${harbor_robot_token}"
                        // 停止并删除之前的容器
                        sshCommand remote: remote, command: "if [ -n \"\$(sudo docker ps -a -q --filter name=${CONTAINER_NAME})\" ];then sudo docker stop ${CONTAINER_NAME} && sudo docker rm ${CONTAINER_NAME};else echo 'Container is not exist';fi"
                        
                        sshCommand remote: remote, command: "sudo docker run -itd --name=${CONTAINER_NAME} ${HARBOR_URL}/${HARBOR_PROJECT}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"
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