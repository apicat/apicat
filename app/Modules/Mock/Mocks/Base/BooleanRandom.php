<?php

namespace App\Modules\Mock\Mocks\Base;

/**
 * 随机布尔值
 */
class BooleanRandom
{
    /**
     * 随机生成true
     *
     * @param int $probability 0-100概率
     * @return boolean
     */
    public static function randomTrue($probability = 50)
    {
        if ($probability  == 0) {
            return false;
        }

        if ($probability  == 100) {
            return true;
        }

        $val = mt_rand(1, 100);
        if ($val <= $probability) {
            return true;
        } else {
            return false;
        }
    }
}
