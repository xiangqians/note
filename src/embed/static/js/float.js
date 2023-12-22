/**
 * @author xiangqian
 * @date 20:20 2023/12/05
 */
;
$(function () {
    let $float = $($('div[class*="float"]')[0])

    let key = 'float'

    function display() {
        storage.setString(key, 'display')
        displayElements($float)
    }

    function hide() {
        storage.setString(key, 'hide')
        hideElements($float)
    }

    function hasDisplay() {
        return storage.getString(key) === 'display'
    }

    if (hasDisplay()) {
        display()
    }

    let $body = $($("body")[0])

    // 按下键盘事件
    $body.on('keydown', function (event) {
        // Ctrl + M，切换显示更多信息
        if (event.ctrlKey && event.key.toUpperCase() === 'M') {
            if (hasDisplay()) {
                hide()
            } else {
                display()
            }

            // 阻止默认行为
            // Prevent the default handler from running.
            event.preventDefault()
            return
        }
    })
})
;
