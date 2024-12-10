pipeline {
    agent any

    // 定义环境变量
    environment {
        // 例如设置项目相关的变量
        PROJECT_NAME = "OceanWang"
    }

    // 触发构建的条件，这里是当 GitHub 仓库有推送（push）事件时触发
    // triggers {
    //     githubPush()
    // }

    // 构建步骤
    stages {
        stage('Checkout') {
            steps {
                // 从 GitHub 仓库检出代码
                checkout([$class: 'GitSCM', 
                          branches: [[name: '*/main']], 
                          userRemoteConfigs: [[url: 'https://github.com/Zy1bREAd/Xdemo-backend.git']]])
            }
        }
        stage('Build') {
            steps {
                // 这里假设是一个基于 Java 的项目，执行构建命令，你需要根据自己项目类型修改
                sh'echo "test build - 1"' 
            }
        }
        stage('Test') {
            steps {
                // 运行测试用例，同样根据项目类型修改
                sh'echo "run test case - 2"' 
            }
        }
        stage('Deploy') {
            steps {
                // 部署步骤，例如将构建好的项目复制到目标服务器等操作，这里只是示例，需要完善
                sh 'echo "test CD - 3"'
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
