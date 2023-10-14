;
(function () {
    return;
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

// body添加右键菜单
function bodyMenu($ul) {
    let $body = $($('body')[0])
    $body.contextmenu(function (event) {
        alert("处理程序.contextmenu()被调用。");
    });
}

function tableBodyMenu() {
}
