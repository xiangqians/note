/**
 * @author xiangqian
 * @date 22:21 2024/03/02
 */
;

/**
 * HTTP 请求
 * @param type
 * @param url
 * @param contentType
 * @param data
 * @param success
 * @param error
 */
function http(type, url, contentType, data, success, error) {
    // AJAX = 异步 JavaScript 和 XML（Asynchronous JavaScript and XML）。
    // 在不重载整个网页的情况下，AJAX 通过后台加载数据，并在网页上进行显示。
    $.ajax({
        // 请求类型：GET、POST、PUT、DELETE
        type: type,

        // 服务器端地址
        url: url,

        // 发送到服务器的数据类型
        // 1、application/x-www-form-urlencoded（默认）
        //    这是 jQuery 默认的 contentType，适用于普通的表单数据。在这种情况下，数据会以 key1=value1&key2=value2 的形式进行编码。
        // 2、application/json
        //    适用于发送 JSON 格式的数据。通常用于RESTful API，发送的数据会被序列化为 JSON 格式。
        // 3、multipart/form-data
        //    用于发送文件（如图片、视频等）。通常与 FormData 一起使用，浏览器会自动设置 contentType 为 multipart/form-data，因此通常不需要手动设置此选项。
        // 4、text/plain
        //    用于发送纯文本数据。这种类型会将请求体发送为原始的文本格式，不做编码。
        contentType: contentType,

        // 不处理发送数据
        processData: false,

        // 请求数据
        data: data,

        // 是否异步请求，true，异步；false，同步。默认 true。
        async: false,

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
            alert(status + " " + error)
        }
    });
};
