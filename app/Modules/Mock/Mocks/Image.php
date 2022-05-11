<?php

namespace App\Modules\Mock\Mocks;

/**
 * 模拟图片生成
 */
class Image
{
    /**
     * 随机生成图片流
     *
     * @param int $width 图片宽度
     * @param int $height 图片高度
     * @param string $type 图片类型
     * @return string
     */
    public static function stream($width = 200, $height = 150, $type = 'jpg')
    {
        return Base\ImageRandom::stream($width, $height, $type);
    }

    /**
     * 随机生成图片base64内容
     *
     * @param int $width 图片宽度
     * @param int $height 图片高度
     * @param string $type 图片类型
     * @return string
     */
    public static function base64($width = 200, $height = 150, $type = 'jpg')
    {
        if (!$content = Base\ImageRandom::stream($width, $height, $type)) {
            return '';
        }

        return Base\ImageRandom::$base64Prefix[$type] . base64_encode($content);
    }

    /**
     * 随机生成图片url
     *
     * @param int $width 图片宽度
     * @param int $height 图片高度
     * @param string $type 图片类型
     * @return string
     */
    public static function url($width = 200, $height = 150, $type = 'jpg')
    {
        return rtrim(env('APP_MOCK_URL'), '/') . '/mock_image/' . $width . 'x' . $height . '.' . $type;
    }
}
