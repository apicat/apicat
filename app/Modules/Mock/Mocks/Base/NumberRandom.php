<?php

namespace App\Modules\Mock\Mocks\Base;

/**
 * 随机数字
 */
class NumberRandom
{
    /**
     * 随机整数
     *
     * @param int $min 最小值
     * @param int $max 最大值
     * @return int
     */
    public static function randomInt(int $min = 0, int $max = 1000)
    {
        return mt_rand($min, $max);
    }

    /**
     * 随机小数
     *
     * @param int $min 整数最小值
     * @param int $max 整数最大值
     * @param int $dmin 小数最小位数
     * @param int $dmax 小数最大位数
     * @return float
     */
    public static function randomFloat(int $min = 0, int $max = 1000, int $dmin = 1, int $dmax = 10)
    {
        $left = mt_rand($min, $max);
        $len = mt_rand($dmin, $dmax);
        $right = mt_rand(1, (pow(10, $len) - 1));

        return (float)($left . '.' . str_pad($right, $len, '0', STR_PAD_LEFT));
    }
}
