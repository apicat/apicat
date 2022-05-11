<?php

namespace App\Modules\Mock\Mocks;

use App\Models\Location;

/**
 * 身份证号码
 */
class IDCard
{
    /**
     * 加权因子
     *
     * @var array
     */
    public static $factor = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2];

    /**
     * 校验码
     *
     * @var array
     */
    public static $checkCode = ['1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'];

    /**
     * 每个月的天数（从索引0开始算一月）
     *
     * @var array
     */
    public static $monthDays = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];

    /**
     * 随机生成一个身份证号码
     *
     * @return string
     */
    public static function random()
    {
        $offset = mt_rand(0, 3341);
        $areaID  = Location::where('leveltype', 3)->offset($offset)->limit(1)->value('id');

        $year = mt_rand(1930, 2020);

        $monthIndex = array_rand(self::$monthDays);
        $day = mt_rand(1, self::$monthDays[$monthIndex]);
        $day = str_pad($day, 2, '0', STR_PAD_LEFT);

        $month = $monthIndex + 1;
        $month = str_pad($month, 2, '0', STR_PAD_LEFT);

        $tail = mt_rand(1, 999);
        $tail = str_pad($tail, 3, '0', STR_PAD_LEFT);

        $idcard = $areaID . $year . $month . $day . $tail;

        $sum = 0;
        for ($i = 0; $i < 17; $i++) {
            $sum += ($idcard[$i] * self::$factor[$i]);
        }

        $remainder = $sum % 11;
        $checkCode = self::$checkCode[$remainder];

        return $idcard . $checkCode;
    }
}
