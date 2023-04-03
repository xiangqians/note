/**
 * @author xiangqian
 * @date 22:24 2022/12/25
 */
;
custom = function () {
    let obj = {}

    // 判断是否是 object
    obj.isObj = function (v) {
        return Object.prototype.toString.call(v) === '[object Object]'
    }

    // 判断是否是 string
    obj.isStr = function (v) {
        return Object.prototype.toString.call(v) === '[object String]'
    }

    /**
     * 人性化文件大小
     * @param size 文件大小，单位：byte
     */
    obj.humanizFileSize = function (size) {

        // B, Byte
        // 1B  = 8b
        // 1KB = 1024B
        // 1MB = 1024KB
        // 1GB = 1024MB
        // 1TB = 1024GB

        if (size <= 0) {
            return "0 B"
        }

        function format(num) {
            // Math.floor()，不四舍五入 ，向下取整
            return Math.floor(num * 100) / 100
        }

        // GB
        let gb = size / (1024 * 1024 * 1024)
        if (gb > 1) {
            return format(gb) + ' GB'
        }

        // MB
        let mb = size / (1024 * 1024)
        if (mb > 1) {
            return format(mb) + ' MB'
        }

        // KB
        let kb = size / 1024
        if (kb > 1) {
            return format(kb) + ' KB'
        }

        // B
        return size + ' B'
    }

    // 获取指定范围的随机数 [m, n)
    obj.random = function (m, n) {
        return Math.round(Math.random() * (n - m) + m)
    }

    // 添加url时间戳
    obj.addTimestamp = function (url) {
        let timestamp = new Date().getTime()
        timestamp += obj.random(-1000, 1000)
        custom.random(1, 1000)
        if (url.indexOf('?') > 0) {
            url += '&t=' + timestamp
        } else {
            url += '?t=' + timestamp
        }
        return url
    }

    obj.reload = function (text) {
        // 使用 document.write() 覆盖当前文档
        document.write(text)
        document.close()

        // 修改当前浏览器地址
        let $html = $('html')
        let url = $html.attr('uri')
        if (url) {
            history.replaceState(undefined, undefined, url)
        }
    }

    // ------------------------------ ajax ------------------------------

    /**
     * ajax
     * @param url       请求url
     * @param method    请求方法：GET, POST, PUT, DELETE
     * @param data      请求数据
     * @param dataType  请求数据类型，form, json
     * @param async     是否异步请求
     * @param success   请求成功回调函数
     * @param error     请求错误回调函数
     */
    obj.ajax = function (url, method, data, dataType, async, success, error) {
        url = obj.addTimestamp(url)

        let param = {
            url: url,
            type: method,
            data: data,
            async: async,
            success: function (resp) {
                if (success) {
                    success(resp)
                }
            },
            error: function (resp) {
                if (error) {
                    error(resp)
                }
            }
        }

        // form
        // application/x-www-form-urlencoded
        if (dataType === "form") {
            param.processData = false
            param.contentType = false
        }
        // json
        else if (dataType === "json") {
            param.dataType = "json"
        }
        // other
        else {
            return
        }

        $.ajax(param)
    }

    obj.ajaxDefault = function (url, method, data, dataType, async) {
        custom.ajax(url, method, data, dataType, async, function (data) {
            custom.reload(data)
        }, function (e) {
            alert(JSON.stringify(e))
            console.error(e)
        })
    }

    // form, a, button
    obj.ajaxE = function ($e) {
        // form
        if ($e.is('form')) {
            let $form = $e
            // console.log($form)
            $($form.find("[type=submit]")[0]).click(function () {
                let action = $form.attr("action")
                // console.log('action', action)
                let method = $form.attr("method").trim().toUpperCase()
                // console.log('method', method)
                let data = new FormData()
                $form.serializeArray().forEach(e => {
                    // console.log(e.name, e.value)
                    data.append(e.name, e.value);
                })

                // file
                let $input = $form.find("input[type='file']");
                if ($input.length > 0) {
                    let files = $input[0].files;
                    // console.log($input.attr('name'), files);
                    if (files.length > 0) {
                        data.append($input.attr('name'), files[0]);
                    }
                }

                // console.log('data', data);
                // data.forEach((value, key) => {
                //     console.log(key, value);
                // })

                custom.ajaxDefault(action, method, data, 'form', true)
                return false
            })
            return
        }

        // a
        // // 处理 ajaxE
        //     $ajaxE.click(function () {
        //         let data = null
        //
        //         // pre
        //         let pre = $ajaxE[0].pre
        //         if (pre) {
        //             let r = pre($ajaxE)
        //             // console.log('r', r)
        //             if (custom.isObj(r)) {
        //                 data = new FormData()
        //                 for (let name in r) {
        //                     data.append(name, r[name])
        //                 }
        //
        //             } else {
        //                 return
        //             }
        //         }
        //         // confirm
        //         else {
        //             let message = $ajaxE.attr("confirm")
        //             if (message) {
        //                 if (!confirm(message)) {
        //                     return
        //                 }
        //             }
        //         }
        //
        //         // ajaxFormData
        //         let url = $ajaxE.attr("href")
        //         // console.log(url)
        //         if (!(url)) {
        //             url = $ajaxE.attr("action")
        //         }
        //         if (!(url)) {
        //             url = $ajaxE.attr("url")
        //         }
        //         // console.log(url)
        //
        //         let method = $ajaxE.attr("method").trim().toUpperCase()
        //
        //         custom.ajaxDefault(url, method, data, "form", true)
        //
        //         // 如果是 <a></a> 标签，则取消 <a></a> 默认行为
        //         return false
        //     })
    }

    // ------------------------------ storage ------------------------------

    // 存储
    function Storage() {
    }

    Storage.prototype.vToStr = function (v) {
        if (!obj.isStr(v)) {
            v = JSON.stringify(v)
        }
        return v
    }

    // 存储数据
    Storage.prototype.set = function (key, value) {
        window.localStorage.setItem(this.vToStr(key), this.vToStr(value));
        return true
    }

    // 获取数据
    Storage.prototype.get = function (key) {
        return window.localStorage.getItem(this.vToStr(key));
    }

    // 删除数据
    Storage.prototype.remove = function (key) {
        window.localStorage.removeItem(this.vToStr(key))
        return true
    }

    obj.storage = new Storage()


    // ------------------------------ pie-chart ------------------------------

    obj.PieChart = function (canvas) {
        this.canvas = canvas
        this.cxt = canvas.getContext('2d')
        this.w = this.cxt.canvas.width
        this.h = this.cxt.canvas.height
        this.x = this.w / 2 + 30
        this.y = this.h / 2
        this.r = 150
        this.line = 20
        this.rectW = 30
        this.rectH = 15
        this.rectL = 10
        this.rectT = 6
    }

    obj.PieChart.prototype.beginPath = function () {
        this.cxt.beginPath()
    }

    obj.PieChart.prototype.getColor = function () {
        function color() {
            let min = 0, max = 255
            return parseInt(Math.random() * (max - min + 1) + min)
        }

        return `rgb(${color()},${color()},${color()})`
    }

    obj.PieChart.prototype.drawArc = function (sAngle, eAngle, color) {
        this.cxt.moveTo(this.x, this.y)
        this.cxt.arc(this.x, this.y, this.r, sAngle, eAngle)
        this.cxt.fillStyle = color
        this.cxt.fill()
    }

    obj.PieChart.prototype.drawLabelDetails = function (sAngle, angle, color, labelDetails) {
        this.beginPath()
        this.endX = Math.cos(sAngle + angle / 2) * (this.r + this.line) + this.x
        this.endY = Math.sin(sAngle + angle / 2) * (this.r + this.line) + this.y
        this.cxt.moveTo(this.x, this.y)
        this.cxt.lineTo(this.endX, this.endY)
        this.cxt.strokeStyle = color
        this.cxt.stroke()

        this.beginPath()
        this.textWidth = this.cxt.measureText(labelDetails).width
        this.cxt.moveTo(this.endX, this.endY)
        this.lineEndX = this.endX > this.x ? this.endX + this.textWidth : this.endX - this.textWidth
        this.cxt.lineTo(this.lineEndX, this.endY)
        this.cxt.strokeStyle = color
        this.cxt.stroke()

        // 绘制标题
        this.beginPath()
        this.cxt.textBaseline = 'bottom'
        this.cxt.fillText(labelDetails, this.x > this.endX ? this.lineEndX : this.endX, this.endY)
    }

    obj.PieChart.prototype.drawRect = function (label, n, rectColor) {
        this.beginPath()
        let rectEndT = this.rectT * (n + 1) + this.rectH * (n)
        this.cxt.fillRect(this.rectL, rectEndT, this.rectW, this.rectH)
        // 配套相应的文字
        this.cxt.font = '12px Miscrosoft Yahei'
        this.cxt.textBaseline = 'middle'
        this.cxt.fillText(label, this.rectL + this.rectW + this.rectT, rectEndT + this.rectH / 2)
        this.cxt.fillStyle = rectColor
        this.cxt.fill()
    }

    /**
     * draw
     * @param data array数据
     * @param getLabel 获取标签
     * @param getLabelDetails 获取标签详情
     * @param getNum 获取数量
     */
    obj.PieChart.prototype.draw = function (data, getLabel, getLabelDetails, getNum) {
        if (!(data) || !(data.length)) {
            return
        }

        // total
        let total = 0
        data.forEach(e => total += getNum(e))

        // percentage（百分比）
        data.forEach(e => {
            e._pct = getNum(e) / total * Math.PI * 2
        })

        let start = 0
        let end = 0
        for (let i = 0, len = data.length; i < len; i++) {
            let color = this.getColor()
            this.beginPath()
            if (i == 0) {
                start = 0
                end = data[i]._pct
            } else {
                start += data[i - 1]._pct
                end += data[i]._pct
            }

            // 绘制弧
            this.drawArc(start, end, color)
            // 绘制标签详情
            this.drawLabelDetails(start, data[i]._pct, color, getLabelDetails(data[i]))

            // 绘制左上角标签
            this.drawRect(getLabel(data[i]), i, color)
        }
    }

    return obj
}()
;

