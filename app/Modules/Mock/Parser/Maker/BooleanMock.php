<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Base\BooleanRandom;

/**
 * 布尔值
 */
class BooleanMock extends Mock
{
    /**
     * 默认概率
     *
     * @var int
     */
    public static $default = 50;

    /**
     * 最小概率
     *
     * @var int
     */
    public static $min = 0;

    /**
     * 最大概率
     *
     * @var int
     */
    public static $max = 100;

    /**
     * 根据规则生成数据
     *
     * @param string $rule mock规则
     * @return boolean
     */
    public static function generate($rule)
    {
        if ($rule === '') {
            // 默认规则
            return BooleanRandom::randomTrue(self::$default);
        }

        if (strtolower($rule) == 'true') {
            return true;
        }
        if (strtolower($rule) == 'false') {
            return false;
        }

        $probability = self::validate($rule);
        return BooleanRandom::randomTrue($probability);
    }

    /**
     * 规则校验
     *
     * @param string $rule mock规则
     * @return int
     */
    public static function validate($rule)
    {
        return self::checkLimits($rule, self::$min, self::$max);
    }
}