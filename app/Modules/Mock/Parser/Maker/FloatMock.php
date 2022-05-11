<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Base\NumberRandom;

/**
 * 小数
 */
class FloatMock extends Mock
{
    /**
     * 默认最小整数
     *
     * @var int
     */
    public static $defaultMin = 0;

    /**
     * 默认最大整数
     *
     * @var int
     */
    public static $defaultMax = 1000;

    /**
     * 最小整数
     *
     * @var int
     */
    public static $min = -1000000;

    /**
     * 最大整数
     *
     * @var int
     */
    public static $max = 1000000;

    /**
     * 默认最小小数位
     *
     * @var int
     */
    public static $defaultDmin = 1;

    /**
     * 默认最大小数位
     *
     * @var int
     */
    public static $defaultDmax = 3;

    /**
     * 最小小数位
     *
     * @var int
     */
    public static $dmin = 1;

    /**
     * 最大小数位
     *
     * @var int
     */
    public static $dmax = 10;

    /**
     * 根据规则生成数据
     *
     * @param string $rule mock规则
     * @return int|float
     */
    public static function generate($rule)
    {
        if (!$rule) {
            // 默认规则
            return NumberRandom::randomFloat(self::$defaultMin, self::$defaultMax, self::$defaultDmin, self::$defaultDmax);
        }

        list($min, $max, $dmin, $dmax) = self::validate($rule);
        return NumberRandom::randomFloat($min, $max, $dmin, $dmax);
    }

    /**
     * 规则校验
     *
     * @param string $rule mock规则
     * @return array
     */
    public static function validate($rule)
    {
        list($integer, $decimal) = explode('.', $rule);

        if (strpos($integer, '~') !== false) {
            // 随机长度字符串
            list($min, $max) = explode('~', $integer);
        } else {
            // 固定长度字符串
            $min = $max = $integer;
        }

        if ($decimal) {
            if (strpos($decimal, '-') !== false) {
                // 随机长度字符串
                list($dmin, $dmax) = explode('-', $decimal);
            } else {
                // 固定长度字符串
                $dmin = $dmax = $decimal;
            }
        } else {
            $dmin = $dmax = 0;
        }

        $min = self::checkLimits($min, self::$min, self::$max);
        $max = self::checkLimits($max, self::$min, self::$max);
        $dmin = self::checkLimits($dmin, self::$dmin, self::$dmax);
        $dmax = self::checkLimits($dmax, self::$dmin, self::$dmax);

        return [$min, $max, $dmin, $dmax];
    }
}
