<template>
    <div class="feedback">
        <div class="feedback-item">
            <i class="iconfont iconsuggest" />建议和反馈

            <div class="feedback-qr">
                <img src="@/assets/image/qr_suggest.jpg" alt="" />
                <p>添加产品经理微信<br />问题、需求及时响应</p>
            </div>
        </div>

        <div class="feedback-item">
            <i class="iconfont icondiscuss" />产品讨论

            <div class="feedback-qr">
                <img src="@/assets/image/qr_discuss.jpg" alt="" />
                <p>加入微信用户讨论群</p>
            </div>
        </div>

        <div class="feedback-item">
            <i class="iconfont iconwechat" /><span>关注公众号</span>

            <div class="feedback-qr">
                <img src="@/assets/image/qr_wechat.jpg" alt="" />
                <p>关注ApiCat公众号<br />随时掌握ApiCat新动态</p>
            </div>
        </div>
    </div>
</template>
<style lang="scss" scoped>
    .feedback {
        padding: 6px 10px;
        position: relative;

        .feedback-item {
            cursor: pointer;
            display: flex;
            align-items: center;
            font-size: 15px;
            line-height: 1.8;
            i {
                font-size: 20px;
                margin-right: 3px;
            }

            &:hover {
                color: var(--primary-color);

                .feedback-qr {
                    display: block;
                }
            }
        }

        .feedback-qr {
            position: absolute;
            width: 180px;
            top: -4px;
            left: -200px;
            background: #fff;
            color: theme('colors.zinc.500');
            text-align: center;
            line-height: 1.5;
            border-radius: 4px;
            display: none;
            padding: 20px 0;
            box-shadow: 1px 1px 14px 0px #d4d4d4;

            img {
                width: 100px;
                height: 100px;
                display: inline-block;
            }
        }
    }
</style>
