<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Base\NumberRandom;

/**
 * 整数
 */
class IntMock extends Mock
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
     * 根据规则生成数据
     *
     * @param string $rule mock规则
     * @return int|float
     */
    public static function generate($rule)
    {
        if (!$rule) {
            // 默认规则
            return NumberRandom::randomInt(self::$defaultMin, self::$defaultMax);
        }

        list($min, $max) = self::validate($rule);
        return NumberRandom::randomInt($min, $max);
    }

    /**
     * 规则校验
     *
     * @param string $rule mock规则
     * @return array
     */
    public static function validate($rule)
    {
        if (strpos($rule, '~') !== false) {
            // 随机长度字符串
            list($min, $max) = explode('~', $rule);
        } else {
            // 固定长度字符串
            $min = $max = $rule;
        }

        $min = self::checkLimits($min, self::$min, self::$max);
        $max = self::checkLimits($max, self::$min, self::$max);

        return [$min, $max];
    }
}
