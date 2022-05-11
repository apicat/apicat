<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Base\StringRandom;

/**
 * 字符串
 */
class StringMock extends Mock
{
    /**
     * 默认最小字符串长度
     *
     * @var int
     */
    public static $defaultMinLen = 3;

    /**
     * 默认最大字符串长度
     *
     * @var int
     */
    public static $defaultMaxLen = 10;

    /**
     * 最小字符串长度
     *
     * @var int
     */
    public static $minLen = 1;

    /**
     * 最大字符串长度
     *
     * @var int
     */
    public static $maxLen = 10000;

    /**
     * 根据规则生成数据
     *
     * @param string $rule mock规则
     * @return string
     */
    public static function generate($rule)
    {
        if (!$rule) {
            // 默认规则
            return StringRandom::randomLenLower(self::$defaultMinLen, self::$defaultMaxLen);
        }

        list($min, $max) = self::validate($rule);

        if ($min == $max) {
            return StringRandom::lowercaseString($min);
        } elseif ($min > $max) {
            return StringRandom::randomLenLower($max, $min);
        } else {
            return StringRandom::randomLenLower($min, $max);
        }
    }

    /**
     * 规则校验
     *
     * @param string $rule mock规则
     * @return array
     */
    public static function validate($rule)
    {
        if (strpos($rule, '-') !== false) {
            // 随机长度字符串
            list($min, $max) = explode('-', $rule);
        } else {
            // 固定长度字符串
            $min = $max = $rule;
        }

        $min = self::checkLimits($min, self::$minLen, self::$maxLen);
        $max = self::checkLimits($max, self::$minLen, self::$maxLen);

        return [$min, $max];
    }
}
