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
 * 存储数据
 * @param key {string}
 * @param value {string}
 */
Storage.prototype.set = function (key, value) {
    window.localStorage.setItem(key, value)
}

/**
 * 获取数据
 * @param key {string}
 * @returns {string}
 */
Storage.prototype.get = function (key) {
    return window.localStorage.getItem(key)
}

/**
 * 删除数据
 * @param key {string}
 */
Storage.prototype.del = function (key) {
    window.localStorage.removeItem(key)
}

storage = new Storage()

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