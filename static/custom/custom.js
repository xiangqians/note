/**
 * @author xiangqian
 * @date 22:24 2022/12/25
 */
;custom = function () {

    // object
    let obj = {}

    /**
     * 判断是否属于某类型
     * @param obj object
     * @param type '[object {type}]'
     * @returns {boolean}
     */
    obj.isType = function (obj, type) {
        return Object.prototype.toString.call(obj) === type
    }

    /**
     * 判断是否是 undefined
     * @param obj
     * @returns {boolean}
     */
    obj.isUndefined = function (obj) {
        return typeof (obj) === 'undefined'
    }
    /**
     * 判断是否是 Object
     * @param obj
     * @returns {boolean}
     */
    obj.isObj = function (obj) {
        return obj.isType(obj, '[object Object]')
    }

    /**
     * 判断是否是 String
     * @param _obj
     * @returns {boolean}
     */
    obj.isStr = function (_obj) {
        return obj.isType(_obj, '[object String]')
    }

    /**
     * 判断是否是 FormData
     * @param obj
     * @returns {boolean}
     */
    obj.isFormData = function (obj) {
        return obj.isType(obj, '[object FormData]')
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

    /**
     * 获取指定范围的随机数 [min, max)
     * @param min
     * @param max
     * @returns {number}
     */
    obj.random = function (min, max) {
        return Math.round(Math.random() * (max - min) + min)
    }

    /**
     * rgb
     * @returns {`rgb(${string},${string},${string})`}
     */
    obj.rgb = function () {
        function color() {
            return obj.random(0, 255 + 1)
        }

        return `rgb(${color()},${color()},${color()})`
    }

    /**
     * 解析url查询参数
     * @returns {Object}
     */
    obj.parseUrlQueryParams = function (url) {
        let params = {}
        let index = url.indexOf("?")
        if (index === -1) {
            return params
        }

        let arr = url.substring(index + 1).split("&")
        for (let i = 0; i < arr.length; i++) {
            let e = arr[i]
            let eArr = e.split("=")
            let name = eArr[0]
            let value = null
            if (eArr.length == 2) {
                value = eArr[1]
            }
            params[name] = value
        }
        return params
    }

    /**
     * 时间戳
     * @returns {number}
     */
    obj.timestamp = function () {
        return new Date().getTime() + obj.random(-1000, 1000)
    }

    /**
     * ajax
     * AJAX = 异步 JavaScript 和 XML（Asynchronous JavaScript and XML）。
     * 在不重载整个网页的情况下，AJAX 通过后台加载数据，并在网页上进行显示。
     * @param url       服务器端地址
     * @param method    请求方法：GET | POST | PUT | DELETE
     * @param data      {FormData} 请求数据
     * @param async     是否异步请求，true，异步；false，同步。默认 true
     * @param success   请求成功回调函数
     * @param error     请求错误回调函数
     */
    obj.ajax = function (url, method, data, async, complete, success, error) {
        // url
        let timestamp = obj.timestamp()
        if (url.indexOf('?') > 0) {
            url += '&t=' + timestamp
        } else {
            url += '?t=' + timestamp
        }

        if (obj.isUndefined(async)) {
            async = false
        }

        if (obj.isUndefined(complete)) {
            // 指定当HTTP请求结束时（请求成功或请求失败的回调函数，此时已经运行完毕）的回调函数。不管请求成功或失败，该回调函数都会执行。它的参数为发出请求的原始对象以及返回的状态信息。
            complete = function (xhr, status) {
            }
        }

        if (obj.isUndefined(success)) {
            // 请求成功回调函数
            success = function (data, status, xhr) {
                // 重新加载document
                // 使用 document.write() 覆盖当前文档
                document.write(data)
                document.close()

                // 修改当前浏览器地址
                let $html = $('html')
                let url = $html.attr('url')
                if (url) {
                    history.replaceState(undefined, undefined, url)
                }
            }
        }

        if (obj.isUndefined(error)) {
            // 请求错误回调函数
            error = function (xhr, status, error) {
                console.error(status, error)
                error = JSON.stringify(error)
                alert(`${status}\n${error}`)
            }
        }

        let params = {
            url: url, type: method, data: data,

            // contentType 发送到服务器的数据类型。
            // contentType:"form"，发送FormData数据。
            // 不设置 contentType:"application/json"，数据是以键值对的形式传递到服务器（data: {name: "test"}）；
            // 设置  contentType:"application/json"，数据是以json串的形式传递到后端（data: '{name: "test"}'），如果传递的是比较复杂的数据（例如多层嵌套数据），这时候就需要设置 contentType:"application/json" 了。
            // contentType: contentType,

            async: async,

            // 预期服务器返回的数据类型，当设置dataType："json"时，如果后端返回了String，则ajax无法执行，去掉后ajax会自动检测返回数据类型。可以设为 text、html、script、json、jsonp和xml，和form
            // dataType: dataType,

            // 等待的最长毫秒数。如果过了这个时间，请求还没有返回，则自动将请求状态改为失败。
            // timeout: 30 * 1000,

            // 浏览器是否缓存服务器返回的数据，默认为true，注：浏览器本身不会缓存POST请求返回的数据，所以即使设为false，也只对HEAD和GET请求有效。
            // cache: cache,

            // 指定发出请求前，所要调用的函数，通常用来对发出的数据进行修改。
            // beforeSend: beforeSend,

            // 指定当HTTP请求结束时（请求成功或请求失败的回调函数，此时已经运行完毕）的回调函数。不管请求成功或失败，该回调函数都会执行。它的参数为发出请求的原始对象以及返回的状态信息。
            complete: complete,

            // 请求成功回调函数
            success: success,

            // 请求错误回调函数
            error: error
        }

        // contentType
        // application/x-www-form-urlencoded
        // 不处理发送数据
        params.processData = false
        // 不设置Content-Type请求头
        params.contentType = false

        // $.ajax(url[, options])
        // $.ajax([options])
        $.ajax(params)
    }

    /**
     * element点击事件
     * form, a, button, ...
     * @param $e
     * @param callback
     */
    obj.clickE = function ($e, callback) {
        // form
        if ($e.is('form')) {
            let $form = $e
            // console.log($form)
            $($form.find("[type=submit]")[0]).click(function () {
                // url
                let url = $form.attr("action")
                // console.log('url', url)

                // method
                let method = $form.attr("method").trim().toUpperCase()
                // console.log('method', method)

                // data
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

                // ajax
                obj.ajax(url, method, data)
                return false
            })
            return
        }

        // e
        $e.click(function () {
            let params = null
            if (callback) {
                params = callback($e)
            }

            //  url & method 属性必须存在
            if (!params || !params.hasOwnProperty('url') || !params.hasOwnProperty('method')) {
                return false
            }

            let url = params.url
            let method = params.method

            // data
            let data = null
            if (params.hasOwnProperty('data')) {
                data = params.data
            }
            console.log(method, url, data)

            // ajax
            obj.ajax(url, method, data)

            // 取消 <a></a> 默认行为
            return false
        })
    }


    return obj
}()

;(function (obj) {
})(custom)


;(function (obj) {
})(custom)

// storage
;(function (obj) {

    /**
     * 存储
     * @constructor
     */
    function Storage() {
    }

    /**
     * 存储数据
     * @param key {string}
     * @param value {object}
     * @returns {boolean}
     */
    Storage.prototype.set = function (key, value) {
        window.localStorage.setItem(key, value ? JSON.stringify(value) : null);
        return true
    }

    /**
     * 获取数据
     * @param key {string}
     * @returns {object}
     */
    Storage.prototype.get = function (key) {
        let value = window.localStorage.getItem(key)
        if (value) {
            return JSON.parse(value)
        }
        return null;
    }

    /**
     * 删除数据
     * @param key {string}
     * @returns {boolean}
     */
    Storage.prototype.del = function (key) {
        window.localStorage.removeItem(key)
        return true
    }

    // new storage
    obj.storage = new Storage()

})(custom)


// float 收缩/展开
;(function (obj) {

    obj.float = function (type, data) {
        if (!(type === 'note' || type === 'img')) {
            return;
        }

        console.log(type, data)

        const key = 'float'

        // 设置float显示
        function setFloatDisplay() {
            custom.storage.set(key, {value: 'block'})
        }

        // 设置float隐藏
        function setFloatHide() {
            custom.storage.set(key, {value: 'none'})
        }

        function getFloat() {
            let value = custom.storage.get(key)
            return value ? value.value : 'none'
        }

        // float是否显示
        function isFloatDisplay() {
            return getFloat() === 'block'
        }

        // btn
        let $btn = $('<button></button>')

        // 显示btn
        function displayBtn() {
            $btn.attr('value', '-')
            $btn.text('-')
            $btn.css('margin-top', '10px')
        }

        // 隐藏btn
        function hideBtn() {
            $btn.attr('value', '+')
            $btn.text('+')
            $btn.css('margin-top', '0px')
        }

        // wrapper div
        let $wrapperDiv = $('<div style="padding: 20px"></div>')

        // 显示div
        function displayWrapperDiv() {
            $wrapperDiv.css('display', 'block')
            setFloatDisplay()
            displayBtn()
        }

        // 隐藏div
        function hideWrapperDiv() {
            $wrapperDiv.css('display', 'none')
            setFloatHide()
            hideBtn()
        }

        $btn.click(function () {
            let value = $btn.attr('value')
            if (value === '-') {
                hideWrapperDiv()
            } else {
                displayWrapperDiv()
            }
        })

        let $float = $($('div[class="float"]')[0])
        $wrapperDiv.html($float.html())
        $float.html('')
        $float.prepend($btn)
        $float.append($wrapperDiv)

        if (isFloatDisplay()) {
            displayWrapperDiv()
        } else {
            hideWrapperDiv()
        }

        // form
        let uploadUrl = ''
        if (type === 'note') {
            uploadUrl = '/note/upload'
        } else if (type === 'img') {
            uploadUrl = '/img/upload'
        }
        uploadUrl += `?t=${obj.timestamp()}`
        let $form = $($float.find('form')[0])
        $form.attr('action', uploadUrl)

        return;

        // hist
        let isHist = url.indexOf('hist') > 0
        // console.log(isHist)
        let idx = -1
        if (isHist) {
            let idxStr = url.substring(url.indexOf('hist/') + 'hist/'.length, url.indexOf('/view'))
            idx = parseInt(idxStr)
        }

        // 如果不是历史记录，则移除form的 disabled 属性
        if (!isHist) {
            $($float.find("input[name='file']")[0]).attr('disabled', false)
            $($float.find("button[type='submit']")[0]).attr('disabled', false)
        }

        // form
        custom.ajaxE($($float.find('form')[0]))

        return

        // select
        let $select = $($("select[name='hist']")[0])
        let options = $select.find("option")
        if (idx != -1) {
            for (let i = 0; i < options.length; i++) {
                let $option = $(options[i])
                let value = parseInt($option.attr('value'))
                // console.log(value)
                if (value === idx) {
                    $option.attr('selected', true)
                    break
                }
            }
        }

        $select.change(function () {
            let value = $select.find("option:selected").attr("value");
            // console.log(value)
            let url = null
            if (value === "-1") {
                url = '/img/{{ $data.Id }}/view'
            } else {
                url = '/img/{{ $data.Id }}/hist/' + value + '/view'
            }
            custom.ajaxReload(url, 'GET', null)
        });

    }

})(custom)

;
