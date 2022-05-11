<?php

namespace App\Modules\Mock\Mocks\Base;

/**
 * 随机英文字符串
 */
class StringRandom
{
    /**
     * 小写英文字母
     *
     * @var string
     */
    public static $lowercase = 'abcdefghijklmnopqrstuvwxyz';

    /**
     * 大写英文字母
     *
     * @var string
     */
    public static $uppercase = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

    /**
     * 小写和大写英文字母
     *
     * @var string
     */
    public static $lnucase = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';

    /**
     * 小写随机字符串
     *
     * @param int $len 字符串长度
     * @return string
     */
    public static function lowercaseString(int $len)
    {
        return self::randomString($len, self::$lowercase);
    }

    /**
     * 大写随机字符串
     *
     * @param int $len 字符串长度
     * @return string
     */
    public static function uppercaseString(int $len)
    {
        return self::randomString($len, self::$uppercase);
    }

    /**
     * 小写和大写随机字符串
     *
     * @param int $len 字符串长度
     * @return string
     */
    public static function lnucaseString(int $len)
    {
        return self::randomString($len, self::$lnucase);
    }

    /**
     * 首字母大写随机字符串
     *
     * @param int $len 字符串长度
     * @return string
     */
    public static function firstUppercaseString(int $len)
    {
        $index = mt_rand(0, 25);
        return self::$uppercase[$index] . self::lowercaseString($len - 1);
    }

    /**
     * 随机长度小写随机字符串
     *
     * @param int $min 最小长度
     * @param int $max 最大长度
     * @return string
     */
    public static function randomLenLower(int $min = 3, int $max = 10)
    {
        $len = mt_rand($min, $max);
        return self::lowercaseString($len);
    }

    /**
     * 随机长度大写随机字符串
     *
     * @param int $min 最小长度
     * @param int $max 最大长度
     * @return string
     */
    public static function randomLenUpper(int $min = 3, int $max = 10)
    {
        $len = mt_rand($min, $max);
        return self::uppercaseString($len);
    }

    /**
     * 随机长度小写和大写随机字符串
     *
     * @param int $min 最小长度
     * @param int $max 最大长度
     * @return string
     */
    public static function randomLenLnu(int $min = 3, int $max = 10)
    {
        $len = mt_rand($min, $max);
        return self::lnucaseString($len);
    }

    /**
     * 随机长度首字母大写随机字符串
     *
     * @param int $min 最小长度
     * @param int $max 最大长度
     * @return string
     */
    public static function randomLenFirstUpper(int $min = 3, int $max = 10)
    {
        $len = mt_rand($min, $max);
        return self::firstUppercaseString($len);
    }

    /**
     * 生成随机字符串
     *
     * @param int $len 字符串长度
     * @param string $source 随机源
     * @return string
     */
    protected static function randomString($len, $source)
    {
        $strLength = mb_strlen($source);
        $maxIndex = $strLength - 1;

        $start = 0;
        $rand = mt_rand(0, $maxIndex);

        $result = [];
        for ($i = 0; $i < $len; $i++) {
            $start += $rand;
            if ($start > $maxIndex) {
                $start -= $maxIndex;
            }

            $result[] = $source[$start];
        }

        return implode('', $result);
    }
}
