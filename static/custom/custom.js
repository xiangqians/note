/**
 * @author xiangqian
 * @date 22:24 2022/12/25
 */
;
custom = function () {

    let _obj = {}

    /**
     * 判断是否属于某类型
     * @param obj object
     * @param type '[object {type}]'
     * @returns {boolean}
     */
    _obj.isType = function (obj, type) {
        return Object.prototype.toString.call(obj) === type
    }

    /**
     * 判断是否是 undefined
     * @param obj
     * @returns {boolean}
     */
    _obj.isUndefined = function (obj) {
        return typeof (obj) === 'undefined'
    }
    /**
     * 判断是否是 Object
     * @param obj
     * @returns {boolean}
     */
    _obj.isObj = function (obj) {
        return _obj.isType(obj, '[object Object]')
    }

    /**
     * 判断是否是 String
     * @param obj
     * @returns {boolean}
     */
    _obj.isStr = function (obj) {
        return _obj.isType(obj, '[object String]')
    }

    /**
     * 判断是否是 FormData
     * @param obj
     * @returns {boolean}
     */
    _obj.isFormData = function (obj) {
        return _obj.isType(obj, '[object FormData]')
    }

    /**
     * 人性化文件大小
     * @param size 文件大小，单位：byte
     */
    _obj.humanizFileSize = function (size) {

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
    _obj.random = function (min, max) {
        return Math.round(Math.random() * (max - min) + min)
    }

    _obj.rgb = function () {
        function color() {
            return _obj.random(0, 255 + 1)
        }

        return `rgb(${color()},${color()},${color()})`
    }

    /**
     * 给url添加时间戳
     * @param url
     * @returns {string}
     */
    _obj.addUrlTimestamp = function (url) {
        let timestamp = new Date().getTime()
        timestamp += _obj.random(-1000, 1000)
        if (url.indexOf('?') > 0) {
            url += '&t=' + timestamp
        } else {
            url += '?t=' + timestamp
        }
        return url
    }

    /**
     * 重新加载document
     * @param text
     */
    _obj.reload = function (text) {
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


    /**
     * 获取url查询参数
     * @returns {Object}
     */
    _obj.urlQueryParams = function (url) {
        let params = {}

        // 获取url中"?"符后的字串
        if (!url) {
            url = location.search
        }
        // console.log(search)

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
     * @param data      请求数据
     * @param contentType 发送到服务器的数据类型。
     *      contentType:"form"，发送FormData数据。
     *      不设置 contentType:"application/json"，数据是以键值对的形式传递到服务器（data: {name: "test"}）；
     *      设置  contentType:"application/json"，数据是以json串的形式传递到后端（data: '{name: "test"}'），如果传递的是比较复杂的数据（例如多层嵌套数据），这时候就需要设置 contentType:"application/json" 了。
     * @param async     是否异步请求，true，异步；false，同步。默认 true
     * @param dataType   预期服务器返回的数据类型，当设置dataType："json"时，如果后端返回了String，则ajax无法执行，去掉后ajax会自动检测返回数据类型。可以设为 text、html、script、json、jsonp和xml，和form
     * @param success   请求成功回调函数
     * @param error     请求错误回调函数
     */
    _obj.ajax = function (url, method, data, contentType, async, dataType, success, error) {
        url = _obj.addUrlTimestamp(url)
        let params = {
            url: url,
            type: method,
            data: data,
            async: async,
            timeout: 30 * 1000, // 等待的最长毫秒数。如果过了这个时间，请求还没有返回，则自动将请求状态改为失败。
            // cache: cache, // 浏览器是否缓存服务器返回的数据，默认为true，注：浏览器本身不会缓存POST请求返回的数据，所以即使设为false，也只对HEAD和GET请求有效。
            // beforeSend: beforeSend, // 指定发出请求前，所要调用的函数，通常用来对发出的数据进行修改。
            // complete: complete, // 指定当HTTP请求结束时（请求成功或请求失败的回调函数，此时已经运行完毕）的回调函数。不管请求成功或失败，该回调函数都会执行。它的参数为发出请求的原始对象以及返回的状态信息。
            success: success,
            error: error
        }

        // contentType
        if (contentType) {
            // application/x-www-form-urlencoded
            if (contentType === 'form') {
                // 不处理发送数据
                params.processData = false
                // 不设置Content-Type请求头
                params.contentType = false
            }
            // other
            else {
                params.contentType = contentType
            }
        }

        // dataType
        if (dataType) {
            params.dataType = dataType
        }

        // $.ajax(url[, options])
        // $.ajax([options])
        $.ajax(params)
    }

    _obj.ajaxSimple = function (url, method, data, success, error) {
        let contentType = null
        if (_obj.isFormData(data)) {
            contentType = 'form'
        }
        _obj.ajax(url, method, data, contentType, false, null, success, error)
    }

    _obj.ajaxReload = function (url, method, data) {
        _obj.ajaxSimple(url, method, data, function (resp) {
            _obj.reload(resp)
        }, function (e) {
            console.error(e)
            alert(e)
        })
    }

    // ajaxE 无所事事
    _obj.ajaxEDoNothing = {
        url: null,
        method: null
    }

    // form, a, button, ...
    _obj.ajaxE = function ($e, callback) {
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
                _obj.ajaxReload(url, method, data)
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

            // url
            let url = null
            if (params && params.hasOwnProperty('url')) {
                url = params.url
            } else {
                // <a><a/>
                url = $e.attr("href")
            }
            // console.log('url', url)
            if (!url) {
                return false
            }

            // method
            let method = null
            if (params && params.hasOwnProperty('method')) {
                method = params.method
            } else {
                method = $e.attr("method").trim().toUpperCase()
            }
            // console.log('method', method)
            if (!method) {
                return false
            }

            // data
            let data = null
            if (params && params.hasOwnProperty('data')) {
                data = params.data
            }
            // console.log('data', data)

            // ajax
            _obj.ajaxReload(url, method, data)

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
        if (!_obj.isStr(v)) {
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
    _obj.storage = new Storage()

    return _obj
}()

;(function (_obj) {

    _obj.note = function (note, uri) {
        note = JSON.parse(note)
        // console.log(note)

        // pdf
        if (note.type === 'pdf') {
            // params
            let params = custom.urlQueryParams(decodeURIComponent(uri))
            // console.log(params)

            // v
            let version = '2.0'
            if (params.v) {
                version = params.v
            }

            // select
            let $select = $($("select[name='version']")[0])
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
                let url = `/note/${note.id}/view?v=${version}`
                custom.ajaxReload(url, 'GET', null)
            });
        }

        // path
        let $path = $($('td[name="path"]')[0])
        $path.html(note.pathLink)

    }

})(custom)

;(function () {
    ;

    // div收缩/展开
    let $divs = $("div[class='float']")
    for (let i = 0; i < $divs.length; i++) {
        let $div = $($divs[i])
        // console.log($div)
        let display = custom.storage.get('displayFloat')
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
                custom.storage.set('displayFloat', 'none')
            }
            // 设置为 -
            else {
                setBtn('-')
                // 显示div
                $wrapperDiv.css('display', 'block')
                custom.storage.set('displayFloat', 'block')
            }
        })
        $div.prepend($btn)
        $div.append($wrapperDiv)
    }

})()
;
