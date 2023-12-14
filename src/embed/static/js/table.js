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
    let $addFolder = $(`<li><form id="addFolder" action="${variable.contextPath}/${variable.table}/addfolder?t=${new Date().getTime()}" method="post"><input name="pid" hidden><input name="name" hidden><a href="javascript:addFolder.submit();">${variable.i18n.addFolder}</a></form></li>`)
    $ul.append($addFolder)
    // 【新增MD文件】
    let $addMdFile = $(`<li><form id="addMdFile" action="${variable.contextPath}/${variable.table}/addmdfile?t=${new Date().getTime()}" method="post"><input name="pid" hidden><input name="name" hidden><a href="javascript:addMdFile.submit();">${variable.i18n.addMdFile}</a></form></li>`)
    $ul.append($addMdFile)
    // 【上传文件】
    let $uploadFile = $(`<li><a href="javascript:void(0);">${variable.i18n.uploadFile}</a></li>`)
    $ul.append($uploadFile)
    // 【上传】
    let $upload = $(`<li><form action="${variable.contextPath}/${variable.table}/upload?t=${new Date().getTime()}" method="post" enctype="multipart/form-data"><input name="file" type="file"/><input name="pid" hidden><button type="submit">${variable.i18n.upload}</button></form></li>`)
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
    let $paste = $(`<li><form id="paste" action="${variable.contextPath}/${variable.table}/paste?t=${new Date().getTime()}" method="post"><input name="fromId" hidden><input name="toId" hidden><a href="javascript:paste.submit();">${variable.i18n.paste}</a></form></li>`)
    $ul.append($paste)
    // 【删除】
    let $del = $(`<li><form id="del" action="${variable.contextPath}/${variable.table}/{id}/del?t=${new Date().getTime()}" method="post"><a href="javascript:del.submit();">${variable.i18n.del}</a></form></li>`)
    $ul.append($del)
    // 【恢复】
    let $restore = $(`<li><form id="restore" action="${variable.contextPath}/${variable.table}/{id}/restore?t=${new Date().getTime()}" method="post"><a href="javascript:restore.submit();">${variable.i18n.restore}</a></form></li>`)
    $ul.append($restore)
    // 【永久删除】
    let $permlyDel = $(`<li><form id="permlyDel" action="${variable.contextPath}/${variable.table}/{id}/permlydel?t=${new Date().getTime()}" method="post"><a href="javascript:permlyDel.submit();">${variable.i18n.permlyDel}</a></form></li>`)
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
                if (variable.pNote && variable.pNote.id >= 0) {
                    displayElements($addFolder, $addMdFile, $uploadFile)
                    if (getCutData()) {
                        displayElements($paste)
                    } else {
                        hideElements($paste)
                    }

                } else {
                    hideElements($addFolder, $addMdFile, $uploadFile, $paste)
                    return
                }
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
    $($addFolder.find('a')[0]).click(function () {
        let pid = variable.pNote.id
        let name = prompt(`${variable.i18n.name}`, '')
        if (!name || (name = name.trim()) === '') {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }

        $($addFolder.find('input[name="pid"]')[0]).val(pid)
        $($addFolder.find('input[name="name"]')[0]).val(name)
        return true
    })

    // 【新增md文件】
    $($addMdFile.find('a')[0]).click(function () {
        let pid = variable.pNote.id
        let name = prompt(`${variable.i18n.name}`, '')
        if (!name || (name = name.trim()) === '') {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }

        $($addMdFile.find('input[name="pid"]')[0]).val(pid)
        $($addMdFile.find('input[name="name"]')[0]).val(name)
        return true
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
        let pid = 0
        if (variable.pNote) {
            pid = variable.pNote.id
        }
        $($upload.find('input[name="pid"]')[0]).val(pid)
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
        if (variable.table === 'note') {
            return
        }

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


    // 根据id获取$tr
    function get$TrById(id) {
        let tr = $(`tr[id="${id}"]`)
        if (tr.length > 0) {
            return $(tr[0])
        }
        return null
    }

    // 移除tr灰色样式
    function addTrGrayClass(id) {
        let $tr = get$TrById(id)
        if ($tr) {
            $tr.addClass('gray')
        }
    }

    // 移除tr灰色样式
    function removeTrGrayClass(id) {
        let $tr = get$TrById(id)
        if ($tr) {
            $tr.removeClass('gray')
        }
    }

    const CUT_STORAGE_KEY = 'cut_data'

    // 获取剪切数据
    function getCutData() {
        return storage.getObject(CUT_STORAGE_KEY)
    }

    // 设置剪切数据
    function setCutData(data) {
        storage.setObject(CUT_STORAGE_KEY, data)
    }

    // 设置已剪切的tr为灰色
    (function () {
        let data = getCutData()
        if (data) {
            addTrGrayClass(data.id)
        }
    })()

    // 【剪切】
    $($cut.find('a')[0]).click(function () {
        let pNote = variable.pNote
        let id = $selectedTr.attr('id')
        let name = `${pNote.namesStr}/${$selectedTr.attr('name')}`

        let data = getCutData()
        if (data) {
            removeTrGrayClass(data.id)
        }

        data = {id: id, name: name}
        setCutData(data)
        addTrGrayClass(data.id)
    })

    // 【粘贴】
    $($paste.find('a')[0]).click(function () {
        let data = getCutData()
        if (data) {
            let pNote = variable.pNote
            if (!confirm(`${data.name} -> ${pNote.id === 0 ? "/" : pNote.namesStr} ?`)) {
                // 隐藏菜单
                hideMenu()
                // 取消 <a></a> 默认行为
                return false
            }

            $($paste.find('input[name="fromId"]')[0]).val(data.id)
            $($paste.find('input[name="toId"]')[0]).val(pNote.id)
            setCutData(null)
            return true
        }
        return false
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
    $($restore.find('a')[0]).click(function () {
        let id = $selectedTr.attr('id')
        let name = $selectedTr.attr('name')
        if (!confirm(`${variable.i18n.restore} ${name} ?`)) {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }

        let $form = $($restore.find('form')[0])
        $form.attr('action', `${variable.contextPath}/${variable.table}/${id}/restore?t=${new Date().getTime()}`)
    })

    // 【永久删除】
    $($permlyDel.find('a')[0]).click(function () {
        let id = $selectedTr.attr('id')
        let name = $selectedTr.attr('name')
        if (!confirm(`${variable.i18n.permlyDel} ${name} ?`)) {
            // 隐藏菜单
            hideMenu()
            // 取消 <a></a> 默认行为
            return false
        }

        let $form = $($permlyDel.find('form')[0])
        $form.attr('action', `${variable.contextPath}/${variable.table}/${id}/permlydel?t=${new Date().getTime()}`)
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