<?php

namespace App\Modules\Mock\Mocks\Base;

use Intervention\Image\Facades\Image;

/**
 * 随机图片
 */
class ImageRandom
{
    /**
     * 支持的图片类型
     *
     * @var array
     */
    public static $types = ['jpg', 'jpeg', 'png'];

    /**
     * 图片类型对应的base64前缀
     *
     * @var array
     */
    public static $base64Prefix = [
        'jpg' => 'data:image/jpg;base64,',
        'jpeg' => 'data:image/jpeg;base64,',
        'png' => 'data:image/png;base64,'
    ];

    /**
     * 默认图片颜色
     *
     * @var array
     */
    public static $colors = [
        '#72BFF9', '#48A0F8', '#3274B5', '#1F4D7B',
        '#99FAEB', '#6CE4CF', '#4CA88F', '#2E6A66', '#D5D5D5',
        '#A4F769', '#82D552', '#54AF32', '#306F1D', '#929292',
        '#FCF170', '#F9D957', '#F2B13E', '#E1792E', '#5E5E5E',
        '#F09B91', '#EC6F57', '#DA3B26', '#A62A17',
        '#F19AC9', '#EB539F', '#C23175', '#8A2052'
    ];

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
        if ($type == 'jpeg') {
            $type = 'jpg';
        }

        if (!in_array($type, self::$types)) {
            return '';
        }

        $color = array_rand(self::$colors);

        return Image::canvas($width, $height, self::$colors[$color])->stream($type, 60);
    }
}
