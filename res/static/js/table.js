/**
 * @author xiangqian
 * @date 22:35 2023/10/15
 */
;
(function () {

    let $body = $($('body')[0])

    // 右键菜单列表
    let $ul = $('<ul class="menu"></ul>')
    $ul.append($('<li>111</li>'))
    $ul.append($('<li>222</li>'))
    $ul.append($('<li>333</li>'))
    $body.prepend($ul)

    // 为tr添加右键菜单事件
    let $trs = $('table tbody tr[name!="noData"]')
    for (let i = 0, len = $trs.length; i < len; i++) {
        let $tr = $($trs[i])
        $tr.contextmenu(function (e) {
            // 阻止默认行为
            e.preventDefault()

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
            $ul.css('display', 'block')

        });
    }

    // 监听 html->body 鼠标点击事件
    $body.click(function (e) {
        // 隐藏菜单
        $ul.css('display', 'none')
    })

    $ul.click(function (e) {
        let tagName = e.target.tagName.toLowerCase()
        if (tagName !== 'li') {
            return
        }
        
        let $li = $(e.target)
        console.log($li)
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


