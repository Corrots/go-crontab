<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .title }}</title>
    <!-- bootstrap + jquery -->

    <!-- vuejs  , reactjs , angular -->
    <script src="https://cdn.bootcss.com/jquery/3.3.1/jquery.min.js"></script>
    <link href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
</head>
<body>
<div class="container-fluid">
    <!-- 页头 -->
    <div class="row">
        <div class="col-md-12">
            <div class="page-header">
                <h1>管理后台<small>Golang分布式Crontab</small></h1>
            </div>
        </div>
    </div>

    <!-- 功能按钮 -->
    <div class="row">
        <div class="col-md-12">
            <button type="button" class="btn btn-primary" id="new-job">新建任务</button>
            <button type="button" class="btn btn-success" id="list-worker">健康节点</button>
        </div>
    </div>

    <!-- 任务列表 -->
    <div class="row">
        <div class="col-md-12">
            <div class="panel panel-default" style="margin-top: 20px">
                <div class="panel-body">
                    <table id="job-list"  class="table table-striped">
                        <thead>
                        <tr>
                            <th>任务名称</th>
                            <th>shell命令</th>
                            <th>cron表达式</th>
                            <th>任务操作</th>
                        </tr>
                        </thead>
                        <tbody>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</div>

<!-- position: fixed -->
<div id="edit-modal" class="modal fade" tabindex="-1" role="dialog">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title">编辑任务</h4>
            </div>
            <div class="modal-body">
                <form>
                    <div class="form-group">
                        <label for="edit-name">任务名称</label>
                        <input type="text" class="form-control" id="edit-name" placeholder="任务名称">
                    </div>
                    <div class="form-group">
                        <label for="edit-command">shell命令</label>
                        <input type="text" class="form-control" id="edit-command" placeholder="shell命令">
                    </div>
                    <div class="form-group">
                        <label for="edit-cronExpr">cron表达式</label>
                        <input type="text" class="form-control" id="edit-cronExpr" placeholder="cron表达式">
                    </div>
                </form>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
                <button type="button" class="btn btn-primary" id="save-job">保存</button>
            </div>
        </div><!-- /.modal-content -->
    </div><!-- /.modal-dialog -->
</div><!-- /.modal -->

<!--  日志模态框 -->
<div id="log-modal" class="modal fade" tabindex="-1" role="dialog">
    <div class="modal-dialog modal-lg" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title">任务日志</h4>
            </div>
            <div class="modal-body">
                <table id="log-list" class="table table-striped">
                    <thead>
                    <tr>
                        <th>shell命令</th>
                        <th>错误原因</th>
                        <th>脚本输出</th>
                        <th>计划开始时间</th>
                        <th>实际调度时间</th>
                        <th>开始执行时间</th>
                        <th>执行结束时间</th>
                    </tr>
                    </thead>
                    <tbody>

                    </tbody>
                </table>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
            </div>
        </div><!-- /.modal-content -->
    </div><!-- /.modal-dialog -->
</div><!-- /.modal -->

<!--  健康节点模态框 -->
<div id="worker-modal" class="modal fade" tabindex="-1" role="dialog">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title">健康节点</h4>
            </div>
            <div class="modal-body">
                <table id="worker-list" class="table table-striped">
                    <thead>
                    <tr>
                        <th>节点IP</th>
                    </tr>
                    </thead>
                    <tbody>

                    </tbody>
                </table>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
            </div>
        </div><!-- /.modal-content -->
    </div><!-- /.modal-dialog -->
</div><!-- /.modal -->

