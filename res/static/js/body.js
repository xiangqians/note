/**
 * @author xiangqian
 * @date 14:32 2023/10/14
 */
;

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