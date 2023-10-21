/**
 * @author xiangqian
 * @date 22:35 2023/10/15
 */
;
(function () {

    let $body = $($('body')[0])

    // 右键菜单列表
    let $ul = $('<ul></ul>')
    $ul.append($(`<a name="rename" href="#"><li>${i18n.rename}</li></a>`))
    $ul.append($(`<a name="cut" href="#"><li>${i18n.cut}</li></a>`))
    $ul.append($(`<a name="del" href="#"><li>${i18n.del}</li></a>`))
    $ul.append($(`<a name="restore" href="#"><li>${i18n.restore}</li></a>`))
    $ul.append($(`<a name="permlyDel" href="#"><li>${i18n.permlyDel}</li></a>`))

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

            // 根据删除状态显示右键菜单列表
            let $rename = $($ul.find('a[name="rename"]')[0])
            let $cut = $($ul.find('a[name="cut"]')[0])
            let $del = $($ul.find('a[name="del"]')[0])
            let $restore = $($ul.find('a[name="restore"]')[0])
            let $permlyDel = $($ul.find('a[name="permlyDel"]')[0])
            let del = $targetTr.attr('del')
            if (del === "0") {
                $rename.removeClass('hide')
                $cut.removeClass('hide')
                $del.removeClass('hide')
                $restore.addClass('hide')
                $permlyDel.addClass('hide')
            } else if (del === "1") {
                $rename.addClass('hide')
                $cut.addClass('hide')
                $del.addClass('hide')
                $restore.removeClass('hide')
                $permlyDel.removeClass('hide')
            }

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
        console.log($targetTr[0], $li[0])
    })

    // 重命名
    let $rename = $($ul.find('a[name="rename"]')[0])
    $rename.click(function (e) {
        let name = $targetTr.attr('name')
        name = prompt('{{ Localize "i18n.name" }}', name)
        if (!name || (name = name.trim()) === '') {
            // 取消 <a></a> 默认行为
            return false
        }

        let href = $rename.attr('href')
        // href = `?_method=PUT&name=${name}`
        // $rename.attr('href', href)
    })

})();


