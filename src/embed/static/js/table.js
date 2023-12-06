/**
 * @author xiangqian
 * @date 22:35 2023/10/15
 */
;
$(function () {
    let $body = $($('body')[0])

    // 右键菜单列表
    let $ul = $('<ul></ul>')
    // 【新增文件夹】
    let $addFolder = $(`<li><a href="javascript:void(0);">${variable.i18n.addFolder}</a></li>`)
    $ul.append($addFolder)
    // 【新增MD文件】
    let $addMdFile = $(`<li><a href="javascript:void(0);">${variable.i18n.addMdFile}</a></li>`)
    $ul.append($addMdFile)
    // 【上传文件】
    let $uploadFile = $(`<li><a href="javascript:void(0);">${variable.i18n.uploadFile}</a></li>`)
    $ul.append($uploadFile)
    // 【上传】
    let $upload = $(`<li><form action="${variable.contextPath}/${variable.table}/upload?t=${new Date().getTime()}" method="post" enctype="multipart/form-data"><input name="file" type="file"/><button type="submit">${variable.i18n.upload}</button></form></li>`)
    setAccept($($upload.find('input[type="file"]')[0]))
    $ul.append($upload)
    // 【重命名】
    let $rename = $(`<li><form id="rename" action="${variable.contextPath}/${variable.table}/rename?t=${new Date().getTime()}" method="post"><input name="id" hidden><input name="name" hidden><a href="javascript:rename.submit();">${variable.i18n.rename}</a></form></li>`)
    $ul.append($rename)
    // 【复制地址】
    let $copyAddress = $(`<li><a href="javascript:void(0);">${variable.i18n.copyAddress}</a></li>`)
    $ul.append($copyAddress)
    // 【剪切】
    let $cut = $(`<li><a href="javascript:void(0);">${variable.i18n.cut}</a></li>`)
    $ul.append($cut)
    // 【粘贴】
    let $paste = $(`<li><a href="javascript:void(0);">${variable.i18n.paste}</a></li>`)
    $ul.append($paste)
    // 【删除】
    let $del = $(`<li><form id="del" action="${variable.contextPath}/${variable.table}/{id}/del?t=${new Date().getTime()}" method="post"><a href="javascript:del.submit();">${variable.i18n.del}</a></form></li>`)
    $ul.append($del)
    // 【恢复】
    let $restore = $(`<li><a href="javascript:void(0);">${variable.i18n.restore}</a></li>`)
    $ul.append($restore)
    // 【永久删除】
    let $permlyDel = $(`<li><a href="javascript:void(0);">${variable.i18n.permlyDel}</a></li>`)
    $ul.append($permlyDel)
    // 【关闭】
    let $close = $(`<li><a href="javascript:void(0);">${variable.i18n.close}</a></li>`)
    $ul.append($close)
    // 遮罩层
    let $div = $('<div class="menu"></div>')
    $div.append($ul)
    $body.prepend($div)

    // 右键选中的<tr>
    let $selectedTr = null

    // 是否是上传文件操作
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
        if ($selectedTr != null) {
            $selectedTr.css('background-color', '')
            $selectedTr = null
        }

        // 允许页面滚动
        $body.css('overflow', '')

        // 隐藏菜单
        $ul.css('display', 'none')
        $div.css('display', 'none')
    }

    // 右键菜单监听
    function contextmenu(event, type) {
        // 阻止默认行为
        event.preventDefault()

        if (type === 'tr') {
            // 设置选中的tr背景颜色
            let $target = $(event.target)
            if (event.target.tagName.toLowerCase() !== 'tr') {
                $selectedTr = $target.closest('tr')
            }
            $selectedTr.css('background-color', '#CCCCCC')

            // 根据删除状态显示右键菜单列表
            hideElements($addFolder, $addMdFile, $uploadFile, $upload, $paste, $close)
            let del = $selectedTr.attr('del')
            if (del === "0") {
                hideElements($restore, $permlyDel)
                displayElements($rename, $del)

                // note
                if (variable.table === 'note') {
                    hideElements($copyAddress)
                    displayElements($cut)
                }
                // 其他：image、audio、video
                else {
                    hideElements($cut)
                    displayElements($copyAddress)
                }

            } else if (del === "1") {
                hideElements($rename, $copyAddress, $cut, $del)
                displayElements($restore, $permlyDel)
            }

        } else if (type === 'table') {
            hideElements($upload, $rename, $copyAddress, $cut, $del, $restore, $permlyDel, $close)

            // note
            if (variable.table === 'note') {
                displayElements($addFolder, $addMdFile, $uploadFile, $paste)
            }
            // 其他：image、audio、video
            else {
                hideElements($addFolder, $addMdFile, $paste)
                displayElements($uploadFile)
            }

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
    let trs = $('table:first-child tbody tr[class!="no-data"]')
    for (let i = 0, len = trs.length; i < len; i++) {
        let $tr = $(trs[i])
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
    $body.click(function () {
        // 如果不是上传文件操作，单击空白处则隐藏菜单
        if (!isUploadFile) {
            hideMenu()
        }
    })

    // 【新增文件夹】
    $addFolder.click(function () {
        let name = prompt(`${variable.i18n.name}`, '')
        if (!name || (name = name.trim()) === '') {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }
        console.log(name)
    })

    // 【新增md文件】
    $addMdFile.click(function () {
        let name = prompt(`${variable.i18n.name}`, '')
        if (!name || (name = name.trim()) === '') {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }
        console.log(name)
    })

    // 【上传文件】
    $($uploadFile.find('a')[0]).click(function () {
        hideElements($addFolder, $addMdFile, $uploadFile, $rename, $copyAddress, $cut, $paste, $del, $restore, $permlyDel)
        displayElements($upload, $close)
        isUploadFile = true
        return false
    })

    // 【上传】
    $($upload.find('[type="submit"]')[0]).click(function () {
        return true
    })

    // 【重命名】
    $($rename.find('a')[0]).click(function () {
        let id = $selectedTr.attr('id')
        let name = $selectedTr.attr('name')
        name = prompt(`${variable.i18n.name}`, name)
        if (!name || (name = name.trim()) === '') {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }

        $($rename.find('input[name="id"]')[0]).val(id)
        $($rename.find('input[name="name"]')[0]).val(name)
        return true
    })

    // 【复制地址】
    let $copyAddressTr = null
    $($copyAddress.find('a')[0]).click(function () {
        $copyAddressTr = $selectedTr
    })
    ;(function () {
        // 销毁clipboard，如果存在的话
        if (window.clipboard) {
            window.clipboard.destroy()
            window.clipboard = null
        }

        let clipboard = new ClipboardJS($copyAddress.find('a')[0], {
            text: function () {
                let id = $copyAddressTr.attr('id')
                let name = $copyAddressTr.attr('name')
                $copyAddressTr = null
                return `![${name}](/${variable.table}/${id})`
            }
        })

        clipboard.on('success', function (e) {
            // console.info('Action:', e.action)
            // console.info('Text:', e.text)
            // console.info('Trigger:', e.trigger)
            alert(variable.i18n.copied)
        })

        clipboard.on('error', function (e) {
            // console.info('Action:', e.action)
            // console.info('Text:', e.text)
            // console.info('Trigger:', e.trigger)
            alert(e)
        })

        window.clipboard = clipboard
    })();

    // 【剪切】
    $cut.click(function () {
    })

    // 【粘贴】
    $paste.click(function () {
    })

    // 【删除】
    $($del.find('a')[0]).click(function () {
        let id = $selectedTr.attr('id')
        let name = $selectedTr.attr('name')
        if (!confirm(`${variable.i18n.del} ${name} ?`)) {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }

        let $form = $($del.find('form')[0])
        $form.attr('action', `${variable.contextPath}/${variable.table}/${id}/del?t=${new Date().getTime()}`)
    })

    // 【恢复】
    $restore.click(function () {
        let name = $selectedTr.attr('name')
        if (!confirm(`${variable.i18n.restore} ${name} ?`)) {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }
    })

    // 【永久删除】
    $permlyDel.click(function () {
        let name = $selectedTr.attr('name')
        if (!confirm(`${variable.i18n.permlyDel} ${name} ?`)) {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }
    })

    // 【关闭】
    $($close.find('a')[0]).click(function () {
        // 隐藏菜单
        hideMenu()
    })

    $(document).keydown(function (event) {
        // 按下 Esc键，并且不是上传文件操作时，隐藏菜单
        if (event.keyCode === 27  // Esc键
            && !isUploadFile) { // 不是上传文件操作时
            // 隐藏菜单
            hideMenu()
        }
    })

})
;