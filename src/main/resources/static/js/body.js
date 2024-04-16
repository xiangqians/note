/**
 * @author xiangqian
 * @date 22:21 2024/03/02
 */
;

/**
 * Http
 *
 * AJAX = 异步 JavaScript 和 XML（Asynchronous JavaScript and XML）。
 * 在不重载整个网页的情况下，AJAX 通过后台加载数据，并在网页上进行显示。
 *
 * @param url {string} 请求地址
 * @param type {string} 请求类型（方法）：GET | POST | PUT | DELETE
 * @param formData {FormData}
 * @param success {function} function (data, status, xhr) {}
 * @param error {function} function (xhr, status, error) {}
 */
function http(url, type, formData, success, error) {
    // $.ajax(url[, options])
    // $.ajax([options])
    $.ajax({
        // 服务器端地址
        url: url,

        // 请求方法：GET | POST | PUT | DELETE
        type: type,

        // 请求数据
        data: formData,

        // contentType 发送到服务器的数据类型。
        // contentType:"form"，发送FormData数据。
        // 不设置 contentType:"application/json"，数据是以键值对的形式传递到服务器（data: {name: "test"}）；
        // 设置  contentType:"application/json"，数据是以json串的形式传递到后端（data: '{name: "test"}'），如果传递的是比较复杂的数据（例如多层嵌套数据），这时候就需要设置 contentType:"application/json" 了。
        // contentType: contentType,

        // 是否异步请求，true-异步，false-同步。默认 true
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
        success: success || function (data, status, xhr) {
        },

        // 请求错误回调函数
        error: error || function (xhr, status, error) {
            alert(`${status}, ${error}`)
        }
    })
}

function showElement($element) {
    $element.css('display', 'inline-block')
}

function hideElement($element) {
    $element.css('display', 'none')
}

$(function () {
    let $trs = $('body > content > left > table tbody tr')

    // 鼠标进入<tr>时执行的函数
    $trs.mouseenter(function (event) {
        // 获取目标<tr>元素
        let $tr = $(event.currentTarget)

        // 获取目标<tr>元素中的最后一个<td>元素
        let $lastTd = $($tr.find('td:last'))

        // 显示<td>元素
        $lastTd.css('visibility', 'visible')
    })

    // 鼠标离开<tr>时执行的函数
    $trs.mouseleave(function (event) {
        // 获取目标<tr>元素
        let $tr = $(event.currentTarget)

        // 获取目标<tr>元素中的最后一个<td>元素
        let $lastTd = $($tr.find('td:last'))

        // 隐藏<td>元素
        $lastTd.css('visibility', 'hidden')
    })
})

;