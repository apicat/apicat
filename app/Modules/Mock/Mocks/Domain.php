<?php

namespace App\Modules\Mock\Mocks;

/**
 * 域名
 */
class Domain
{
    /**
     * 常用域名后缀
     *
     * @var array
     */
    public static $commonSubffix = [
        '.top', '.cn', '.com', '.net', '.xyz', '.icu', '.shop', '.club', '.cc',
        '.vip', '.ltd', '.site', '.ink', '.pub', '.co', '.cloud', '.ren', '.asia',
        '.work', '.fit', '.biz', '.art', '.love', '.online', '.info', '.wang', '.fans',
        '.store', '.red', '.mobi', '.kim', '.com.cn', '.net.cn', '.gov.cn', '.link', '.tech',
        '.pro', '.xin'
    ];

    /**
     * 随机生成一个域名
     *
     * @return string
     */
    public static function random()
    {
        $hasChild = mt_rand(0, 1);
        $main = Base\StringRandom::randomLenLower(4, 20);
        $index = array_rand(self::$commonSubffix);

        if (!$hasChild) {
            return $main . self::$commonSubffix[$index];
        }

        return 'www.' . $main . self::$commonSubffix[$index];
    }
}
