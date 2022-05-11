<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Image;

/**
 * 图片
 */
class ImageMock extends Mock
{
    /**
     * 默认宽度
     *
     * @var int
     */
    public static $defaultWidth = 200;

    /**
     * 默认高度
     *
     * @var int
     */
    public static $defaultHeight = 150;

    /**
     * 默认类型
     *
     * @var int
     */
    public static $defaultType = 'jpeg';

    /**
     * 最小宽度
     *
     * @var int
     */
    public static $minWidth = 20;

    /**
     * 最大宽度
     *
     * @var int
     */
    public static $maxWidth = 1024;

    /**
     * 最小高度
     *
     * @var int
     */
    public static $minHeight = 20;

    /**
     * 最大高度
     *
     * @var int
     */
    public static $maxHeight = 1024;

    /**
     * 允许的类型
     *
     * @var array
     */
    public static $allowTypes = ['jpeg', 'jpg', 'png'];

    /**
     * 根据规则生成图片流
     *
     * @param string $rule mock规则
     * @return boolean
     */
    public static function stream($rule)
    {
        if (!$rule) {
            // 默认规则
            return Image::stream(self::$defaultWidth, self::$defaultHeight, self::$defaultType);
        }

        list($width, $height, $type) = self::validate($rule);
        return Image::stream($width, $height, $type);
    }

    /**
     * 根据规则生成base64图片内容
     *
     * @param string $rule mock规则
     * @return boolean
     */
    public static function base64($rule)
    {
        if (!$rule) {
            // 默认规则
            return Image::base64(self::$defaultWidth, self::$defaultHeight, self::$defaultType);
        }

        list($width, $height, $type) = self::validate($rule);
        return Image::base64($width, $height, $type);
    }

    /**
     * 根据规则生成图片url
     *
     * @param string $rule mock规则
     * @return boolean
     */
    public static function url($rule)
    {
        if (!$rule) {
            // 默认规则
            return Image::url(self::$defaultWidth, self::$defaultHeight, self::$defaultType);
        }

        list($width, $height, $type) = self::validate($rule);
        return Image::url($width, $height, $type);
    }

    /**
     * 规则校验
     *
     * @param string $rule mock规则
     * @return array
     */
    public static function validate($rule)
    {
        if (strpos($rule, ',') !== false) {
            // 有指定类型
            list($size, $type) = explode(',', $rule);
        } else {
            // 无指定类型
            $size = $rule;
            $type = self::$defaultType;
        }

        if (strpos($size, '*') !== false) {
            // 随机长度字符串
            list($width, $height) = explode('*', $size);
        } else {
            // 固定长度字符串
            $width = $height = $size;
        }

        $width = self::checkLimits($width, self::$minWidth, self::$maxWidth);
        $height = self::checkLimits($height, self::$minHeight, self::$maxHeight);
        $type = self::checkIn($type, self::$allowTypes, self::$defaultType);

        return [$width, $height, $type];
    }
}