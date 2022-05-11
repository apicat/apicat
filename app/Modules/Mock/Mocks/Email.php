<?php

namespace App\Modules\Mock\Mocks;

/**
 * 电子邮箱
 */
class Email
{
    /**
     * 常用邮箱后缀
     *
     * @var array
     */
    public static $commonSubffix = [
        'gmail.com', 'yahoo.com', 'yahoo.com.cn', 'msn.com', 'hotmail.com',
        'aol.com', 'ask.com', 'live.com', 'qq.com', 'vip.qq.com', 'foxmail.com',
        'sina.com', 'vip.sina.com', 'sina.com.cn', 'vip.sina.com.cn', 'sohu.com',
        'outlook.com', 'icloud.com', 'korea.com', 'opera.com', 'zoho.com', 'tom.com',
        '0355.net', '163.com', 'vip.163.com', '163.net', 'netease.com', '126.com',
        '126.net', '128.com', '189.com', '263.net', '3721.net', '56.com', 'eyou.com',
        '21cn.com', 'chianren.com', 'yeah.net', 'googlemail.com'
    ];

    /**
     * 随机生成一个邮箱地址
     *
     * @return string
     */
    public static function random()
    {
        $len = mt_rand(1, 64);

        if ($len > 1) {
            $hasSpecialChar = mt_rand(0, 1);
        } else {
            $hasSpecialChar = false;
        }

        $local = Base\StringRandom::lowercaseString($len);

        if ($hasSpecialChar) {
            $specialChars = '._-+';
            $index = mt_rand(0, 3);
            $position = mt_rand(1, $len);
            
            $local = substr($local , 0, $position) . $specialChars[$index] . substr($local, $position, $len - $position);
        }

        $subffixIndex = array_rand(self::$commonSubffix);

        return $local . '@' . self::$commonSubffix[$subffixIndex];
    }
}
