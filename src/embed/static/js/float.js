/**
 * @author xiangqian
 * @date 20:20 2023/12/05
 */
;
$(function () {
    let $float = $($('div[class*="float"]')[0])

    let key = 'float'

    function hasHide() {
        return storage.getString(key) === 'hide'
    }

    function display() {
        storage.setString(key, '')
        displayElements($float)
    }

    function hide() {
        storage.setString(key, 'hide')
        hideElements($float)
    }

    if (!hasHide()) {
        display()
    }

    let $body = $($("body")[0])

    // 按下键盘事件
    $body.on('keydown', function (event) {
        // Ctrl + M，切换显示更多信息
        if (event.ctrlKey && event.key.toUpperCase() === 'M') {
            if (hasHide()) {
                display()
            } else {
                hide()
            }

            // 阻止默认行为
            // Prevent the default handler from running.
            event.preventDefault()
            return
        }
    })
})
;
