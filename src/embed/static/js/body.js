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