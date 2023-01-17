<!DOCTYPE html>
<html lang="zh-CN">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no">
    <meta name="csrf-token" content="">
    <link rel="shortcut icon" type="image/png" href="/favicon.ico">
    <title>Not Found - ApiCat</title>
    <meta name="keywords" content="API开发,API文档,API协作,API调试,API测试,API模拟,API Mock,API导入导出,API对接,文档管理,团队协作,项目协作">
    <meta name="description" content="ApiCat是一款API协作开发提效软件，其简化了团队成员间的协作流程，提供了优质的API文档、API Mock、数据文件导入导出等功能，让开发者可以更快更好的完成开发工作。">
    <meta name="renderer" content="webkit|ie-comp|ie-stand">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <link rel="stylesheet" href="/static/stylesheet/common.css">
</head>

<body class="bg-white">
    <header class="header">
        <div class="container">
            <a class="logo" href="/">
                <img src="/static/image/logo.svg" alt="ApiCat">
                <span class="logo">
                    <span class="logo-text logo-apicat">ApiCat</span></span>
            </a>
        </div>
    </header>

    <main class="container">
        <div class="text-center">
            <img class="not-found-image" src="/static/image/404@2x.png" alt="">
            <p class="not-found-tip">
                啊哦，网页走丢了，正在努力寻找中……<a href="/home" id="home">回到首页</a>
            </p>
        </div>
    </main>

    <script>
        (function() {
            var isLogin = localStorage.getItem("api.cat.token") || null;
            var backHomeBtn = document.getElementById("home");
            backHomeBtn &&
                backHomeBtn.setAttribute(
                    "href",
                    isLogin ? "/home" : "/login"
                );
        })();
    </script>

</body>

</html>
