<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Chinese;

/**
 * 中文词语
 */
class CwordMock extends Mock
{
    /**
     * 默认最小汉字个数
     *
     * @var int
     */
    public static $defaultMin = 1;

    /**
     * 默认最大汉字个数
     *
     * @var int
     */
    public static $defaultMax = 4;

    /**
     * 最小汉字个数
     *
     * @var int
     */
    public static $min = 1;

    /**
     * 最大汉字个数
     *
     * @var int
     */
    public static $max = 10;

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
            return Chinese::word(self::$defaultMin, self::$defaultMax);
        }

        list($min, $max) = self::validate($rule);
        return Chinese::word($min, $max);
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
            // 随机句子个数
            list($min, $max) = explode('-', $rule);
        } else {
            // 固定句子个数
            $min = $max = $rule;
        }

        $min = self::checkLimits($min, self::$min, self::$max);
        $max = self::checkLimits($max, self::$min, self::$max);

        return [$min, $max];
    }
}
