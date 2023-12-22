/**
 * @author xiangqian
 * @date 14:32 2023/10/14
 */
;

/**
 * 存储
 * @constructor
 */
function Storage() {
}

/**
 * 存储字符串数据
 * @param key {string}
 * @param value {string}
 */
Storage.prototype.setString = function (key, value) {
    if (value != null) {
        value = (value + "").trim()
    }
    window.localStorage.setItem(key, value)
}

/**
 * 获取字符串数据
 * @param key {string}
 * @returns {string | null}
 */
Storage.prototype.getString = function (key) {
    let value = window.localStorage.getItem(key)
    if (value == null || value === '') {
        return null
    }
    return value
}


/**
 * 存储对象数据
 * @param key {string}
 * @param value {object}
 */
Storage.prototype.setObject = function (key, value) {
    if (value != null) {
        value = JSON.stringify(value)
    }
    this.setString(key, value)
}

/**
 * 获取对象数据
 * @param key {string}
 * @returns {object}
 */
Storage.prototype.getObject = function (key) {
    let value = this.getString(key)
    if (value == null) {
        return null
    }
    return JSON.parse(value)
}

/**
 * 删除数据
 * @param key {string}
 */
Storage.prototype.del = function (key) {
    window.localStorage.removeItem(key)
}

storage = new Storage()

function Http() {
}

/**
 * Http POST
 *
 * AJAX = 异步 JavaScript 和 XML（Asynchronous JavaScript and XML）。
 * 在不重载整个网页的情况下，AJAX 通过后台加载数据，并在网页上进行显示。
 *
 * @param url
 * @param data {}
 * @param success function (data, status, xhr) {}
 * @param error function (xhr, status, error) {}
 */
Http.post = function (url, data, success, error) {
    let formData = new FormData()
    for (let key in data) {
        formData.append(key, data[key])
    }

    // $.ajax(url[, options])
    // $.ajax([options])
    $.ajax({
        // 服务器端地址
        url: url,

        // 请求方法：GET | POST | PUT | DELETE
        type: 'POST',

        // 请求数据
        data: formData,

        // contentType 发送到服务器的数据类型。
        // contentType:"form"，发送FormData数据。
        // 不设置 contentType:"application/json"，数据是以键值对的形式传递到服务器（data: {name: "test"}）；
        // 设置  contentType:"application/json"，数据是以json串的形式传递到后端（data: '{name: "test"}'），如果传递的是比较复杂的数据（例如多层嵌套数据），这时候就需要设置 contentType:"application/json" 了。
        // contentType: contentType,

        // 是否异步请求，true，异步；false，同步。默认 true
        async: false,

        // contentType
        // application/x-www-form-urlencoded
        // 不处理发送数据
        processData: false,

        // 不设置Content-Type请求头
        contentType: false,

        // 预期服务器返回的数据类型，当设置dataType："json"时，如果后端返回了String，则ajax无法执行，去掉后ajax会自动检测返回数据类型。可以设为 text、html、script、json、jsonp和xml，和form
        // dataType: dataType,

        // 等待的最长毫秒数。如果过了这个时间，请求还没有返回，则自动将请求状态改为失败。
        // timeout: 30 * 1000,

        // 浏览器是否缓存服务器返回的数据，默认为true，注：浏览器本身不会缓存POST请求返回的数据，所以即使设为false，也只对HEAD和GET请求有效。
        // cache: cache,

        // 指定发出请求前，所要调用的函数，通常用来对发出的数据进行修改。
        // beforeSend: beforeSend,

        // 指定当HTTP请求结束时（请求成功或请求失败的回调函数，此时已经运行完毕）的回调函数。不管请求成功或失败，该回调函数都会执行。它的参数为发出请求的原始对象以及返回的状态信息。
        complete: function (xhr, status) {
        },

        // 请求成功回调函数
        success: success,

        // 请求错误回调函数
        error: error
    })
}

// 是否是隐藏元素
function hasHideElement($element) {
    return $element.hasClass('hide')
}

// 显示元素
function displayElements() {
    if (arguments.length === 0) {
        return
    }

    for (let i = 0, len = arguments.length; i < len; i++) {
        let $element = arguments[i]
        $element.removeClass('hide')
    }
}

// 隐藏元素
function hideElements() {
    if (arguments.length === 0) {
        return
    }

    for (let i = 0, len = arguments.length; i < len; i++) {
        let $element = arguments[i]
        $element.addClass('hide')
    }
}

;