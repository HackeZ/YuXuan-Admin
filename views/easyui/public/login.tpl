{{template "../public/header.tpl"}}

<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <title>登录 - 御轩寝室</title>
    <link href="../../static/css/vendors.css" rel="stylesheet">
    <link href="../../static/css/user_login.css" rel="stylesheet">
    <style type="text/css">
        #particles-js {
            position: absolute;
            width: 100%;
            height: 100%;
            background-color: rgb(247, 250, 252);
            background-image: url("");
            background-repeat: no-repeat;
            background-size: cover;
            background-position: 50% 50%;
            z-index: -1;
        }
        
        .login-tips {
            margin-top: 26px;
            margin-bottom: -30px;
            color: red;
        }
    </style>
</head>

<body>

    <div class="navbar">
        <div class="nav">
            <ul>
                <li>
                    <a href="##">前台首页</a>
                    <span class="cursor"></span>
                </li>
            </ul>
        </div>
    </div>
    <div id="particles-js"></div>
    <div class="user-container">
        <div class="user-container-title">
            <h2>御轩寝室</h2>
            <h3>后台管理登录</h3>
        </div>
        <div class="form-container" id="form-container">

            <form id="form" method="post" autocomplete="off">
                <div class="form-item">
                    <label for="username">用户名或邮箱</label>
                    <input id="id_username" maxlength="30" name="username" required="required" type="text" />
                </div>
                <div class="form-item">
                    <label for="password">密码</label>
                    <input id="id_password" name="password" required="required" type="password" />
                </div>
                <div class="form-login-tips">
                    <p class="login-tips"></p>
                </div>
                <div class="form-item">
                    <input type="button" name="login" id="login-btn" value="登录" onclick="fromsubmit()">
                </div>
                <div class="register-tab">
                    <a href="##">御轩寝室 - 出品</a>
                </div>
            </form>
        </div>
    </div>
    <script src="../../static/js/vendors.js"></script>
    <script src="http://cdn.bootcss.com/particles.js/2.0.0/particles.min.js"></script>
    <script src="../../static/js/user_login.js"></script>
    <script type="text/javascript">
        particlesJS.load("particles-js", "../../static/assets/particles.json", function() {
            console.log('callback - particles.js config loaded');
        });
    </script>
    <script type="text/javascript">
        var URL = "/public"

        function fromsubmit() {
            // $("#form").submit(function() {

            // });
            $.ajax({
                type: "POST",
                url: URL + '/login?isajax=1',
                data: $("#form").serialize(),
                async: false,
                error: function() {
                    $(".login-tips").text("网络出问题了，请联系管理员！");
                    $("#id_username").val("").focus();
                    $("#id_password").val("");
                },
                success: function(data) {
                    if (data.status) {
                        location.href = URL + "/index"
                    } else {
                        $(".login-tips").text(data.info);
                        $("#id_username").val("").focus();
                        $("#id_password").val("");
                    }
                }
            });
        }
        //这个就是键盘触发的函数
        var SubmitOrHidden = function(evt) {
            evt = window.event || evt;
            if (evt.keyCode == 13) { //如果取到的键值是回车
                fromsubmit();
            }
        }
        window.document.onkeydown = SubmitOrHidden; //当有键按下时执行函数
    </script>
</body>

</html>