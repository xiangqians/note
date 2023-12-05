/**
 * @author xiangqian
 * @date 20:20 2023/12/05
 */
;
$(function () {
    let $float = $($('div[class*="float"]')[0])

    let $div = $('<div class="hide" style="padding: 10px;"></div>')
    $div.html($float.html())

    $float.html('')

    let $btn = $('<button></button>')
    $float.append($btn)
    $float.append($div)

    // 显示
    function display() {
        $btn.attr('value', '-')
        $btn.text('-')
        displayElements($div)
        storage.set('float', 'block')
    }

    // 隐藏
    function hide() {
        $btn.attr('value', '+')
        $btn.text('+')
        hideElements($div)
        storage.set('float', 'none')
    }

    if (storage.get('float') === 'none') {
        hide()
    } else {
        display()
    }
    displayElements($float)

    $btn.click(function () {
        let value = $btn.attr('value')
        if (value === '-') {
            hide()
        } else {
            display()
        }
    })
})
;
