/**
 * @author xiangqian
 * @date 22:35 2023/10/15
 */
;
$(function () {
    let $body = $($('body')[0])

    // 右键菜单列表
    let $ul = $('<ul></ul>')
    $ul.append($(`<a name="addFolder" href="#"><li>${i18n.addFolder}</li></a>`))
    $ul.append($(`<a name="addMdFile" href="#"><li>${i18n.addMdFile}</li></a>`))
    $ul.append($(`<a name="uploadFile" href="#"><li>${i18n.uploadFile}</li></a>`))
    $ul.append(`<li name="upload"><form action="#" method="post" enctype="multipart/form-data"><input name="file" type="file"/><button type="submit">${i18n.upload}</button></form></li>`)
    $ul.append($(`<a name="rename" href="#"><li>${i18n.rename}</li></a>`))
    $ul.append($(`<a name="cut" href="#"><li>${i18n.cut}</li></a>`))
    $ul.append($(`<a name="paste" href="#"><li>${i18n.paste}</li></a>`))
    $ul.append($(`<a name="del" href="#"><li>${i18n.del}</li></a>`))
    $ul.append($(`<a name="restore" href="#"><li>${i18n.restore}</li></a>`))
    $ul.append($(`<a name="permlyDel" href="#"><li>${i18n.permlyDel}</li></a>`))
    $ul.append($(`<a name="close" href="#"><li>${i18n.close}</li></a>`))
    let $div = $('<div class="menu"></div>')
    $div.append($ul)
    $body.prepend($div)

    let $targetTr = null

    let isUploadFile = false

    // 显示菜单
    function displayMenu() {
        // 禁止页面滚动
        $body.css('overflow', 'hidden')

        // 显示菜单
        $div.css('display', 'block')
        $ul.css('display', 'block')
    }

    // 隐藏菜单
    function hideMenu() {
        isUploadFile = false

        // 重置已设置选中的tr背景颜色
        if ($targetTr != null) {
            $targetTr.css('background-color', '')
            $targetTr = null
        }

        // 允许页面滚动
        $body.css('overflow', '')

        // 隐藏菜单
        $ul.css('display', 'none')
        $div.css('display', 'none')
    }

    let $addFolder = $($ul.find('a[name="addFolder"]')[0])
    let $addMdFile = $($ul.find('a[name="addMdFile"]')[0])
    let $uploadFile = $($ul.find('a[name="uploadFile"]')[0])
    let $upload = $($ul.find('li[name="upload"]')[0])
    let $rename = $($ul.find('a[name="rename"]')[0])
    let $cut = $($ul.find('a[name="cut"]')[0])
    let $paste = $($ul.find('a[name="paste"]')[0])
    let $del = $($ul.find('a[name="del"]')[0])
    let $restore = $($ul.find('a[name="restore"]')[0])
    let $permlyDel = $($ul.find('a[name="permlyDel"]')[0])
    let $close = $($ul.find('a[name="close"]')[0])

    function contextmenu(event, type) {
        // 阻止默认行为
        event.preventDefault()

        if (type === 'tr') {
            // 设置选中的tr背景颜色
            let $target = $(event.target)
            if (event.target.tagName.toLowerCase() !== 'tr') {
                $targetTr = $target.parent('tr')
            }
            $targetTr.css('background-color', '#CCCCCC')

            // 根据删除状态显示右键菜单列表
            $addFolder.addClass('hide')
            $addMdFile.addClass('hide')
            $uploadFile.addClass('hide')
            $upload.addClass('hide')
            $paste.addClass('hide')
            $close.addClass('hide')
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

        } else if (type === 'table') {
            $upload.addClass('hide')
            $rename.addClass('hide')
            $cut.addClass('hide')
            $del.addClass('hide')
            $restore.addClass('hide')
            $permlyDel.addClass('hide')
            $close.addClass('hide')
            $addFolder.removeClass('hide')
            $addMdFile.removeClass('hide')
            $uploadFile.removeClass('hide')
            $paste.removeClass('hide')

        } else {
            throw new Error('未知类型 ' + type)
        }

        // 设置菜单位置
        let x = event.clientX
        let y = event.clientY
        if (x >= document.documentElement.clientWidth - $ul.offsetWidth) {
            x = document.documentElement.clientWidth - $ul.offsetWidth
        }
        if (y >= document.documentElement.clientHeight - $ul.offsetHeight) {
            y = document.documentElement.clientHeight - $ul.offsetHeight
        }
        $ul.css('left', `${x}px`)
        $ul.css('top', `${y}px`)

        // 显示菜单
        displayMenu()
    }

    // 为 table:first-child tr 添加右键菜单事件
    let $trs = $('table:first-child tbody tr[class!="no-data"]')
    for (let i = 0, len = $trs.length; i < len; i++) {
        let $tr = $($trs[i])
        $tr.contextmenu(function (event) {
            contextmenu(event, 'tr')
        })
    }

    // 为 table:last-child 添加右键菜单事件
    let $table = $('table:last-child')
    $table.contextmenu(function (event) {
        contextmenu(event, 'table')
    })

    // 为 div 添加右键菜单事件
    $div.contextmenu(function (event) {
        // 阻止默认行为
        event.preventDefault()
    })

    // 监听 html->body 鼠标点击事件
    $body.click(function (event) {
        if (!isUploadFile) {
            // 隐藏菜单
            hideMenu()
        }
    })

    // 新增文件夹
    $addFolder.click(function () {
        let name = prompt(`${i18n.name}`, '')
        if (!name || (name = name.trim()) === '') {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }
        console.log(name)
    })

    // 新增md文件
    $addMdFile.click(function () {
        let name = prompt(`${i18n.name}`, '')
        if (!name || (name = name.trim()) === '') {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }
        console.log(name)
    })

    // 上传文件
    $uploadFile.click(function () {
        $addFolder.addClass('hide')
        $addMdFile.addClass('hide')
        $uploadFile.addClass('hide')
        $rename.addClass('hide')
        $cut.addClass('hide')
        $paste.addClass('hide')
        $del.addClass('hide')
        $restore.addClass('hide')
        $permlyDel.addClass('hide')
        $upload.removeClass('hide')
        $close.removeClass('hide')
        isUploadFile = true
        return false
    })

    // 上传
    $($ul.find('button[type="submit1"]')[0]).click(function () {
        console.log('--1-')
        return false
    })

    // 上传
    $($ul.find('from')[0]).click(function () {
        console.log('--2-')
        return false
    })

    // 重命名
    $rename.click(function (event) {
        let name = $targetTr.attr('name')
        name = prompt(`${i18n.name}`, name)
        if (!name || (name = name.trim()) === '') {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }

        let href = $rename.attr('href')
        // href = `?_method=PUT&name=${name}`
        // $rename.attr('href', href)
    })

    // 剪切
    $cut.click(function () {

    })

    // 粘贴
    $paste.click(function () {

    })

    // 删除
    $del.click(function () {
        let name = $targetTr.attr('name')
        if (!confirm(`${i18n.del} ${name} ?`)) {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }
    })

    // 恢复
    $restore.click(function () {
        let name = $targetTr.attr('name')
        if (!confirm(`${i18n.restore} ${name} ?`)) {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }
    })

    // 永久删除
    $permlyDel.click(function () {
        let name = $targetTr.attr('name')
        if (!confirm(`${i18n.permlyDel} ${name} ?`)) {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }
    })

    // 关闭
    $close.click(function () {
        // 隐藏菜单
        hideMenu()
    })

    $(document).keydown(function (event) {
        // Esc按键
        if (event.keyCode === 27 // Esc按键
            && !isUploadFile) {
            // 隐藏菜单
            hideMenu()
        }
    })

})
;