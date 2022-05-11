<?php

namespace App\Modules\Mock\Mocks;

/**
 * IP
 */
class IP
{
    /**
     * a-e类ip地址范围
     *
     * @var array
     */
    public static $types = [
        'a' => [0, 127],
        'b' => [128, 191],
        'c' => [192, 223],
        'd' => [224, 239],
        'e' => [240, 255]
    ];

    /**
     * 随机生成一个IP地址
     *
     * @param string $type ip类型a-e
     * @return string
     */
    public static function random($type = null)
    {
        if ($type) {
            $type = strtolower($type);
        }

        if (!$type or !isset(self::$types[$type])) {
            $type = array_rand(self::$types);
        }

        return self::generate($type);
    }

    /**
     * 随机生成一个IP地址
     *
     * @param string $type ip类型
     * @return string
     */
    protected static function generate($type)
    {
        $arr = [];
        $arr[] = mt_rand(self::$types[$type][0], self::$types[$type][1]);
        $arr[] = mt_rand(0, 255);
        $arr[] = mt_rand(0, 255);
        if ($type == 'e') {
            $arr[] = mt_rand(0, 254);
        } else {
            $arr[] = mt_rand(0, 255);
        }
        return implode('.', $arr);
    }
}
