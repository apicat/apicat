<?php

namespace App\Modules\Mock\Mocks;

/**
 * 链接地址
 */
class Url
{
    /**
     * 随机生成一个链接地址
     *
     * @return string
     */
    public static function random()
    {
        $protocols = ['http://', 'https://'];
        $index = mt_rand(0, 1);
        $domain = Domain::random();

        return $protocols[$index] . $domain;
    }
}
