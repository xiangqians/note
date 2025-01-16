/**
 * @author xiangqian
 * @date 21:55 2024/12/06
 */
;

// 自动保存
function AutoSave(url, $title, $textarea) {
    this.url = url;

    this.$title = $title;
    this.title = $title.text();

    this.$textarea = $textarea;
    this.hash = md5($textarea.val());

    // 最近保存时间
    this.lastSavedTime = null;
}

AutoSave.prototype.getContent = function () {
    return this.$textarea.val();
}

/**
 * 判断内容是否已发生改变
 */
AutoSave.prototype.hasChanged = function () {
    return this.hash !== md5(this.getContent());
}

AutoSave.prototype.save0 = function () {
    let self = this;
    let content = self.getContent();
    http("PUT", self.url + "?t=" + new Date().getTime(), "text/plain", content, function (data, status, xhr) {
        self.hash = md5(content);

        // 当前时间
        let date = new Date();
        // 当前小时
        let hours = date.getHours();
        // 当前分钟
        let minutes = date.getMinutes();
        // 当前秒
        let seconds = date.getSeconds();
        // 最近保存时间
        // 使用 padStart 方法确保每个数字都是两位
        self.lastSavedTime = [hours.toString().padStart(2, "0"), minutes.toString().padStart(2, "0"), seconds.toString().padStart(2, "0")].join(":");

        self.$title.text(self.title + "（已自动保存）");
    });
}

AutoSave.prototype.save = function () {
    const self = this;

    // 判断内容是否已发生改变
    if (self.hasChanged()) {
        self.save0();
    } else if (self.lastSavedTime != null) {
        self.$title.text(self.title + "（最近保存 " + self.lastSavedTime + "）");
    }

    // 等待 2 秒后再次执行
    setTimeout(() => {
        // 使用箭头函数，this 指向当前对象
        self.save();
    }, 2000);
}

AutoSave.prototype.start = function () {
    let self = this;
    // 在浏览器窗口或标签页即将关闭、刷新、导航到其他页面时触发，用来提醒用户是否确定离开当前页面，防止用户丢失未保存的数据。
    window.onbeforeunload = function (event) {
        let e = window.event || event
        // 判断内容是否已发生改变
        if (self.hasChanged()) {
            e.returnValue = "数据发生改变"
            return e.returnValue
        }
    }

    // 自动保存
    this.save();
}
;
