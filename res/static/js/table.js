/**
 * @author xiangqian
 * @date 22:35 2023/10/15
 */
;
(function () {

    let $body = $($('body')[0])

    // 右键菜单列表
    let $ul = $('<ul></ul>')
    $ul.append($('<a href="#"><li>重命名</li></a>'))
    $ul.append($('<a href="#"><li>删除</li></a>'))
    $ul.append($('<a href="#"><li>恢复</li></a>'))
    $ul.append($('<a href="#"><li>永久删除</li></a>'))

    let $div = $('<div></div>')
    $div.append($ul)
    $body.prepend($div)

    let $targetTr = null

    // 为tr添加右键菜单事件
    let $trs = $('table tbody tr[class!="noData"]')
    for (let i = 0, len = $trs.length; i < len; i++) {
        let $tr = $($trs[i])
        $tr.contextmenu(function (e) {
            // 阻止默认行为
            e.preventDefault()

            // 设置选中的tr背景颜色
            let $target = $(e.target)
            if (e.target.tagName.toLowerCase() !== 'tr') {
                $targetTr = $target.parent('tr')
            }
            $targetTr.css('background-color', '#CCCCCC')

            // 禁止页面滚动
            $body.css('overflow', 'hidden')

            // 设置菜单位置
            let x = e.clientX
            let y = e.clientY
            if (x >= document.documentElement.clientWidth - $ul.offsetWidth) {
                x = document.documentElement.clientWidth - $ul.offsetWidth
            }
            if (y >= document.documentElement.clientHeight - $ul.offsetHeight) {
                y = document.documentElement.clientHeight - $ul.offsetHeight
            }
            $ul.css('left', `${x}px`)
            $ul.css('top', `${y}px`)

            // 显示菜单
            $div.css('display', 'block')
            $ul.css('display', 'block')
        });
    }

    // 监听 html->body 鼠标点击事件
    $body.click(function (e) {
        // 重置已设置选中的tr背景颜色
        if ($targetTr != null) {
            $targetTr.css('background-color', '')
        }

        // 允许页面滚动
        $body.css('overflow', '')

        // 隐藏菜单
        $ul.css('display', 'none')
        $div.css('display', 'none')
    })

    // 菜单点击事件
    $ul.click(function (e) {
        let tagName = e.target.tagName.toLowerCase()
        // console.log(e.target)
        if (tagName !== 'li') {
            return
        }

        let $li = $(e.target)
        console.log($li[0])
    })

    return

    document.addEventListener("contextmenu", function (evt) {
        evt.preventDefault()
        list.style.display = "block"
        var x = evt.clientX
        var y = evt.clientY
        if (x >= document.documentElement.clientWidth - list.offsetWidth)
            x = document.documentElement.clientWidth - list.offsetWidth
        if (y >= document.documentElement.clientHeight - list.offsetHeight)
            y = document.documentElement.clientHeight - list.offsetHeight
        list.style.left = x + "px"
        list.style.top = y + "px"
    })
    document.addEventListener("click", (e) => {
        list.style.display = "none"
    })
    list.onclick = function (evt) {
        console.log(evt.target)
        if (evt.target.className === "aaa") {
            console.log(111)
        }

    }
})();


