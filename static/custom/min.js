/**
 * @author xiangqian
 * @date 22:24 2022/12/25
 */
;
custom = function () {
    let obj = {}

    // 获取指定范围的随机数 [m, n)
    obj.random = function (m, n) {
        return Math.round(Math.random() * (n - m) + m)
    }

    // url 添加时间戳
    obj.urlAddTimestamp = function (url) {
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
        url = obj.urlAddTimestamp(url)

        let param = {
            url: url,
            type: method,
            data: data,
            async: async,
            success: function (data) {
                if (success) {
                    success(data)
                }
            },
            error: function (e) {
                if (error) {
                    error(e)
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

    // 判断是否是 object
    obj.isObj = function (v) {
        return Object.prototype.toString.call(v) === '[object Object]'
    }

    // 判断是否是 string
    obj.isStr = function (v) {
        return Object.prototype.toString.call(v) === '[object String]'
    }

    obj.ajaxDefault = function (url, method, data, dataType, async) {
        obj.ajax(url, method, data, dataType, async, function (data) {
            obj.reload(data)
        }, function (e) {
            alert(JSON.stringify(e))
        })
    }

    obj.ajaxE = function ($e) {
        let data = null

        // pre
        let pre = $e[0].pre
        if (pre) {
            let r = pre($e)
            // console.log('r', r)
            if (obj.isObj(r)) {
                data = new FormData()
                for (let name in r) {
                    data.append(name, r[name])
                }

            } else {
                return
            }
        }
        // confirm
        else {
            let message = $e.attr("confirm")
            if (message) {
                if (!confirm(message)) {
                    return
                }
            }
        }

        // ajaxFormData
        let url = $e.attr("href")
        // console.log(url)
        if (!(url)) {
            url = $e.attr("action")
        }
        if (!(url)) {
            url = $e.attr("url")
        }
        // console.log(url)

        let method = $e.attr("method").trim().toUpperCase()

        obj.ajaxDefault(url, method, data, "form", true)
    }

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

    return obj
}()
;

(function () {

    // 处理 ajaxE
    let $ajaxEArr = $('[ajaxE]')
    // console.log('$ajaxEArr', $ajaxEArr)
    for (let i = 0, len = $ajaxEArr.length; i < len; i++) {
        let $ajaxE = $($ajaxEArr[i])
        // console.log($ajaxE)
        let tagName = $ajaxE.prop('tagName')
        // <form></form>
        if ((tagName.toLowerCase() === 'input' || tagName.toLowerCase() === 'button') && $ajaxE.attr('type') === 'submit') {
            let $input = $ajaxE
            for (let $parent = $input.parent(); !$parent.is('body'); $parent = $parent.parent()) {
                if ($parent.is('form')) {
                    let $form = $parent
                    $input.click(function () {
                        let action = $form.attr("action")
                        let method = $form.attr("method").trim().toUpperCase()
                        let data = new FormData()
                        $form.serializeArray().forEach(e => {
                            data.append(e.name, e.value);
                        })
                        custom.ajaxDefault(action, method, data, 'form', true)
                        return false
                    })
                    break
                }
            }
        } else {
            $ajaxE.click(function () {
                custom.ajaxE($ajaxE)

                // 如果是 <a></a> 标签，则取消 <a></a> 默认行为
                return false
            })
        }
    }

    // 为普通的 <a></a> url添加时间戳
    let $aArr = $('a:not([ajaxE])[href^="/"]')
    // console.log('$aArr', $aArr)
    for (let i = 0, len = $aArr.length; i < len; i++) {
        let $a = $($aArr[i])
        // console.log($a)
        let href = $a.attr('href')
        $a.attr('href', custom.urlAddTimestamp(href))
    }

})()
;
