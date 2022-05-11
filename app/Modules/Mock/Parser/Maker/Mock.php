<?php

namespace App\Modules\Mock\Parser\Maker;

class Mock
{
    /**
     * 范围校验
     *
     * @param int $value 校验值
     * @param int $min 最小值
     * @param int $max 最大值
     * @return int
     */
    public static function checkLimits($value, int $min, int $max)
    {
        $value = (int)$value;

        if ($value < $min) {
            return $min;
        } elseif ($value > $max) {
            return $max;
        } else {
            return $value;
        }
    }

    /**
     * 校验是否在数组中
     *
     * @param string $value 校验值
     * @param array $arr 数组
     * @param string $default 默认值
     * @return string
     */
    public static function checkIn($value, array $arr, $default)
    {
        if (!in_array($value, $arr)) {
            return $default;
        }
        
        return $value;
    }
}
