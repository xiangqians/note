/**
 * @author xiangqian
 * @date 22:24 2022/12/25
 */
;
custom = function () {

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
     * @param obj
     * @returns {boolean}
     */
    obj.isStr = function (obj) {
        return obj.isType(obj, '[object String]')
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
        let timestamp = new Date().getTime()
        timestamp += obj.random(-1000, 1000)
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
            url: url,
            type: method,
            data: data,

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

            // 是否中断
            if (params && params.abort) {
                return false
            }

            // url
            let url = null
            if (params && params.hasOwnProperty('url')) {
                url = params.url
            }
            // <a><a/>
            else {
                url = $e.attr("href")
            }
            // console.log('url', url)

            // method
            let method = null
            if (params && params.hasOwnProperty('method')) {
                method = params.method
            } else {
                method = $e.attr("method").trim().toUpperCase()
            }
            // console.log('method', method)

            // data
            let data = null
            if (params && params.hasOwnProperty('data')) {
                data = params.data
            }
            // console.log('data', data)

            // ajax
            obj.ajax(url, method, data)

            // 取消 <a></a> 默认行为
            return false
        })
    }

    /**
     * 存储
     * @constructor
     */
    function Storage() {
    }

    /**
     * value to string
     * @param v
     * @returns {string}
     * @private
     */
    Storage.prototype._vToStr = function (v) {
        if (!obj.isStr(v)) {
            v = JSON.stringify(v)
        }
        return v
    }

    /**
     * 存储数据
     * @param key
     * @param value
     * @returns {boolean}
     */
    Storage.prototype.set = function (key, value) {
        window.localStorage.setItem(this._vToStr(key), this._vToStr(value));
        return true
    }

    /**
     * 获取数据
     * @param key
     * @returns {string}
     */
    Storage.prototype.get = function (key) {
        return window.localStorage.getItem(this._vToStr(key));
    }

    /**
     * 删除数据
     * @param key
     * @returns {boolean}
     */
    Storage.prototype.remove = function (key) {
        window.localStorage.removeItem(this._vToStr(key))
        return true
    }

    // new storage
    obj.storage = new Storage()

    return obj
}()

// float
;(function (obj) {

    // float 收缩/展开
    const key = 'displayFloat'
    let $divs = $('div[class="float"]')
    for (let i = 0; i < $divs.length; i++) {
        let $div = $($divs[i])
        // console.log($div)
        let display = custom.storage.get(key)
        if (!display) {
            display = 'none' // 隐藏
        }
        let $btn = $('<button></button>')

        function setBtn(value) {
            $btn.attr('value', value)
            $btn.text(value)
            if (value === '+') {
                $btn.css('margin-top', '0px')
            } else {
                $btn.css('margin-top', '10px')
            }
        }

        let value = display === 'none' ? '+' : '-'
        setBtn(value)

        let $wrapperDiv = $('<div style="padding: 20px"></div>')
        $wrapperDiv.css('display', display)
        $wrapperDiv.html($div.html())
        $div.html('')
        $btn.click(function () {
            let value = $btn.attr('value')
            // 设置为 +
            if (value === '-') {
                setBtn('+')
                // 隐藏div
                $wrapperDiv.css('display', 'none')
                custom.storage.set(key, 'none')
            }
            // 设置为 -
            else {
                setBtn('-')
                // 显示div
                $wrapperDiv.css('display', 'block')
                custom.storage.set(key, 'block')
            }
        })
        $div.prepend($btn)
        $div.append($wrapperDiv)
    }

    // float
    obj.float = function (url, entity, type) {
        let $float = $($('div[class="float"]')[0])

        entity = JSON.parse(entity)
        // console.log(entity)

        // pdf
        if (entity.type === 'pdf') {
            // params
            let params = custom.urlQueryParams(decodeURIComponent(url))
            // console.log(params)

            // v
            let version = '2.0'
            if (params.v) {
                version = params.v
            }

            // select
            let $select = $($float.find("select[name='version']")[0])
            let options = $select.find("option")
            for (let i = 0; i < options.length; i++) {
                let $option = $(options[i])
                if ($option.attr('value') === version) {
                    $option.attr('selected', true)
                    break
                }
            }

            $select.change(function () {
                let version = $select.find("option:selected").attr("value");
                // console.log(version)
                let url = `/${type}/${entity.id}/view?v=${version}`
                custom.ajaxReload(url, 'GET', null)
            });
        }

        let uploadUrl = null

        // note
        if (type === 'note') {
            // path
            let $path = $($float.find('td[name="path"]')[0])
            $path.html(entity.pathLink)
            $($float.find('tr[name="path"]')[0]).css('display', '')

            // upload
            uploadUrl = '/note/upload'
        }
        // img
        else if (type === 'img') {
            // upload
            uploadUrl = '/img/upload'
        }

        // upload form
        let $form = $($float.find('form')[0])
        $form.attr('action', uploadUrl)

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
