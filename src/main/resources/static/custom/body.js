/**
 * @author xiangqian
 * @date 22:21 2024/03/02
 */
;

// 鼠标进入和离开<tr>时，显示和隐藏操作
$(function () {
    let $trs = $('body > center > table tbody tr');

    function callback(event) {
        // 获取目标<tr>元素
        let $tr = $(event.currentTarget);

        // 获取目标<tr>元素中的最后一个<td>元素
        let $lastTd = $($tr.find('td:last'));

        // 事件类型
        const type = event.type;

        // 鼠标进入<tr>时执行的函数，显示<td>元素
        if (type === 'mouseenter') {
            $lastTd.css('visibility', 'visible');
            return;
        }

        // 鼠标离开<tr>时执行的函数，隐藏<td>元素
        if (type === 'mouseleave') {
            $lastTd.css('visibility', 'hidden');
            return;
        }
    }

    // 鼠标进入<tr>时执行的函数
    $trs.mouseenter(callback);

    // 鼠标离开<tr>时执行的函数
    $trs.mouseleave(callback);
});

// prompt
$(function () {
    let $inputs = $('body > center form input[prompt-message]');
    for (let i = 0, length = $inputs.length; i < length; i++) {
        let $input = $($inputs[i]);
        // console.log($input);

        // 查找元素最近的父级<form>元素
        let $form = $input.closest('form');

        // 添加<form>单击事件
        $form.click(function (event) {
            let value = prompt($input.attr('prompt-message'), $input.attr('prompt-default'));

            // 用户取消操作或者输入为空
            if (value == null || (value = value.trim()) == '') {
                return false;
            }

            // 设置值
            $input.val(value);

            return true;
        });
    }
});

// confirm
$(function () {
    let $buttons = $('body > center form button[type="submit"][confirm-message]');
    for (let i = 0, length = $buttons.length; i < length; i++) {
        let $button = $($buttons[i]);
        // console.log($button);

        // 查找元素最近的父级<form>元素
        let $form = $button.closest('form');

        // 添加<form>单击事件
        $form.click(function (event) {
            return confirm($button.attr('confirm-message'));
        });
    }
});
