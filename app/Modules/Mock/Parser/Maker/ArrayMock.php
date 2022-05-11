<?php

namespace App\Modules\Mock\Parser\Maker;

class ArrayMock extends Mock
{
    /**
     * 默认最小循环次数
     *
     * @var int
     */
    public static $defaultMin = 1;

    /**
     * 默认最大循环次数
     *
     * @var int
     */
    public static $defaultMax = 5;

    /**
     * 最小循环次数
     *
     * @var int
     */
    public static $min = 0;

    /**
     * 最大循环次数
     *
     * @var int
     */
    public static $max = 50;

    /**
     * 根据规则生成数据
     *
     * @param string $rule mock规则
     * @param array $subParams 子参数
     * @return array
     */
    public static function generate($rule, $subParams)
    {
        if ($rule === '') {
            $rule = self::$defaultMin . '-' . self::$defaultMax;
        }

        $count = self::validate($rule);
        if ($count == 0) {
            return [];
        }

        $result = [];
        for ($i = 0; $i < $count; $i++) {
            foreach($subParams as $v) {
                $result[] = MockRouter::generateData($v);
            }
        }

        return $result;
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
            // 随机数量
            list($min, $max) = explode('-', $rule);
        } else {
            // 固定数量
            $min = $max = $rule;
        }

        $min = self::checkLimits($min, self::$min, self::$max);
        $max = self::checkLimits($max, self::$min, self::$max);

        if ($min == $max) {
            return $min;
        }

        return mt_rand($min, $max);
    }
}