<script>
        // 页面加载完成后, 回调函数
        $(document).ready(function() {
            // 时间格式化函数
            function timeFormat(millsecond) {
                // 前缀补0: 2018-08-07 08:01:03.345
                function paddingNum(num, n) {
                    var len = num.toString().length
                    while (len < n) {
                        num = '0' + num
                        len++
                    }
                    return num
                }
                var date = new Date(millsecond)
                var year = date.getFullYear()
                var month = paddingNum(date.getMonth() + 1, 2)
                var day = paddingNum(date.getDate(), 2)
                var hour = paddingNum(date.getHours(), 2)
                var minute = paddingNum(date.getMinutes(), 2)
                var second = paddingNum(date.getSeconds(), 2)
                var millsecond = paddingNum(date.getMilliseconds(), 3)
                return year + "-" + month + "-" + day + " " + hour + ":" + minute + ":" + second + "." + millsecond
            }

            // 1, 绑定按钮的事件处理函数
            // 用javascript委托机制, DOM事件冒泡的一个关键原理

            // 编辑任务
            $("#job-list").on("click", ".edit-job", function(event) {
                // 取当前job的信息，赋值给模态框的input
                $('#edit-name').val($(this).parents('tr').children('.job-name').text())
                $('#edit-command').val($(this).parents('tr').children('.job-command').text())
                $('#edit-cronExpr').val($(this).parents('tr').children('.job-cronExpr').text())
                // 弹出模态框
                $('#edit-modal').modal('show')
            })
            // 删除任务
            $("#job-list").on("click", ".delete-job", function(event) { // javascript bind
                var jobName = $(this).parents("tr").children(".job-name").text()
                $.ajax({
                    url: '/job/' + jobName,
                    type: 'DELETE',
                    dataType: 'json',
                    // data: {name: jobName},
                    complete: function() {
                        window.location.reload()
                    }
                })
            })
            // 杀死任务
            $("#job-list").on("click", ".kill-job", function(event) {
                var jobName = $(this).parents("tr").children(".job-name").text()
                $.ajax({
                    url: '/job/kill' + jobName,
                    type: 'post',
                    dataType: 'json',
                    // data: {name: jobName},
                    complete: function() {
                        window.location.reload()
                    }
                })
            })
            // 保存任务
            $('#save-job').on('click', function() {
                var jobInfo = {name: $('#edit-name').val(), command: $('#edit-command').val(), expression: $('#edit-cronExpr').val()}
                $.ajax({
                    url: '/job/save',
                    type: 'post',
                    contentType: "application/json; charset=utf-8",
                    dataType: 'json',
                    data: JSON.stringify(jobInfo),
                    complete: function() {
                        window.location.reload()
                    }
                })
            })
            // 新建任务
            $('#new-job').on('click', function() {
                $('#edit-name').val("")
                $('#edit-command').val("")
                $('#edit-cronExpr').val("")
                $('#edit-modal').modal('show')
            })
            // 查看任务日志
            $("#job-list").on("click", ".log-job", function(event) {
                // 清空日志列表
                $('#log-list tbody').empty()

                // 获取任务名
                var jobName = $(this).parents('tr').children('.job-name').text()

                // 请求/job/log接口
                $.ajax({
                    url: "/job/log",
                    dataType: 'json',
                    data: {name: jobName},
                    success: function(resp) {
                        if (resp.code != "00") {
                            return
                        }
                        // 遍历日志
                        var logList = resp.data
                        for (var i = 0; i < logList.length; ++i) {
                            var log = logList[i]
                            var tr = $('<tr>')
                            tr.append($('<td>').html(log.command))
                            tr.append($('<td>').html(log.err))
                            tr.append($('<td>').html(log.output))
                            tr.append($('<td>').html(timeFormat(log.planTime)))
                            tr.append($('<td>').html(timeFormat(log.scheduleTime)))
                            tr.append($('<td>').html(timeFormat(log.startTime)))
                            tr.append($('<td>').html(timeFormat(log.endTime)))
                            console.log(tr)
                            $('#log-list tbody').append(tr)
                        }
                    }
                })

                // 弹出模态框
                $('#log-modal').modal('show')
            })

            // 健康节点按钮
            $('#list-worker').on('click', function() {
                // 清空现有table
                $('#worker-list tbody').empty()

                // 拉取节点
                $.ajax({
                    url: '/worker/list',
                    dataType: 'json',
                    success: function(resp) {
                        if (resp.code != "00") {
                            return
                        }

                        var workerList = resp.data
                        // 遍历每个IP, 添加到模态框的table中
                        for (var i = 0; i < workerList.length; ++i) {
                            var workerIP = workerList[i]
                            var tr = $('<tr>')
                            tr.append($('<td>').html(workerIP))
                            $('#worker-list tbody').append(tr)
                        }
                    }
                })

                // 弹出模态框
                $('#worker-modal').modal('show')
            })

            // 2，定义一个函数，用于刷新任务列表
            function rebuildJobList() {
                // /job/list
                $.ajax({
                    url: '/job/list',
                    dataType: 'json',
                    success: function(resp) {
                        if (resp.code != "00") {  // 服务端出错了
                            return
                        }
                        // 任务数组
                        var jobList = resp.data
                        // 清理列表
                        $('#job-list tbody').empty()
                        // 遍历任务, 填充table
                        for (var i = 0; i < jobList.length; ++i) {
                            var job = jobList[i];
                            var tr = $("<tr>")
                            tr.append($('<td class="job-name">').html(job.name))
                            tr.append($('<td class="job-command">').html(job.command))
                            tr.append($('<td class="job-cronExpr">').html(job.expression))
                            var toolbar = $('<div class="btn-toolbar">')
                                    .append('<button class="btn btn-info edit-job">编辑</button>')
                                    .append('<button class="btn btn-danger delete-job">删除</button>')
                                    .append('<button class="btn btn-warning kill-job">强杀</button>')
                                    .append('<button class="btn btn-success log-job">日志</button>')
                            tr.append($('<td>').append(toolbar))
                            $("#job-list tbody").append(tr)
                        }
                    }
                })
            }
            rebuildJobList()
        })
    </script>

</body>
</html>