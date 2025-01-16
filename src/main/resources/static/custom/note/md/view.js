/**
 * @author xiangqian
 * @date 21:08 2024/11/22
 */
;
$(function () {
    // 预览页
    let $toc = $("#toc");

    // 预览页
    let $preview = $("#preview");

    // 编辑页
    let $editor = $("#editor");

    // 是否显示编辑页
    let isEditor = false;

    // 移除以 "col-xs-" 开头的class
    function removeColXsClass(index, className) {
        return (className.match(/(^|\s)col-xs-\S+/g) || []).join(" ")
    }

    function dblclick(event) {
        // 获取选中文本
        let selection = window.getSelection().toString().trim();
        // console.log(selection);
        // 如果选中了文本，则不触发事件
        if (selection !== "") {
            return;
        }

        isEditor = !isEditor;

        $editor.removeClass(removeColXsClass);
        $preview.removeClass(removeColXsClass);

        if (isEditor) {
            // console.log("编辑");
            $editor.attr("edit", "");
            $editor.css("display", "");
            $editor.addClass("col-xs-5");
            $preview.addClass("col-xs-5");
        } else {
            // console.log("预览");
            $editor.removeAttr("edit");
            $editor.css("display", "none");
            $preview.addClass("col-xs-10");
        }

        // 阻止事件的默认行为
        // event.preventDefault();
    }

    // 预览页绑定双击事件
    $toc.dblclick(dblclick);
    $preview.dblclick(dblclick);

    // 前端js获取剪切板中的文本或者图片文件数据
    // 监听paste事件，在进行粘贴的时候会触发
    // https://developer.mozilla.org/en-US/docs/Web/API/Element/paste_event
    document.addEventListener("paste", function (event) {
        // 获取粘贴板数据项
        const items = (event.clipboardData || event.originalEvent.clipboardData).items;
        // console.log("items", items);

        // dropEffect
        // String，默认是 none
        // let dropEffect = clipboardData.dropEffect;
        // console.log("dropEffect", dropEffect);

        // effectAllowed
        // String，默认是 uninitialized
        // let effectAllowed = clipboardData.effectAllowed;
        // console.log("effectAllowed", effectAllowed);

        // files
        // FileList，文件列表
        // let files = clipboardData.files;
        // console.log("files", files);

        // items
        // DataTransferItemList，剪切板中的各项数据，Chrome有该属性，Safari没有。
        //
        // DataTransferItem属性：
        // 1、kind：一般为string或者file
        // 2、type：具体的数据类型，例如具体是哪种类型字符串或者哪种类型的文件，即MIME-Type
        // DataTransferItem方法：
        // 1、getAsFile：如果kind是file，可以用该方法获取到文件
        // 2、getAsString：入参是回调函数；如果kind是string，可以用该方法获取到字符串，字符串需要用回调函数得到，回调函数的第一个参数就是剪切板中的字符串
        //
        // let items = clipboardData.items;
        // console.log("items", items);

        // types
        // Array，剪切板中的数据类型 该属性在Safari下比较混乱
        // 一般types中常见的值有：text/plain（普通字符串）、text/html（带有样式的html）、Files（文件）。
        // let types = clipboardData.types;
        // console.log("types", types);

        const mimeTypes = {
            "image/png": "png",
            "image/jpeg": "jpg",  // 对于 jpeg 文件，mime 类型仍是 "image/jpeg"
            "image/gif": "gif",
            "image/webp": "webp",
            "image/x-icon": "ico"
        };

        for (let i = 0, length = items.length; i < length; i++) {
            let item = items[i];
            // console.log("item", item);
            const kind = item.kind;
            if (kind === "file") {
                const type = item.type;
                // console.log(type)
                if (/^image\/[jpeg|png|gif|jpg]/.test(type)) {
                    // blob 就是从剪切板获得的文件 可以进行上传或其他操作
                    let blob = item.getAsFile();
                    if (blob.size === 0) {
                        continue;
                    }

                    // 上传图片
                    let formData = new FormData();
                    formData.append("file", blob, "image." + mimeTypes[type]);
                    http("POST", window.contextPath + "/image/upload?t=" + new Date().getTime(), false, formData, function (data, status, xhr) {
                        // console.log(data);

                        let message = data["message"];
                        if (message != undefined) {
                            alert(message);
                            return;
                        }

                        let id = data.id;
                        let name = data.name;
                        // let type = data.type;
                        document.execCommand("insertText", false, "![" + name + "](/image/" + id + ")");
                        return;
                    });

                    // 阻止事件的默认行为
                    event.preventDefault();
                    return;
                } else {
                    alert("不支持粘贴此数据项类型 " + type);
                }
            }
        }
    });

});