(function () {
    ;

    // 为普通的 <a></a> url添加时间戳
    // let $aArr = $('a:not([ajaxE])[href^="/"]')
    // // console.log('$aArr', $aArr)
    // for (let i = 0, len = $aArr.length; i < len; i++) {
    //     let $a = $($aArr[i])
    //     // console.log($a)
    //     let href = $a.attr('href')
    //     $a.attr('href', custom.addTimestamp(href))
    // }
    //
    // // 为普通的 <img></img> src添加时间戳
    // let $imgArr = $('img')
    // for (let i = 0, len = $imgArr.length; i < len; i++) {
    //     let $img = $($imgArr[i])
    //     let src = $img.attr('src')
    //     $img.attr('src', custom.addTimestamp(src))
    // }

    // div收缩/展开
    let $divs = $("div[class='float']")
    for (let i = 0; i < $divs.length; i++) {
        let $div = $($divs[i])
        // console.log($floatDiv)
        let $btn = $('<button value="+">+</button>')
        let $wrapperDiv = $('<div></div>')
        // 默认隐藏div
        $wrapperDiv.css('display', 'none')
        $wrapperDiv.html($div.html())
        $div.html('')
        $btn.click(function () {
            let value = $btn.attr('value')
            // 设置为 +
            if (value === '-') {
                $btn.attr('value', '+')
                $btn.text('+')
                // 隐藏div
                $wrapperDiv.css('display', 'none')
            }
            // 设置为 -
            else {
                $btn.attr('value', '-')
                $btn.text('-')
                // 显示div
                $wrapperDiv.css('display', 'block')
            }
        })
        $div.prepend($btn)
        $div.append($wrapperDiv)
    }

})()
;
