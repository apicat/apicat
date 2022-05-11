<?php

namespace App\Modules\Mock\Mocks;

/**
 * 手机号码
 */
class Mobile
{
    /**
     * 手机号码前缀
     *
     * @var array
     */
    public static $prefix = [
        '134', '135', '136', '137', '138', '139', '150', '151', '152',
        '157', '158', '159', '178', '182', '183', '184', '187', '188',
        '195', '197', '198', '130', '131', '132', '155', '156', '166',
        '171', '175', '176', '185', '186', '196', '133', '153', '173',
        '177', '180', '181', '189', '190', '191', '193', '199', '192'
    ];

    /**
     * 随机生成手机号
     *
     * @return string
     */
    public static function random()
    {
        $prefix = array_rand(self::$prefix);
        $number = mt_rand(0, 99999999);
        $result = self::$prefix[$prefix] . str_pad($number, 8, '0', STR_PAD_LEFT);
        return $result;
    }
}
