/**
 * @author xiangqian
 * @date 22:24 2022/12/25
 */
;
Custom = function () {
    let obj = {}

    obj.reload = function (text) {
        // 使用 document.write() 覆盖当前文档
        document.write(text)
        document.close()

        // 修改当前浏览器地址
        let $html = $('html')
        let url = $html.attr('url')
        if (url) {
            history.replaceState(undefined, undefined, url)
        }
    }

    obj.addTimestamp = function (url) {
        let timestamp = new Date().getTime()
        if (url.indexOf('?') > 0) {
            url += '&t=' + timestamp
        } else {
            url += '?t=' + timestamp
        }
        return url
    }

    obj.ajaxFormData = function (url, method, data, async, success, error) {
        url = obj.addTimestamp(url)

        // console.log(method, url, formData)
        // application/x-www-form-urlencoded
        $.ajax({
            url: url,
            type: method,
            data: data,
            processData: false,
            contentType: false,
            async: async,
            success: function (resp) {
                if (success) {
                    success(resp)
                }
            },
            error: function (e) {
                if (error) {
                    error(e)
                }
            }
        })
    }

    obj.ajaxJsonData = function (url, method, data, async, success, error) {
        url = obj.addTimestamp(url)
        $.ajax({
            url: url,
            type: method,
            data: data,
            dataType: "json",
            async: async,
            success: function (resp) {
                if (success) {
                    success(resp)
                }
            },
            error: function (e) {
                if (error) {
                    error(e)
                }
            }
        })
    }

    return obj
}()
;

(function () {
    function ajaxFormData(url, method, data) {
        Custom.ajaxFormData(url, method, data, true, function (resp) {
            Custom.reload(resp)
        }, function (e) {
            alert(e)
        })
    }

    let request = function ($e) {
        let formData = null
        let flag = true

        // pre func
        let pre = $e[0]._pre_
        if (pre) {
            let r = pre($e)
            let rarr = null
            let rl = 0
            if (Object.prototype.toString.call(r) === '[object Boolean]') {
                flag = r
            } else if (Object.prototype.toString.call(r) === '[object Array]' && (rl = (rarr = r).length) > 0) {
                flag = rarr[0] ? true : false
                if (flag && rl > 1) {
                    formData = new FormData()
                    for (let ri = 1; ri < rl; ri++) {
                        let robj = rarr[ri]
                        for (let name in robj) {
                            formData.append(name, robj[name])
                        }
                    }
                }
            }
        }
        // confirm
        else {
            let message = $e.attr("confirm")
            if (message) {
                flag = confirm(message)
            }
        }

        // ajaxFormData
        if (flag) {
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
            ajaxFormData(url, method, formData)
        }
    }

    // <a></a>
    // let $as = $('a[method]')
    let $as = $('a')
    for (let i = 0, l = $as.length; i < l; i++) {
        let $a = $($as[i])
        // console.log($a)
        let method = $a.attr('method')
        if (method) {
            $a.click(function () {
                request($a)

                // 取消 <a></a> 默认行为
                return false
            })
        }
        // 为 href 添加时间戳，预防被Chrome浏览器劫持，认为是”重复提交“。
        else {
            let href = $a.attr('href')
            if (href.indexOf("https://github.com") >= 0) {
                continue
            }

            let timestamp = new Date().getTime()
            if (href.indexOf('?') > 0) {
                href += '&t=' + timestamp
            } else {
                href += '?t=' + timestamp
            }
            $a.attr('href', href)
        }
    }

    // <form></form>
    let $formInputs = $('input[type="submit"]')
    for (let i = 0, l = $formInputs.length; i < l; i++) {
        let $input = $($formInputs[i])
        for (let $parent = $input.parent(); !$parent.is('body'); $parent = $parent.parent()) {
            if ($parent.is('form')) {
                let $form = $parent
                let method = $form.attr("method").trim().toUpperCase()
                if (method !== "POST") {
                    $input.click(function () {
                        let url = $form.attr("action")
                        let method = $form.attr("method").trim().toUpperCase()
                        let formData = new FormData()
                        $form.serializeArray().forEach(e => {
                            formData.append(e.name, e.value);
                        })
                        ajaxFormData(url, method, formData)
                        return false
                    })
                }
                break
            }
        }
    }

    // type=radio, method
    let $radioInputs = $('input[type="radio"][method]')
    for (let i = 0, l = $radioInputs.length; i < l; i++) {
        let $radioInput = $($radioInputs[i])
        $radioInput.click(function () {
            request($radioInput)
        })
    }

    // <select><option></option></select>
    let $selects = $('select[method]')
    for (let i = 0, l = $selects.length; i < l; i++) {
        let $select = $($selects[i])
        // console.log($select)
        $select.on('change', function () {
            // 获取选中的 <option></option>
            let $option = $($select.find('option:selected')[0]);
            // console.log($option)
            $select.attr('url', $option.attr('url'))
            request($select)
        })
    }

})()
;
