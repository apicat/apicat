<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\English;

/**
 * 英文段落
 */
class ParagraphMock extends Mock
{
    /**
     * 默认最小句子个数
     *
     * @var int
     */
    public static $defaultMin = 3;

    /**
     * 默认最大句子个数
     *
     * @var int
     */
    public static $defaultMax = 7;

    /**
     * 最小句子个数
     *
     * @var int
     */
    public static $min = 1;

    /**
     * 最大句子个数
     *
     * @var int
     */
    public static $max = 20;

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
            return English::paragraph(self::$defaultMin, self::$defaultMax);
        }

        list($min, $max) = self::validate($rule);
        return English::paragraph($min, $max);
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
